package config

import (
	"encoding/json"
	"github.com/zdarovich/einvoice/constants"
	"io/ioutil"
	"os"
)

//Config ...
type Config struct {
	ProdEnv              bool
	ProviderData         *providerData
	Server               *server
	WorkerPoolGoroutines int

	Db *db
}

type server struct {
	Port                 int
	ReadHeaderTimeoutSec int
	ReadTimeoutSec       int
	WriteTimeoutSec      int
	LogFileEnable        bool
	LogFileName          string
}

type providerData struct {
	EarveldajaWSDLURL              string
	EArveldajaStatusUpdateDelaySec int
	EarveldajaProviderCode         string
}

type db struct {
	User     string
	Password string
	IP       string
	Port     string
}

//New / Init / Get
func New() (*Config, error) {
	conf := &Config{}
	_, err := os.Stat(constants.ConfigFile)

	if os.IsNotExist(err) {
		return nil, err
	}

	b, err := ioutil.ReadFile(constants.ConfigFile)
	if err := json.Unmarshal([]byte(string(b)), conf); err != nil {
		return nil, err
	}

	return conf, nil
}
