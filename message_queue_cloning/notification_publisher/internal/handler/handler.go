package handler

import (
	"net/http"
	"notification_publisher/internal/entity"
	"notification_publisher/internal/usecase"

	"github.com/labstack/echo/v4"
)

type MessageHandler struct {
	usecase *usecase.MessageUseCase
}

func NewMessageHandler(usecase *usecase.MessageUseCase) *MessageHandler {
	return &MessageHandler{usecase: usecase}
}

func (h *MessageHandler) PublishMessage(c echo.Context) error {
	var message entity.Message
	if err := c.Bind(&message); err != nil {
		return c.JSON(http.StatusBadRequest, entity.Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid request",
		})
	}

	if err := h.usecase.PublishMessage(c.Request().Context(), "notifications", message); err != nil {
		return c.JSON(http.StatusInternalServerError, entity.Response{
			Code:    http.StatusInternalServerError,
			Message: "Failed to publish message",
		})
	}

	return c.JSON(http.StatusOK, entity.Response{
		Code:    http.StatusOK,
		Message: "Message published successfully",
	})
}
