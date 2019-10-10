package models

const (
	Accepted Status = "ACCEPTED"
	Declined Status = "DECLINED"
	Pending  Status = "PENDING"
)

type ExportRecord struct {
	Id           string
	Status       string
	ProviderName string
	ClientCode   string
	DocumentHash string
	DocumentId   int
}

type ImportRecord struct {
	Id           int
	ProviderName string
	ClientCode   string
	DocumentHash string
	DocumentId   int
}

type Status string
