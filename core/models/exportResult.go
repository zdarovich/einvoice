package models

type ExportMapping struct {
	DocumentID         int
	DocumentProviderID string
	Err                error
}
