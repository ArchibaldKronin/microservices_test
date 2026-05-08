package model

type OrderPaidEvent struct {
	EventUuid       string
	OrderUuid       string
	UserUuid        string
	PaymentMethod   PaymentMethod
	TransactionUuid string
}

type ShipAssembledEvent struct {
	EventUuid    string
	OrderUuid    string
	UserUuid     string
	BuildTimeSec int64
}
