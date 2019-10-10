package app

import (
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/zdarovich/einvoice/config"
	"github.com/zdarovich/einvoice/middleware"

	"net/http"
	"sync"
)

type IMuxRouter interface {
	InitRouter(config *config.Config, db *sql.DB) (*mux.Router, error)
}
type router struct{}

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

func (router *router) InitRouter(config *config.Config, db *sql.DB) (*mux.Router, error) {

	documentController, err := ServiceContainer().InjectDocumentController(config, db)
	if err != nil {
		return nil, err
	}

	middlewareChain := []middleware.Func{
		middleware.Logging,
		middleware.Recovery,
	}

	routes := []route{
		route{
			"Send sale document",
			"POST",
			"/api/v1/sale/export",
			documentController.PostDocuments,
		},
		route{
			"Request sale document from provider",
			"POST",
			"/api/v1/sale/import",
			documentController.GetSaleDocuments,
		},
		route{
			"Request purchase document from provider",
			"POST",
			"/api/v1/purchase/import",
			documentController.GetPurchaseDocuments,
		},
	}
	r := mux.NewRouter().StrictSlash(true)

	for _, route := range routes {

		r.Methods(route.Method, "OPTIONS").
			Path(route.Pattern).
			Name(route.Name).
			Handler(middleware.Cover(route.HandlerFunc, middlewareChain...))
	}
	return r, nil
}

var (
	m          *router
	routerOnce sync.Once
)

func MuxRouter() IMuxRouter {
	if m == nil {
		routerOnce.Do(func() {
			m = &router{}
		})
	}
	return m
}
