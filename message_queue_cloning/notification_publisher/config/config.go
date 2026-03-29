package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port         string `envconfig:"PORT" default:"8080"`
	RabbitMQURL  string `envconfig:"RABBITMQ_URL" default:"amqp://guest:guest@localhost:5672/"`
	ExchangeName string `envconfig:"EXCHANGE_NAME" default:"notifications"`
}

func LoadConfig(cfg *Config) error {
	if err := envconfig.Process("ecommerce", cfg); err != nil {
		log.Fatalf("failed to load config: %v", err)
		return err
	}
	log.Printf("Using Exchange: %s", cfg.ExchangeName)
	return nil
}
