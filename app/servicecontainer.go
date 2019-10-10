package app

import (
	"database/sql"
	"fmt"
	_ "github.com/CovenantSQL/go-sqlite3-encrypt"
	"github.com/gammazero/workerpool"
	_ "github.com/go-sql-driver/mysql"
	"github.com/zdarovich/einvoice/config"
	"github.com/zdarovich/einvoice/controllers"
	"github.com/zdarovich/einvoice/core/provider/earveldaja"
	"github.com/zdarovich/einvoice/repositories"
	"github.com/zdarovich/einvoice/services"
	"sync"
)

type IServiceContainer interface {
	InjectDocumentController(conf *config.Config, db *sql.DB) (*controllers.DocumentController, error)
	InjectConfiguration() (*config.Config, error)
	InjectDatabase(config *config.Config) (*sql.DB, error)
}

type kernel struct{}

func (k *kernel) InjectConfiguration() (*config.Config, error) {
	return config.New()
}

func (k *kernel) InjectDatabase(config *config.Config) (*sql.DB, error) {
	return sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/invoice_db", config.Db.User, config.Db.Password, config.Db.IP, config.Db.Port))
}

func (k *kernel) InjectDocumentController(conf *config.Config, db *sql.DB) (*controllers.DocumentController, error) {
	privKey, err := GetPrivateKey("/cert/server_earveldaja.key")
	if err != nil {
		return nil, err
	}
	creds := earveldaja.NewProviderCredentials(privKey, conf.ProviderData.EarveldajaProviderCode)
	cli := earveldaja.NewHTTPClient(conf.ProviderData.EarveldajaWSDLURL)
	service := earveldaja.NewService(cli, creds)
	repo := &repositories.DocumentRepository{db}
	workerPool := workerpool.New(conf.WorkerPoolGoroutines)
	documentService := &services.DocumentService{service, repo, workerPool, conf.ProviderData.EArveldajaStatusUpdateDelaySec}
	documentController := controllers.DocumentController{documentService}
	return &documentController, nil
}

var (
	k             *kernel
	containerOnce sync.Once
)

func ServiceContainer() IServiceContainer {
	if k == nil {
		containerOnce.Do(func() {
			k = &kernel{}
		})
	}
	return k
}
