package notification

import (
	"fmt"
	"notification_consumer/internal/entity"
)

func SendEmail(message entity.Message) {
	fmt.Println("sending email ...", message.Content)
}
