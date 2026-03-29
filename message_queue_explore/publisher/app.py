import os

import pika
from fastapi import FastAPI, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from loguru import logger

from publisher.config import PublisherConfig
from publisher.entity import Message
from publisher.repository import RabbitMQRepository
from publisher.usecase import MessageUseCase


def create_app() -> FastAPI:
    """Create the FastAPI application with routes and dependencies."""
    config = PublisherConfig()

    connection_parameters = pika.URLParameters(config.rabbitmq_url)
    repository = RabbitMQRepository(connection_parameters, config.exchange_name)
    usecase = MessageUseCase(repository)

    app = FastAPI(title="Message Queue Publisher")
    app.add_middleware(CORSMiddleware, allow_origins=["*"], allow_methods=["*"], allow_headers=["*"])

    @app.post("/publish")
    def publish_message(payload: Message) -> dict:
        """HTTP endpoint for publishing notification messages."""
        try:
            usecase.publish(payload)
        except ValueError as err:
            raise HTTPException(status_code=400, detail=str(err))
        except pika.exceptions.AMQPError as err:
            logger.error("RabbitMQ error while publishing: %s", err)
            raise HTTPException(status_code=500, detail="failed to publish message")

        return {"code": 200, "message": "Message published successfully"}

    return app


if __name__ == "__main__":
    app = create_app()
    port = int(os.getenv("PORT", "8080"))
    logger.info("Starting publisher on port %s", port)
    import uvicorn

    uvicorn.run(app, host="0.0.0.0", port=port)
