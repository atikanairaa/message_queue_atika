import pytest

from publisher.entity import Message
from publisher.usecase import MessageUseCase


class DummyRepo:
    def __init__(self):
        self.published = None

    def publish_message(self, routing_key, message):
        self.published = (routing_key, message)


def test_publish_success():
    dummy = DummyRepo()
    usecase = MessageUseCase(dummy)
    message = Message(order_id="1", user_id="u1", content="hello", timestamp="2025-01-01T00:00:00Z")

    usecase.publish(message)

    assert dummy.published is not None
    assert dummy.published[0] == ""
    assert dummy.published[1].content == "hello"


def test_publish_content_empty_raises_error():
    dummy = DummyRepo()
    usecase = MessageUseCase(dummy)
    message = Message(order_id="1", user_id="u1", content="   ", timestamp="2025-01-01T00:00:00Z")

    with pytest.raises(ValueError):
        usecase.publish(message)
