import argparse
import os
import sys

import pika
from loguru import logger

from consumer.config import ConsumerConfig
from consumer.repository import RabbitMQRepository
from consumer.usecase import MessageUseCase

# ensure logs are output to stdout with timestamps and service context
logger.remove()
logger.add(
    sink=sys.stdout,
    format="{time:YYYY-MM-DD HH:mm:ss} | {level} | {message}",
    level="INFO",
)


def main() -> None:
    """Entrypoint for consumer service."""
    parser = argparse.ArgumentParser(description="Consumer service for notifications")
    parser.add_argument("--service", choices=["EMAIL", "SMS", "FCM"], help="Service type to consume messages for")
    args = parser.parse_args()

    config = ConsumerConfig()
    service_name = args.service or os.getenv("SERVICE_NAME", "CONSUMER")

    connection_parameters = pika.URLParameters(config.rabbitmq_url)
    repository = RabbitMQRepository(connection_parameters, config.exchange_name)
    usecase = MessageUseCase(repository)

    logger.info("Starting consumer %s", service_name)

    try:
        usecase.consume(service_name)
    except Exception as err:
        logger.exception("Failed consumer: %s", err)
        raise


if __name__ == "__main__":
    main()
