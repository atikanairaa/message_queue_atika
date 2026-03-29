# README - Notification Consumer

## 1. Clone Aplikasi Consumer

```sh
git clone https://github.com/kusnadi8605/notification_consumer
cd notification_consumer
go mod init notification_consumer
go mod tidy
go mod vendor
```

## 2. Menjalankan Aplikasi Consumer

```sh
go run cmd/consumer/email/main.go
go run cmd/consumer/sms/main.go
go run cmd/consumer/fcm/main.go
```

## 3. Respon yang Diharapkan Ketika Ada yang Publish Message

### Email Consumer
```sh
Received message: {OrderID:12345 UserID:67890 Content:New order received Timestamp:2025-03-11T10:00:00Z}
sending email New order received
```

### SMS Consumer
```sh
Received message: {OrderID:12345 UserID:67890 Content:New order received Timestamp:2025-03-11T10:00:00Z}
sending SMS New order received
```

### FCM Consumer
```sh
Received message: {OrderID:12345 UserID:67890 Content:New order received Timestamp:2025-03-11T10:00:00Z}
sending FCM New order received
```

