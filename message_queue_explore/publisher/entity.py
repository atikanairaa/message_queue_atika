from pydantic import BaseModel


class Message(BaseModel):
    """Domain model for a notification message."""

    order_id: str
    user_id: str
    content: str
    timestamp: str
