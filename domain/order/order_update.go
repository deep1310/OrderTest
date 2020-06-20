package order

type OrderUpdateRequest struct {
	OrderId int64  `json:"orderId"`
	Status  string `json:"status"`
}

