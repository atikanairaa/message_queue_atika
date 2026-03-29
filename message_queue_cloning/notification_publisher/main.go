package main

import (
	"log"
	"notification_publisher/config"

	"notification_publisher/internal/handler"
	"notification_publisher/internal/repository"
	"notification_publisher/internal/usecase"
	"notification_publisher/pkg/rabbitmq"

	"github.com/labstack/echo/v4"
)

func main() {
	// Load configuration
	var cfg config.Config
	if err := config.LoadConfig(&cfg); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	e := echo.New()

	// Setup RabbitMQ connection
	rmqConn, err := rabbitmq.NewRabbitMQConnection(cfg.RabbitMQURL)
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %v", err)
	}
	defer rmqConn.Close()

	repo := repository.NewRabbitMQRepository(rmqConn, cfg.ExchangeName)
	useCase := usecase.NewMessageUseCase(repo)
	h := handler.NewMessageHandler(useCase)

	e.POST("/publish", h.PublishMessage)

	e.Logger.Fatal(e.Start(":" + cfg.Port))
}
