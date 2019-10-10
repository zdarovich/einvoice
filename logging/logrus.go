package logging

import (
	"net/http"
	"os"
	"time"

	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
)

//Init ...
func Init(isLogFile bool) {
	if isLogFile {
		f, err := os.OpenFile("einvoice.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			logrus.Error(err)
			logrus.SetOutput(colorable.NewColorableStdout())
		} else {
			logrus.SetOutput(f)
		}
	} else {
		logrus.SetOutput(colorable.NewColorableStdout())

	}
	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.StampMilli,
	})
}

//HTTP level log like in previous logging
func HTTP(req *http.Request, res *http.Response, err error, duration time.Duration, method string) {

	logrus.Infof(
		"Request %s %s %s",
		req.Method,
		method,
		req.URL.String(),
	)
	if err != nil {
		logrus.Error(err)
		return
	}

	duration /= time.Millisecond
	logrus.Infof(
		"Response status=%d durationMs=%d",
		res.StatusCode,
		duration,
	)
}
