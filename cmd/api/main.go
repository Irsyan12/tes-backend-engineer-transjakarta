package main

import (
	"fmt"
	"log"

	"fleet-management/internal/config"
	"fleet-management/internal/database"
	"fleet-management/internal/handler"
	mqttsubscriber "fleet-management/internal/mqtt"
	"fleet-management/internal/rabbitmq"
	"fleet-management/internal/repository/postgres"
	"fleet-management/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {

	cfg := config.Load()

	db, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	log.Println("PostgreSQL connected successfully")

	// Setup Dependencies
	locationRepo := postgres.NewLocationRepository(db)
	locationService := service.NewLocationService(locationRepo)
	
	// RabbitMQ Producer
	rmqProducer, err := rabbitmq.NewProducer(cfg.RabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ Producer: %v", err)
	}
	defer rmqProducer.Close()

	// RabbitMQ Worker
	rmqWorker, err := rabbitmq.NewWorker(cfg.RabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ Worker: %v", err)
	}
	defer rmqWorker.Close()

	// Setup Handlers
	locationAPIHandler := handler.NewLocationAPIHandler(locationService) // location.go

	// START BACKGROUND SERVICES (Pendekatan 1 Aplikasi)
	
	// Mulai RabbitMQ Worker
	go rmqWorker.Start()

	// Mulai MQTT Subscriber
	subscriber := mqttsubscriber.NewSubscriber(cfg, locationService, rmqProducer)
	subscriber.Start()
	defer subscriber.Stop() // Disconnect rapi kalau server mati

	// START REST API
	app := gin.Default()

	app.GET("/health", handler.HealthCheck)
	
	// Endpoint Utama
	app.GET("/vehicles/:vehicle_id/location", locationAPIHandler.GetLatestLocation)
	app.GET("/vehicles/:vehicle_id/history", locationAPIHandler.GetHistory)

	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("Server starting on port %s", cfg.AppPort)
	if err := app.Run(addr); err != nil {
		log.Fatal(err)
	}
}