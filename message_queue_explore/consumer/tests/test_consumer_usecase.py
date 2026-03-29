import pytest

from consumer.entity import Message
from consumer.usecase import MessageUseCase


class DummyRepo:
    def __init__(self):
        self.consumed = False

    def consume(self, service_name, handler):
        self.consumed = True
        msg = Message(order_id="1", user_id="u1", content="hi", timestamp="2025-01-01T00:00:00Z")
        handler(msg)


def test_consume_email_executes_repo_and_handler(monkeypatch):
    dummy = DummyRepo()
    usecase = MessageUseCase(dummy)

    called = {"email": False}

    def fake_send_email(message):
        called["email"] = True

    monkeypatch.setattr("consumer.usecase.send_email", fake_send_email)
    monkeypatch.setattr("consumer.usecase.send_sms", lambda x: None)
    monkeypatch.setattr("consumer.usecase.send_fcm", lambda x: None)

    usecase.consume("EMAIL")

    assert dummy.consumed
    assert called["email"]


def test_consume_unknown_service_raises():
    dummy = DummyRepo()
    usecase = MessageUseCase(dummy)

    with pytest.raises(ValueError):
        usecase.consume("UNKNOWN")
