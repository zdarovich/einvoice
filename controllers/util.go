package controllers

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/zdarovich/einvoice/models"
	"net/http"
)

func errorHandler(w http.ResponseWriter, err error) {
	logrus.Error(err)
	resp := new(models.ServerErrorResponse)
	resp.Message = err.Error()
	resp.StatusCode = http.StatusInternalServerError

	respJSON, err := json.Marshal(resp)
	if err != nil {
		logrus.Error("Can not encode json error response")
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(respJSON)
}

func response(w http.ResponseWriter, model interface{}) {
	respJSON, err := json.Marshal(model)
	if err != nil {
		logrus.Error("Can not encode json error response")
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respJSON)
}
