package order

type PaymentSessionUpdateReq struct {
	PaymentSessionId int64
	AmountPaid       float64
	AmountLeft       float64
	Status           string
}

type PaymentItemUpdateReq struct {
	PaymentItemId    int64  `json:"paymentItemId"`
	PaymentSessionId int64  `json:"paymentSessionId"`
	Status           string `json:"status"`
}
