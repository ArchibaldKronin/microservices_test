package model

type PaymentMethod string

const (
	PaymentMethodUNKNOWN       PaymentMethod = "UNKNOWN"
	PaymentMethodCARD          PaymentMethod = "CARD"
	PaymentMethodSBP           PaymentMethod = "SBP"
	PaymentMethodCREDITCARD    PaymentMethod = "CREDIT_CARD"
	PaymentMethodINVESTORMONEY PaymentMethod = "INVESTOR_MONEY"
)

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
