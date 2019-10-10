package main

import (
	"github.com/braintree/manners"
	"github.com/gorilla/handlers"
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
	"github.com/zdarovich/einvoice/app"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)

func main() {
	c, err := app.ServiceContainer().InjectConfiguration()
	if err != nil {
		logrus.Error(err)
		return
	}

	if c.Server.LogFileEnable {
		f, err := os.OpenFile(c.Server.LogFileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
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

	db, err := app.ServiceContainer().InjectDatabase(c)
	if err != nil {
		logrus.Error(err)
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		logrus.Error(err)
		return
	}

	r, err := app.MuxRouter().InitRouter(c, db)
	if err != nil {
		logrus.Error(err)
		return
	}

	headers := handlers.AllowedHeaders([]string{
		"Origin",
		"Access-Control-Request-Method",
		"Access-Control-Request-Headers",
		"X-Requested-With",
		"Content-Type"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"*"})

	server := &http.Server{
		ReadHeaderTimeout: time.Duration(c.Server.ReadHeaderTimeoutSec) * time.Second,
		ReadTimeout:       time.Duration(c.Server.ReadTimeoutSec) * time.Second,
		WriteTimeout:      time.Duration(c.Server.WriteTimeoutSec) * time.Second,
		Addr:              ":" + strconv.Itoa(c.Server.Port),
		Handler:           handlers.CORS(headers, methods, origins)(r),
	}
	logrus.Info("server started")
	err = manners.NewWithServer(server).ListenAndServe()
	if err != nil {
		logrus.Error(err)
		return
	}
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, os.Kill)
	go listenForShutdown(ch)
}

func listenForShutdown(ch <-chan os.Signal) {
	<-ch
	manners.Close()
}
