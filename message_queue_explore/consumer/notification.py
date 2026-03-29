from consumer.entity import Message
from loguru import logger


def send_email(message: Message) -> None:
    """Simulated email sending action for notification message."""
    logger.info("sending email %s", message.content)


def send_sms(message: Message) -> None:
    """Simulated SMS sending action for notification message."""
    logger.info("sending sms %s", message.content)


def send_fcm(message: Message) -> None:
    """Simulated FCM sending action for notification message."""
    logger.info("sending fcm %s", message.content)
