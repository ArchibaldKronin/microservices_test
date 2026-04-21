package model

type Order struct {
	OrderID       string         `db:"order_id"`
	UserID        string         `db:"user_id"`
	PartIDs       []string       `db:"part_ids"`
	TotalPrice    float64        `db:"total_price"`
	TransactionID *string        `db:"transaction_id"`
	PaymentMethod *PaymentMethod `db:"payment_method"`
	Status        OrderStatus    `db:"status"`
}
