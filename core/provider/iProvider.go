package provider

import (
	"github.com/zdarovich/einvoice/core/models"
	"time"
)

type IProvider interface {
	PostDocuments(param map[string]string, docs []models.SaleDocument) ([]models.ExportMapping, error)
	GetDocumentsStatus(param map[string]string, exported []models.ExportMapping, since time.Time) ([]models.ExportMapping, error)
	GetSaleInvoices(param map[string]string, since time.Time) ([]models.SaleDocument, error)
	GetPurchaseInvoices(param map[string]string, since time.Time) ([]models.SaleDocument, error)
}
