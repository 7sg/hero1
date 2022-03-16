package main

import (
	"hero1/go-services/historical-events/internal/config"
	"hero1/go-services/historical-events/internal/infra/db/mongo"
	"hero1/go-services/historical-events/internal/infra/server"
	"hero1/go-services/historical-events/internal/infra/server/grpc"
	"hero1/go-services/historical-events/internal/infra/server/http"
	eventRepo "hero1/go-services/historical-events/internal/repository/event"
	"hero1/go-services/historical-events/internal/service/event"
	"log"
	"sync"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed on loading config, err:%+v", err)
	}

	dbClient, err := db.NewClient(cfg)
	if err != nil {
		log.Fatalf("failed on creating db client, err:%+v", err)
	}

	//repository
	eventRepository := eventRepo.New(dbClient, cfg.DB)

	//service
	eventService := event.New(eventRepository)

	grpcClient := grpc.New(cfg, eventService)
	httpClient := http.New(cfg, eventService)

	eventServer := server.New(grpcClient, httpClient)
	go eventServer.StartGRPC()
	go eventServer.StartHTTP()

	// Block forever
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
