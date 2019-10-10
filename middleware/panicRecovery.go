package middleware

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/zdarovich/einvoice/models"
	"net/http"
	"runtime/debug"
)

//Recovery ...
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if err := recover(); err != nil {
				logrus.Error(err)
				logrus.Error(string(debug.Stack()))
				serverErr(w, "There was an internal server error. Please contact our support.", http.StatusInternalServerError)
				return
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func serverErr(w http.ResponseWriter, msg string, status int) {
	r := &models.ServerErrorResponse{
		Message:    msg,
		StatusCode: status,
	}

	rJSON, e := json.Marshal(r)
	if e != nil {
		logrus.Error("Can not encode json error response", e.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(rJSON)
}
