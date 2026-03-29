from consumer.notification import send_email, send_fcm, send_sms
from consumer.repository import RabbitMQRepository


class MessageUseCase:
    """Use case for routing queue messages to notification senders."""

    def __init__(self, repository: RabbitMQRepository) -> None:
        self._repository = repository

    def consume(self, service_name: str) -> None:
        """Start consuming messages for given service name."""

        def _handler(message):
            if service_name.upper() == "EMAIL":
                send_email(message)
            elif service_name.upper() == "SMS":
                send_sms(message)
            elif service_name.upper() == "FCM":
                send_fcm(message)
            else:
                raise ValueError(f"Unsupported service type: {service_name}")

        self._repository.consume(service_name, _handler)
