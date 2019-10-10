package earveldaja

import (
	"bytes"
	"crypto/rsa"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/zdarovich/einvoice/core/erro"
	"github.com/zdarovich/einvoice/core/models"
	"io/ioutil"
	"time"

	log "github.com/sirupsen/logrus"
)

type service struct {
	cli   *HTTPClient
	creds providerCredentials
}

type providerCredentials struct {
	privateKey   *rsa.PrivateKey
	providerCode string
}

func NewProviderCredentials(privateKey *rsa.PrivateKey, providerCode string) providerCredentials {
	return providerCredentials{
		privateKey:   privateKey,
		providerCode: providerCode,
	}
}

func NewService(cli *HTTPClient, credentials providerCredentials) *service {
	return &service{cli: cli, creds: credentials}
}

func (s *service) PostDocuments(param map[string]string, docs []models.SaleDocument) ([]models.ExportMapping, error) {
	result := s.sendPipeline(docs)
	return result, nil
}

func (s *service) GetDocumentsStatus(param map[string]string, exported []models.ExportMapping, since time.Time) ([]models.ExportMapping, error) {
	gilq, err := generateGetInvoiceListQuery(s.creds.providerCode, s.creds.privateKey, since)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	env := Envelope{
		Xsi:     "http://www.w3.org/2001/XMLSchema-instance",
		Xsd:     "http://www.w3.org/2001/XMLSchema",
		Soap:    "http://schemas.xmlsoap.org/soap/envelope/",
		XMLName: soapenv,
		Body: Body{
			GetInvoiceListService: &GetInvoiceListService{
				Xmlns:               "http://arireg-rmp.x-road.ee/producer/",
				GetInvoiceListQuery: *gilq,
			},
		}}
	payload, err := xml.MarshalIndent(env, "", "    ")
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Debug(string(payload))

	req, err := s.cli.NewRequest("POST", bytes.NewBuffer(payload))
	if err != nil {
		log.Error(err)
		return nil, err
	}
	resp, err := s.cli.Do(req)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Debug(string(body))

	risp := &ResponseEnvelope{}

	if err := xml.Unmarshal(body, risp); err != nil {
		return nil, err
	}
	if risp.Body.Fault != nil {
		return nil, erro.NewProviderError(risp.Body.Fault.Code, risp.Body.Fault.Description)
	}

	return getUpdateResponse(exported, &risp.Body.GetInvoiceListServiceResponse.GetInvoiceListQueryResponse), nil
}

func (s *service) GetSaleInvoices(param map[string]string, since time.Time) ([]models.SaleDocument, error) {
	return nil, errors.New("implement me")
}

func (s *service) GetPurchaseInvoices(param map[string]string, since time.Time) ([]models.SaleDocument, error) {
	return nil, errors.New("implement me")
}

func (s *service) sendDocument(riq *ReceiveInvoiceQuery) (*ReceiveInvoiceServiceResponse, error) {
	if riq == nil {
		return nil, errors.New("receive invoice query is nil")
	}

	env := Envelope{
		Xsi:     "http://www.w3.org/2001/XMLSchema-instance",
		Xsd:     "http://www.w3.org/2001/XMLSchema",
		Soap:    "http://schemas.xmlsoap.org/soap/envelope/",
		XMLName: soapenv,
		Body: Body{
			ReceiveInvoiceService: &ReceiveInvoiceService{
				Xmlns:               "http://arireg-rmp.x-road.ee/producer/",
				ReceiveInvoiceQuery: *riq,
			},
		}}
	payload, err := xml.MarshalIndent(env, "", "    ")
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Debug(string(payload))

	req, err := s.cli.NewRequest("POST", bytes.NewBuffer(payload))
	if err != nil {
		log.Error(err)
		return nil, err
	}
	resp, err := s.cli.Do(req)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Debug(string(body))

	risp := &ResponseEnvelope{}

	if err := xml.Unmarshal(body, risp); err != nil {
		return nil, err
	}
	if risp.Body.Fault != nil {
		return nil, erro.NewProviderError(risp.Body.Fault.Code, risp.Body.Fault.Description)
	} else if risp.Body.ReceiveInvoiceServiceResponse.Response.Status == "FAIL" {
		return nil, erro.NewProviderError(risp.Body.ReceiveInvoiceServiceResponse.Response.Status, risp.Body.ReceiveInvoiceServiceResponse.Response.Message)
	}

	return risp.Body.ReceiveInvoiceServiceResponse, nil
}

func (s *service) createForwardingContract(cfcq *CreateForwardingContractQuery) (*CreateForwardingContractQueryResponse, error) {
	if cfcq == nil {
		return nil, errors.New("create forwarding contract query is nil")
	}

	env := Envelope{
		Xsi:     "http://www.w3.org/2001/XMLSchema-instance",
		Xsd:     "http://www.w3.org/2001/XMLSchema",
		Soap:    "http://schemas.xmlsoap.org/soap/envelope/",
		XMLName: soapenv,
		Body: Body{
			CreateForwardingContractService: &CreateForwardingContractService{
				Xmlns:                         "http://arireg-rmp.x-road.ee/producer/",
				CreateForwardingContractQuery: *cfcq,
			},
		}}
	payload, err := xml.MarshalIndent(env, "", "    ")
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Debug(string(payload))

	req, err := s.cli.NewRequest("POST", bytes.NewBuffer(payload))
	if err != nil {
		log.Error(err)
		return nil, err
	}
	resp, err := s.cli.Do(req)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Debug(string(body))

	risp := &ResponseEnvelope{}

	if err := xml.Unmarshal(body, risp); err != nil {
		return nil, err
	}
	if risp.Body.Fault != nil {
		return nil, erro.NewProviderError(risp.Body.Fault.Code, risp.Body.Fault.Description)
	}

	return &risp.Body.CreateForwardingContractServiceResponse.Response, nil
}

func (s *service) sendPipeline(documents []models.SaleDocument) []models.ExportMapping {
	done := make(chan interface{})
	defer close(done)

	type ReceiveResult struct {
		riq *ReceiveInvoiceQuery
	}
	generator := func(done <-chan interface{}, documents []models.SaleDocument) <-chan ReceiveInvoiceQuery {
		receiveStream := make(chan ReceiveInvoiceQuery)
		go func(docs []models.SaleDocument) {
			defer close(receiveStream)
			for _, doc := range docs {
				log.Trace(fmt.Sprintf("generate receive query docID %d", doc.ID))
				riq, err := generateReceiveInvoiceQuery(s.creds.providerCode, s.creds.privateKey, &doc)
				if err != nil {
					log.Error(err)
					continue
				}
				log.Trace(riq)

				select {
				case <-done:
					return
				case receiveStream <- *riq:
				}
			}

		}(documents)
		return receiveStream
	}
	type ReceiveInvoiceServiceResponseResult struct {
		docID int
		riq   *ReceiveInvoiceServiceResponse
		err   error
	}
	send := func(done <-chan interface{}, receiveStream <-chan ReceiveInvoiceQuery) <-chan ReceiveInvoiceServiceResponseResult {
		sendStream := make(chan ReceiveInvoiceServiceResponseResult)
		go func() {
			defer close(sendStream)
			for riq := range receiveStream {
				log.Trace(fmt.Sprintf("send receive query docID %d", riq.DocId))
				riqr, err := s.sendDocument(&riq)
				select {
				case <-done:
					return
				case sendStream <- ReceiveInvoiceServiceResponseResult{docID: riq.DocId, riq: riqr, err: err}:
				}
			}
		}()
		return sendStream
	}

	orDone := func(done <-chan interface{}, c <-chan ReceiveInvoiceServiceResponseResult) <-chan ReceiveInvoiceServiceResponseResult {
		valStream := make(chan ReceiveInvoiceServiceResponseResult)
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
	receiveStream := generator(done, documents)
	pipeline := send(done, receiveStream)

	var res []models.ExportMapping
	for result := range orDone(done, pipeline) {
		exportRes := new(models.ExportMapping)
		exportRes.DocumentID = result.docID
		exportRes.Err = result.err
		exportRes.Err = result.err
		exportRes.DocumentProviderID = result.riq.Response.SenderInvoiceID
		res = append(res, *exportRes)

	}
	return res
}
