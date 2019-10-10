package interfaces

import "github.com/zdarovich/einvoice/models"

type IDocumentRepository interface {
	GetExportRecord(clientCode string, documentId string, provider string) (models.ExportRecord, error)
	CreateExportRecord(clientCode string, dataHash string, documentId string, provider string, status models.Status) error
	UpdateExportRecordStatus(id string, status models.Status) error
	GetImportRecordByHash(clientCode string, dataHash string, provider string) (models.ExportRecord, error)
	CreateImportRecord(clientCode string, dataHash string, documentId string, provider string) error
}
