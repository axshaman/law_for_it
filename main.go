package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"law_for_it/internal/api"
	"law_for_it/internal/config"
	"law_for_it/internal/document"
	"law_for_it/internal/payment"
	"law_for_it/internal/state"
	"law_for_it/internal/storage"
)

var version = "0.1.0"

func main() {
	configPath := flag.String("config", "", "path to configuration file (JSON)")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	store, err := storage.New(cfg.Database.Path)
	if err != nil {
		log.Fatalf("failed to initialise storage: %v", err)
	}

	manager, err := state.NewManager(store, state.New())
	if err != nil {
		log.Fatalf("failed to initialise application state: %v", err)
	}

	paymentSvc, err := payment.NewService(manager, cfg.Payment.Currency)
	if err != nil {
		log.Fatalf("failed to init payment service: %v", err)
	}

	docSvc := document.NewService(manager, cfg.Document, paymentSvc)

	server := api.New(docSvc, paymentSvc, api.ServerInfo{
		Name:    "law_for_it",
		Version: version,
	})

	address := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Starting server on %s", address)
	if err := http.ListenAndServe(address, server); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
