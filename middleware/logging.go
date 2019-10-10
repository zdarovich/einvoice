package middleware

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

//Logging ...
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		log.Infof("SERVER: new %s request", r.RequestURI)
		start := time.Now()
		next.ServeHTTP(w, r)
		elapsed := time.Since(start)

		log.Infof("SERVER: %s request, elapsed: %s", r.RequestURI, elapsed.String())
	})
}
