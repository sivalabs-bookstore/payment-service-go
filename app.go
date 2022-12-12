package main

import (
	"database/sql"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/sivalabs-bookstore/payment-service-go/config"
	"github.com/sivalabs-bookstore/payment-service-go/database"
	"github.com/sivalabs-bookstore/payment-service-go/payments"
	"net/http"
	"os"
)

type App struct {
	Router            *mux.Router
	db                *sql.DB
	paymentController *payments.PaymentController
}

func NewApp(config config.AppConfig) *App {
	app := &App{}
	app.init(config)
	return app
}

func (app *App) init(config config.AppConfig) {
	//logFile := initLogging()
	//defer logFile.Close()
	app.initLogging()

	app.db = database.GetDb(config)

	creditCardRepo := payments.NewCreditCardRepo(app.db)
	app.paymentController = payments.NewPaymentController(creditCardRepo)

	app.Router = app.setupRoutes()
}

func (app *App) setupRoutes() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/api/payments/validate", app.paymentController.ValidatePayment).Methods(http.MethodPost)
	return router
}

func (app *App) initLogging() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{})
	log.SetLevel(log.InfoLevel)
}

func (app *App) initFileLogging() *os.File {
	logFile, err := os.OpenFile("payments.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening log file: %v", err)
	}
	log.SetOutput(logFile)
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)
	return logFile
}
