package models

import (
	core "github.com/zdarovich/einvoice/core/models"
	"time"
)

type ServerErrorResponse struct {
	Message    string    `json:"message"`
	StatusCode int       `json:"statusCode"`
	DateTime   time.Time `json:"dateTime"`
}

type DocumentResponse struct {
	DocumentStatus []core.ExportMapping `json:"documentStatus,omitempty"`
	Documents      []core.SaleDocument  `json:"documents,omitempty"`
	DateTime       time.Time            `json:"dateTime,omitempty"`
}

type DocumentStatusResponse struct {
	DocumentStatus []ExportRecord `json:"documentStatus,omitempty"`
	DateTime       time.Time      `json:"dateTime,omitempty"`
}

type DocumentStatusRequest struct {
	DocumentIDs []int `json:"documentIDs"`
}

type DocumentImportRequest struct {
	DateTime string `json:"since"`
}
