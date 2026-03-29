package main

import (
	"context"
	"log"
	"notification_consumer/config"
	"notification_consumer/internal/repository"
	"notification_consumer/internal/usecase"
	"notification_consumer/pkg/rabbitmq"
)

func main() {
	// Load config
	var cfg config.Config
	if err := config.LoadConfig(&cfg); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Setup RabbitMQ connection
	rmqConn, err := rabbitmq.NewRabbitMQConnection(cfg.RabbitMQURL)
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %v", err)
	}
	defer rmqConn.Close()

	// Ambil nama service dari ENV atau default "Consumer"
	// serviceName := os.Getenv("SERVICE_NAME")
	// if serviceName == "" {
	// 	serviceName = "Consumer"
	// }

	serviceName := "SMS"

	log.Printf("[%s] Listening for messages...", serviceName)
	repo := repository.NewRabbitMQRepository(rmqConn, cfg.ExchangeName)
	useCase := usecase.NewMessageUseCase(repo)

	// Konsumsi pesan dari RabbitMQ
	if err := useCase.ConsumeMessagesSms(context.Background(), serviceName); err != nil {
		log.Fatalf("[%s] Failed to consume messages: %v", serviceName, err)
	}

	// Blok agar tidak keluar
	select {}
}
