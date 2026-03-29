from pydantic import BaseModel


class Message(BaseModel):
    """Domain model for notifications consumed from queue."""

    order_id: str
    user_id: str
    content: str
    timestamp: str
