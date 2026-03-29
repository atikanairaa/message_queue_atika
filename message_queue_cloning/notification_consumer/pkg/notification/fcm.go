package notification

import (
	"fmt"
	"notification_consumer/internal/entity"
)

func SendFcm(message entity.Message) {
	fmt.Println("sending fcm ... ", message.Content)
}
