import os


class ConsumerConfig:
    """Application configuration for consumer service loaded from environment variables."""

    def __init__(self) -> None:
        self.rabbitmq_url = os.getenv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/%2F")
        self.exchange_name = os.getenv("EXCHANGE_NAME", "notifications")
        self.service_name = os.getenv("SERVICE_NAME", "CONSUMER")

        self._validate()

    def _validate(self) -> None:
        if not self.rabbitmq_url:
            raise ValueError("RABBITMQ_URL must be configured")
        if not self.exchange_name:
            raise ValueError("EXCHANGE_NAME must be configured")
        if not self.service_name:
            raise ValueError("SERVICE_NAME must be configured")
