package services

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gammazero/workerpool"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	core "github.com/zdarovich/einvoice/core/models"
	"github.com/zdarovich/einvoice/core/provider"
	"github.com/zdarovich/einvoice/interfaces"
	"github.com/zdarovich/einvoice/models"
	"strconv"
	"time"
)

type DocumentService struct {
	provider.IProvider
	interfaces.IDocumentRepository
	WorkerPool  *workerpool.WorkerPool
	UpdateDelay int
}

func (ds *DocumentService) ExportDocuments(docs []core.SaleDocument) (*models.DocumentResponse, error) {
	param := make(map[string]string)
	result, err := ds.PostDocuments(param, docs)
	if err != nil {
		return nil, err
	}

	exportMap := make(map[int]*core.ExportMapping)
	for _, res := range result {
		exportMap[res.DocumentID] = &res
	}

	done := make(chan interface{})
	defer close(done)

	type DocumentHash struct {
		doc  *core.SaleDocument
		hash string
	}
	hashGenerator := func(done <-chan interface{}, documents []core.SaleDocument) <-chan DocumentHash {
		receiveStream := make(chan DocumentHash)
		go func(docs []core.SaleDocument) {
			defer close(receiveStream)
			for _, doc := range docs {
				logrus.Trace(fmt.Sprintf("generate hash of docID %d", doc.ID))
				docByte, err := json.Marshal(doc)
				if err != nil {
					logrus.Error(err)
				}
				hash := sha256.Sum256(docByte)
				hashBase64 := base64.StdEncoding.EncodeToString(hash[:])
				docHash := new(DocumentHash)
				docHash.doc = &doc
				docHash.hash = hashBase64
				select {
				case <-done:
					return
				case receiveStream <- *docHash:
				}
			}

		}(documents)
		return receiveStream
	}
	type ExportRecordResult struct {
		docID int
		err   error
	}

	saveDB := func(done <-chan interface{}, receiveStream <-chan DocumentHash) <-chan ExportRecordResult {
		sendStream := make(chan ExportRecordResult)
		go func() {
			defer close(sendStream)
			for rs := range receiveStream {
				logrus.Trace(fmt.Sprintf("save export result of docID %d", rs.doc.ID))
				var err error
				if exportMap[rs.doc.ID].Err != nil {
					err = ds.createUpdateExportRecords(rs.doc, rs.hash, models.Pending)
				} else {
					err = ds.createUpdateExportRecords(rs.doc, rs.hash, models.Declined)
				}
				select {
				case <-done:
					return
				case sendStream <- ExportRecordResult{docID: rs.doc.ID, err: err}:
				}
			}
		}()
		return sendStream
	}

	orDone := func(done <-chan interface{}, c <-chan ExportRecordResult) <-chan ExportRecordResult {
		valStream := make(chan ExportRecordResult)
		go func() {
			defer close(valStream)
			for {
				select {
				case <-done:
					return
				case v, ok := <-c:
					if ok == false {
						return
					}
					select {
					case valStream <- v:
					case <-done:
					}
				}
			}
		}()
		return valStream
	}
	receiveStream := hashGenerator(done, docs)
	pipeline := saveDB(done, receiveStream)
	for result := range orDone(done, pipeline) {
		logrus.Error(errors.Wrap(result.err, fmt.Sprintf("export record of document ID %d was not saved in db", result.docID)))
	}
	resp := new(models.DocumentResponse)
	resp.DocumentStatus = result
	resp.DateTime = time.Now()
	return resp, nil
}

func (ds *DocumentService) createUpdateExportRecords(document *core.SaleDocument, hash string, status models.Status) error {
	res, err := ds.GetExportRecord("", strconv.Itoa(document.ID), "")
	if err != nil {
		return errors.Wrap(err, "could not get export record")
	}
	if res.Id == "" {
		if err := ds.CreateExportRecord("", hash, strconv.Itoa(document.ID), "", status); err != nil {
			return errors.Wrap(err, "one of the exported docs was not logged")
		}
	} else {
		if err := ds.UpdateExportRecordStatus(res.Id, status); err != nil {
			return errors.Wrap(err, "one of the exported docs was not logged")
		}
	}
	return nil
}

func (ds *DocumentService) ProvideDocumentStatus(req *models.DocumentStatusRequest) (*models.DocumentStatusResponse, error) {
	if len(req.DocumentIDs) == 0 {
		return nil, errors.New("no document ids provided")
	}

	var result []models.ExportRecord
	for _, id := range req.DocumentIDs {
		erecord, err := ds.GetExportRecord("", strconv.Itoa(id), "")
		if err != nil {
			return nil, errors.New("unable to get document statuses")
		}
		result = append(result, erecord)
	}
	response := new(models.DocumentStatusResponse)
	response.DateTime = time.Now()
	response.DocumentStatus = result
	return response, nil
}

func (ds *DocumentService) UpdateDocumentStatus(exported []core.ExportMapping, since time.Time) {
	ds.WorkerPool.Submit(func() {
		if len(exported) == 0 {
			logrus.Info("no exported documents for update")
			return
		}
		uuid := fmt.Sprint(getGID())
		logrus.Infof("goroutine %s: start", uuid)
		defer logrus.Infof("goroutine %s: completed", uuid)

		logrus.Infof("goroutine %s: sleep for %d sec", uuid, ds.UpdateDelay)

		time.Sleep(time.Duration(ds.UpdateDelay) * time.Second)

		logrus.Infof(fmt.Sprintf("goroutine %s: update export status", uuid))

		param := make(map[string]string)
		result, err := ds.GetDocumentsStatus(param, exported, since)
		if err != nil {
			logrus.Error(err)
			return
		}
		for _, mapping := range result {
			res, err := ds.GetExportRecord("", strconv.Itoa(mapping.DocumentID), "")
			if err != nil {
				logrus.Error(err)
				continue
			}
			var status models.Status
			if mapping.Err != nil {
				status = models.Declined
			} else {
				status = models.Accepted
			}
			if err := ds.UpdateExportRecordStatus(res.Id, status); err != nil {
				logrus.Error(err)
				continue
			}
		}

	})
}

func (ds *DocumentService) ImportSaleInvoice(req *models.DocumentImportRequest) (*models.DocumentResponse, error) {
	since, err := time.Parse("2 Jan 2006 15:04:05", req.DateTime)
	if err != nil {
		return nil, errors.New("since datetime not provided")
	}
	param := make(map[string]string)
	result, err := ds.GetSaleInvoices(param, since)
	if err != nil {
		return nil, err
	}
	resp := new(models.DocumentResponse)
	resp.Documents = result
	resp.DateTime = time.Now()
	return resp, nil
}

func (ds *DocumentService) ImportPurchaseInvoice(req *models.DocumentImportRequest) (*models.DocumentResponse, error) {
	since, err := time.Parse("2 Jan 2006 15:04:05", req.DateTime)
	if err != nil {
		return nil, errors.New("since datetime not provided")
	}
	param := make(map[string]string)

	result, err := ds.GetPurchaseInvoices(param, since)
	if err != nil {
		return nil, err
	}
	resp := new(models.DocumentResponse)
	resp.Documents = result
	resp.DateTime = time.Now()
	return resp, nil
}
