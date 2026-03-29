package notification

import (
	"fmt"
	"notification_consumer/internal/entity"
)

func SendSms(message entity.Message) {
	fmt.Println("sending sms ...", message.Content)
}
