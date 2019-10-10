package controllers

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	core "github.com/zdarovich/einvoice/core/models"
	"github.com/zdarovich/einvoice/interfaces"
	"github.com/zdarovich/einvoice/models"
	"net/http"
)

type DocumentController struct {
	interfaces.IDocumentService
}

func (dc *DocumentController) PostDocuments(res http.ResponseWriter, req *http.Request) {
	logrus.Info("Post documents request")
	request := new([]core.SaleDocument)
	if err := json.NewDecoder(req.Body).Decode(request); err != nil {
		errorHandler(res, err)
		return
	}
	tres, err := dc.ExportDocuments(*request)
	if err != nil {
		errorHandler(res, err)
		return
	}
	response(res, tres)

	// schedule document status update
	dc.UpdateDocumentStatus(tres.DocumentStatus, tres.DateTime)
}

func (dc *DocumentController) GetDocumentStatus(res http.ResponseWriter, req *http.Request) {
	logrus.Info("Get documents status")
	request := new(models.DocumentStatusRequest)
	if err := json.NewDecoder(req.Body).Decode(request); err != nil {
		errorHandler(res, err)
		return
	}
	tres, err := dc.ProvideDocumentStatus(request)
	if err != nil {
		errorHandler(res, err)
		return
	}
	response(res, tres)
}

func (dc *DocumentController) GetSaleDocuments(res http.ResponseWriter, req *http.Request) {
	logrus.Info("Get sale documents request")
	request := new(models.DocumentImportRequest)
	if err := json.NewDecoder(req.Body).Decode(request); err != nil {
		errorHandler(res, err)
		return
	}
	tres, err := dc.ImportSaleInvoice(request)
	if err != nil {
		errorHandler(res, err)
		return
	}
	response(res, tres)
}

func (dc *DocumentController) GetPurchaseDocuments(res http.ResponseWriter, req *http.Request) {
	logrus.Info("Get sale documents request")
	request := new(models.DocumentImportRequest)
	if err := json.NewDecoder(req.Body).Decode(request); err != nil {
		errorHandler(res, err)
		return
	}
	tres, err := dc.ImportSaleInvoice(request)
	if err != nil {
		errorHandler(res, err)
		return
	}
	response(res, tres)
}
