import json
from typing import Any

import pika
from loguru import logger

from publisher.entity import Message


class RabbitMQRepository:
    """Repository for RabbitMQ message publishing."""

    def __init__(self, connection_parameters: pika.ConnectionParameters, exchange_name: str) -> None:
        self._connection_parameters = connection_parameters
        self._exchange_name = exchange_name

    def publish_message(self, routing_key: str, message: Message) -> None:
        """Publish a message to the configured fanout exchange."""
        logger.debug("Publishing message to exchange %s", self._exchange_name)
        connection = pika.BlockingConnection(self._connection_parameters)
        try:
            channel = connection.channel()
            channel.exchange_declare(
                exchange=self._exchange_name,
                exchange_type="fanout",
                durable=True,
            )

            payload = message.json().encode("utf-8")
            channel.basic_publish(
                exchange=self._exchange_name,
                routing_key=routing_key,
                body=payload,
                properties=pika.BasicProperties(content_type="application/json"),
            )
            logger.info("Message published: %s", message)
        finally:
            connection.close()
