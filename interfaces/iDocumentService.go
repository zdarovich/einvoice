package interfaces

import (
	core "github.com/zdarovich/einvoice/core/models"
	"github.com/zdarovich/einvoice/models"
	"time"
)

type IDocumentService interface {
	ExportDocuments(docs []core.SaleDocument) (*models.DocumentResponse, error)
	UpdateDocumentStatus(exported []core.ExportMapping, since time.Time)
	ProvideDocumentStatus(req *models.DocumentStatusRequest) (*models.DocumentStatusResponse, error)
	ImportSaleInvoice(req *models.DocumentImportRequest) (*models.DocumentResponse, error)
	ImportPurchaseInvoice(req *models.DocumentImportRequest) (*models.DocumentResponse, error)
}
