import json
from typing import Callable

import pika
from loguru import logger

from consumer.entity import Message


class RabbitMQRepository:
    """Repository for RabbitMQ consumer interactions."""

    def __init__(self, connection_parameters: pika.ConnectionParameters, exchange_name: str) -> None:
        self._connection_parameters = connection_parameters
        self._exchange_name = exchange_name

    def consume(self, service_name: str, handler: Callable[[Message], None]) -> None:
        """Subscribe to the fanout exchange and dispatch messages to handler."""
        logger.debug("Connecting RabbitMQ for consumer %s", service_name)

        connection = pika.BlockingConnection(self._connection_parameters)
        channel = connection.channel()

        channel.exchange_declare(exchange=self._exchange_name, exchange_type="fanout", durable=True)
        queue = channel.queue_declare(queue=service_name, durable=True)
        channel.queue_bind(queue=queue.method.queue, exchange=self._exchange_name)

        logger.info("[%s] Bound queue '%s' to exchange '%s'", service_name, queue.method.queue, self._exchange_name)
        logger.info("[%s] Waiting for messages...", service_name)

        for method_frame, properties, body in channel.consume(queue=queue.method.queue, inactivity_timeout=1):
            if body is None:
                continue
            try:
                message_dict = json.loads(body)
                message = Message(**message_dict)
            except Exception as err:
                logger.error("Failed to parse message: %s", err)
                channel.basic_ack(delivery_tag=method_frame.delivery_tag)
                continue

            try:
                logger.info("[%s] Received message: %s", service_name, message)
                handler(message)
                channel.basic_ack(delivery_tag=method_frame.delivery_tag)
            except Exception as err:
                logger.error("[%s] Handler error: %s", service_name, err)
                # do not ack to requeue on failure
        
        channel.close()
        connection.close()
