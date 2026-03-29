# Analisis Arsitektur dan Kode

## 1. Ringkasan proyek

Proyek ini adalah implementasi contoh pattern pub/sub menggunakan RabbitMQ dengan dua implementasi:
- `message_queue_cloning/notification_publisher`: Publisher HTTP Go + RabbitMQ
- `message_queue_cloning/notification_consumer`: Konsumer Go (EMAIL, SMS, FCM)

Kedua modul menggunakan exchange `fanout` untuk mendistribusikan payload pesan ke semua antrean terkait.

## 2. Struktur high-level

### notification_publisher
- `main.go`: entrypoint. load config, koneksi RabbitMQ, inject dependensi, register endpoint POST `/publish`.
- `config/config.go`: config struct (PORT, RABBITMQ_URL, EXCHANGE_NAME) via `envconfig`.
- `pkg/rabbitmq/connection.go`: koneksi AMQP standar dengan `amqp.Dial`.
- `internal/repository/rabbitmq_repo.go`: deklarasi exchange + publish JSON ke exchange dengan `PublishWithContext`.
- `internal/usecase/message_usecase.go`: lapisan domain yang memanggil repository.
- `internal/handler/handler.go`: HTTP binding request body ke `entity.Message`, validasi, respond JSON.
- `internal/entity/entity.go`: model message.

### notification_consumer
- `cmd/email/main.go`, `cmd/sms/main.go`, `cmd/fcm/main.go`: masing-masing service khusus membaca config, koneksi RMQ, instantiate repo/usecase, panggil metode konsumsi, `select {}` lockmain.
- `config/config.go`: sama (PORT optional, RABBITMQ_URL, EXCHANGE_NAME) via envconfig.
- `pkg/rabbitmq/connection.go`: koneksi amqp.
- `internal/repository/rabbitmq_repo.go`: deklarasi exchange fanout, queue declare (berbasis serviceName), bind queue ke exchange, consume.
    - Konsumsi non-auto-ack (`false`), manual `msg.Ack(false)` saat sukses.
    - Loop goroutine untuk pesan; unmarshal JSON, log, kirim ke `pkg/notification` handler.
- `internal/usecase/message_usecase.go`: lapisan usecase memanggil repos.
- `internal/entity/entity.go`: model message identik.
- `pkg/notification/{email,sms,fcm}.go`: stub fungsi output `fmt.Println`.

## 3. Alur eksekusi

### 3.1 Publisher
1. Start dengan `go run main.go`.
2. `LoadConfig` baca env: `PORT` default 8080, `RABBITMQ_URL` default `amqp://guest:guest@localhost:5672/`, `EXCHANGE_NAME` default `notifications`.
3. Buat koneksi RabbitMQ.
4. Endpoint `POST /publish` -> `MessageHandler.PublishMessage`:
   - bind JSON ke `Message`
   - usecase `PublishMessage(ctx, "notifications", message)`
   - repository `PublishMessage` -> `ExchangeDeclare(fanout)` dan `PublishWithContext`
5. Kembali JSON success/failed.

### 3.2 Consumer (EMAIL/SMS/FCM)
1. Start service masing-masing dari folder `cmd/{email,sms,fcm}`.
2. Load config + RabbitMQ connection.
3. `useCase.ConsumeMessagesX(context.Background(), serviceName)` -> repository:
   - `ExchangeDeclare(fanout)` (idempotent) 
   - `QueueDeclare(serviceName, durable true, ... )`
   - `QueueBind(queue, "", exchange,...)`
   - `Consume(queue, serviceName, autoAck false, ... )`
4. Goroutine konsumsi:
   - `json.Unmarshal(msg.Body, &entity.Message)`
   - log message
   - `notification.SendX(message)` -> print ke stdout
   - `msg.Ack(false)`

## 4. Pola arsitektur

- Clean architecture (layered): config -> infrastructure (RabbitMQ) -> repository -> usecase -> handler.
- Dependency injection secara manual di `main.go`.
- Exchange fanout + queue dedicated setiap kanal konsumer.
- Non-blocking pesanan (goroutine consumer async) dan ack manual -> safe ketika error.

## 5. Kelebihan dan perbaikan potensial

### Kelebihan
- Mudah diikuti (clear separation concern)
- Sederhana untuk edukasi
- Konsumen bisa scale dengan menjalankan multiple instance/queue

### Perbaikan
- Tambah retry/dlq bila gagal unmarshal atau pemrosesan.
- Konfigurasi queue/exchange hard-coded `fanout` vs env.
- Kirim `serviceName` dynamic via env `SERVICE_NAME` (baris dikomentari di sms/fcm/email) agar lebih fleksibel.
- Tambah test unit/integration untuk usecase/repository.
- Kelola graceful shutdown dan context cancel ketika service stop.

## 6. Diagram singkat
```
[HTTP Client] -> /publish -> [Publisher (Go)] -> RabbitMQ exchange "notifications" (fanout)
                                                         /            |           \
                                                        /             |            \
                                                       v              v             v
                                             [Email queue EMAIL] [SMS queue SMS] [FCM queue FCM]
                                                       |             |             |
                                                       v             v             v
                                                 [Email consumer] [SMS consumer] [FCM consumer]
```