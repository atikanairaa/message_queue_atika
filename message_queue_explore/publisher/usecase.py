from publisher.entity import Message
from publisher.repository import RabbitMQRepository


class MessageUseCase:
    """Use case layer for publisher business rules."""

    def __init__(self, repository: RabbitMQRepository) -> None:
        self._repository = repository

    def publish(self, message: Message) -> None:
        """Publish message as business operation and enforce domain invariants."""
        if not message.content.strip():
            raise ValueError("message content cannot be empty")

        self._repository.publish_message(routing_key="", message=message)
