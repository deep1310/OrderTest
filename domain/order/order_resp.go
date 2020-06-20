package order

type OrderResp struct {
	OrderItem          Order          `json:"order_item"`
	PaymentSessionData PaymentSession `json:"paymentSession"`
	PaymentItemData    []PaymentItem  `json:"paymentItems"`
}
