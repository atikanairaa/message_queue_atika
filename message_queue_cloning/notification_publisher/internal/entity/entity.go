package entity

type Message struct {
	OrderID   string `json:"order_id"`
	UserID    string `json:"user_id"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}
