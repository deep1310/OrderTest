package order

type CreateOrderReq struct {
	Quantity  int     `json:"quantity"`
	ProductId int64   `json:"productId"`
	UserId    int64   `json:"userId"`
	Fare      float64 `json:"fare"`
	DiscountAmt      float64 `json:"discountAmt"`
}
