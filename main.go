package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/sivalabs-bookstore/payment-service-go/cmd"
	"github.com/sivalabs-bookstore/payment-service-go/internal/config"
	"net/http"
	"time"
)

func main() {
	cfg := config.GetConfig()
	app := cmd.NewApp(cfg)

	port := fmt.Sprintf(":%d", cfg.AppPort)
	srv := &http.Server{
		Handler:        app.Router,
		Addr:           port,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("listening on port %d", cfg.AppPort)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
