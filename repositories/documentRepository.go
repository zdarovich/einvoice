package repositories

import (
	"database/sql"
	"github.com/zdarovich/einvoice/models"

	log "github.com/sirupsen/logrus"
)

type DocumentRepository struct {
	Db *sql.DB
}

func (repository *DocumentRepository) GetExportRecord(clientCode string, documentId string, provider string) (models.ExportRecord, error) {
	var er = models.ExportRecord{}
	stmtOuts, err := repository.Db.Prepare("select ID, (SELECT Status from InvoiceStatus WHERE ID = InvoiceStatusID) AS Status, DataHash, DocumentID, ProviderName, ClientCode from ExportRecord where ClientCode = ? AND DocumentID = ? AND ProviderName = ? ")
	defer stmtOuts.Close()
	if err != nil {
		return er, err
	}

	err = stmtOuts.QueryRow(clientCode, documentId, provider).Scan(&er.Id, &er.Status, &er.DocumentHash, &er.DocumentId, &er.ProviderName, &er.ClientCode)

	if err != nil {
		if err == sql.ErrNoRows {
			return er, nil
		}
		return er, err

	}
	return er, nil
}

func (repository *DocumentRepository) GetAllExportRecord(clientCode string, documentId string) (models.ExportRecord, error) {
	var result = models.ExportRecord{}
	stmtOuts, err := repository.Db.Prepare("select ID, (SELECT Status from InvoiceStatus WHERE ID = InvoiceStatusID) AS Status, DataHash, DocumentID, ProviderName, ClientCode from ExportRecord where ClientCode = ? AND DocumentID = ?")
	defer stmtOuts.Close()
	if err != nil {
		return result, err
	}
	log.Debugf("clientCode: %s, documentId: %s", clientCode, documentId)
	err = stmtOuts.QueryRow(clientCode, documentId).Scan(&result.Id, &result.Status, &result.DocumentHash, &result.DocumentId, &result.ProviderName, &result.ClientCode)

	if err != nil {
		if err == sql.ErrNoRows {
			return result, nil
		}
		return result, err

	}
	return result, nil
}

func (repository *DocumentRepository) CreateExportRecord(clientCode string, dataHash string, documentId string, provider string, status models.Status) error {
	stmtIns, err := repository.Db.Prepare("insert into ExportRecord(InvoiceStatusID, DataHash, DocumentID, ProviderName, ClientCode) values((SELECT ID from InvoiceStatus WHERE Status = ?), ?, ?, ?, ?)")
	defer stmtIns.Close()

	if err != nil {
		return err
	}
	_, err = stmtIns.Exec(status, dataHash, documentId, provider, clientCode)
	return err
}

func (repository *DocumentRepository) UpdateExportRecordStatus(id string, status models.Status) error {
	stmtIns, err := repository.Db.Prepare("update ExportRecord set InvoiceStatusID = (SELECT ID from InvoiceStatus WHERE Status = ?) where ID = ?")
	defer stmtIns.Close()
	log.Debugf("status: %s, id: %s", status, id)
	if err != nil {
		return err
	}
	_, err = stmtIns.Exec(status, id)
	return err
}

func (repository *DocumentRepository) GetImportRecordByHash(clientCode string, dataHash string, provider string) (models.ExportRecord, error) {
	var er = models.ExportRecord{}
	stmtOuts, err := repository.Db.Prepare("select * from ImportRecord where ClientCode = ? AND DocumentHash = ? AND ProviderName = ? ")
	if err != nil {
		return er, err
	}
	defer stmtOuts.Close()

	err = stmtOuts.QueryRow(clientCode, dataHash, provider).Scan(&er.Id, &er.DocumentHash, &er.DocumentId, &er.ProviderName, &er.ClientCode)
	if err != nil {
		if err == sql.ErrNoRows {
			return er, nil
		} else {
			return er, err
		}
	}
	return er, nil
}

func (repository *DocumentRepository) CreateImportRecord(clientCode string, dataHash string, documentId string, provider string) error {
	stmtIns, err := repository.Db.Prepare("insert into ImportRecord(DocumentHash, DocumentID, ProviderName, ClientCode) values(?, ?, ?, ?)")
	defer stmtIns.Close()

	if err != nil {
		return err
	}
	_, err = stmtIns.Exec(dataHash, documentId, provider, clientCode)
	return err
}
