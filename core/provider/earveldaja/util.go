package earveldaja

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	b64 "encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/zdarovich/einvoice/core/einvoice"
	"github.com/zdarovich/einvoice/core/models"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	NOT_EXPORTED string = "NOT OK"
	EXPORTED     string = "OK"
)

func generateReceiveInvoiceQuery(providerCode string, privateKey *rsa.PrivateKey, salesDoc *models.SaleDocument) (*ReceiveInvoiceQuery, error) {
	if err := validateFields(providerCode, privateKey, salesDoc); err != nil {
		return nil, err
	}
	pdf, err := generatePdf(salesDoc)
	if err != nil {
		return nil, err
	}
	pdfBase64 := b64.StdEncoding.EncodeToString(pdf)
	einv, err := einvoice.NewEinvoice(salesDoc, salesDoc.Customer, salesDoc.CompanyInfo, salesDoc.VatRate)
	if err != nil {
		return nil, err
	}
	xmlBytes, err := xml.Marshal(einv)
	if err != nil {
		return nil, err
	}
	xmlBase64 := b64.StdEncoding.EncodeToString(xmlBytes)
	timeNow := time.Now()
	transferTime := timeNow.Format("2006-01-02T15:04:05.000-07:00")
	signature, err := generateReceiveSignature(timeNow, providerCode, xmlBase64, pdfBase64, privateKey)
	if err != nil {
		return nil, err
	}
	return &ReceiveInvoiceQuery{
		DocId:            salesDoc.ID,
		ProviderCode:     providerCode,
		TransferDatetime: transferTime,
		SenderName:       salesDoc.CompanyInfo.Name,
		SenderRegcode:    salesDoc.CompanyInfo.Code,
		ReceiverName:     salesDoc.Customer.FullName,
		ReceiverRegcode:  salesDoc.Customer.Code,
		SenderInvoiceID:  salesDoc.Number,
		InvoicePdfFile:   pdfBase64,
		InvoiceXMLFile:   xmlBase64,
		InvoiceSignature: signature,
	}, nil
}

func validateFields(providerCode string, privateKey *rsa.PrivateKey, salesDoc *models.SaleDocument) error {
	if providerCode == "" {
		return errors.New("provider code is empty")
	} else if privateKey == nil {
		return errors.New("private key is nil")
	} else if salesDoc == nil {
		return errors.New("sales document is nil")
	} else if salesDoc.Customer == nil {
		return errors.New("customer is nil")
	} else if salesDoc.CompanyInfo == nil {
		return errors.New("companyInfo is nil")
	} else if salesDoc.VatRate == nil {
		return errors.New("vatRate is nil")
	} else {
		return nil
	}
}

func generateGetInvoiceListQuery(providerCode string, privateKey *rsa.PrivateKey, exportStart time.Time) (*GetInvoiceListQuery, error) {
	result := &GetInvoiceListQuery{}
	timeNow := time.Now()
	transferTime := timeNow.Format("2006-01-02T15:04:05.000-07:00")
	fromDatetime := exportStart.Format("2006-01-02T15:04:05.000-07:00")
	toDateTime := exportStart.Add(1 * time.Hour).Format("2006-01-02T15:04:05.000-07:00")
	signature, err := generateGetInvoiceSignature(timeNow, providerCode, exportStart, privateKey)
	if err != nil {
		return nil, err
	}
	result.Signature = signature
	result.TransferDatetime = transferTime
	result.FromDatetime = fromDatetime
	result.ProviderCode = providerCode
	result.ToDatetime = toDateTime
	return result, nil
}

func generateCreateForwardingContractQuery(providerCode string, privateKey *rsa.PrivateKey, receiverRegcode string) (*CreateForwardingContractQuery, error) {
	result := &CreateForwardingContractQuery{}
	timeNow := time.Now()
	transferTime := timeNow.Format("2006-01-02T15:04:05.000-07:00")
	signature, err := generateCreateForwardingContractSignature(timeNow, providerCode, receiverRegcode, privateKey)
	if err != nil {
		return nil, err
	}
	result.Signature = signature
	result.TransferDatetime = transferTime
	result.ProviderCode = providerCode
	result.ForwarderRegcode = receiverRegcode
	return result, nil
}

func getUpdateResponse(exported []models.ExportMapping, fiqr *GetInvoiceListQueryResponse) []models.ExportMapping {
	isAccepted := make(map[int]bool)
	for _, invoiceListing := range fiqr.InvoiceListing.Invoices {
		isAccepted[invoiceListing.InvoiceId] = true
	}
	for _, mapping := range exported {
		providerDocID, err := strconv.Atoi(mapping.DocumentProviderID)
		if err != nil {
			log.Error("wrong received id")
			continue
		}
		if isAccepted[providerDocID] {
			mapping.Err = errors.New("not exported")
		} else {
			mapping.Err = nil
		}
	}
	return exported
}

func generatePdf(salesDoc *models.SaleDocument) ([]byte, error) {
	// not implemented yet
	token := make([]byte, 4)
	_, err := rand.Read(token)
	if err != nil {
		return nil, err
	}
	return token, nil

}

func generateReceiveSignature(transferTime time.Time, providerCode string, xml string, pdf string, privateKey *rsa.PrivateKey) (string, error) {
	transferTimeStr := transferTime.Format("02.01.2006 15:04:05")
	message := fmt.Sprintf("%s#%s#%s#%s", providerCode, transferTimeStr, xml, pdf)
	return generateSignature(message, privateKey)
}

func generateGetInvoiceSignature(transferTime time.Time, providerCode string, fromDatetime time.Time, privateKey *rsa.PrivateKey) (string, error) {
	transferTimeStr := transferTime.Format("02.01.2006 15:04:05")
	fromDatetimeStr := fromDatetime.Format("02.01.2006 15:04:05")
	message := fmt.Sprintf(`%s#%s#%s`, providerCode, transferTimeStr, fromDatetimeStr)
	return generateSignature(message, privateKey)
}

func generateCreateForwardingContractSignature(transferTime time.Time, providerCode string, receiverRegcode string, privateKey *rsa.PrivateKey) (string, error) {
	transferTimeStr := transferTime.Format("02.01.2006 15:04:05")
	message := fmt.Sprintf(`%s#%s#%s`, providerCode, transferTimeStr, receiverRegcode)
	return generateSignature(message, privateKey)
}

func generateSignature(message string, privateKey *rsa.PrivateKey) (string, error) {
	msg := []byte(message)

	hashed := sha256.Sum256(msg)

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return "", err
	}
	sig := base64.StdEncoding.EncodeToString(signature)
	return sig, nil
}
