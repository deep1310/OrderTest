package order

import (
	"createorder/utils/errors"
	"github.com/jinzhu/gorm"
	"strings"
	"time"
)

type Order struct {
	OrderId   int64     `gorm:"primary_key ;column:id" ;json:"orderId"`
	UserId    int64     `gorm:"not null ;column:user_id" json:"userId"`
	Quantity  int       `gorm:"column:quantity" json:"quantity"`
	ProductId int64     `gorm:"column:product_id" json:"productId"`
	Fare      float64   `gorm:"column:fare" json:"fare"`
	Discount  float64   `gorm:"column:discount"json:"discount"`
	Status    string    `gorm:"column:status" json:"status"`
	CreatedAt time.Time `gorm:"column:created_at" json:"-"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"-"`
}

type PaymentSession struct {
	PaymentSessionId int64     `gorm:"primary_key ;column:id" ;json:"paymentSessionId"`
	OrderId          int64     `gorm:"column:order_id" ;json:"orderId"`
	OrderAmount      float64   `gorm:"column:order_amount" ;json:"orderAmount"`
	AmountPaid       float64   `gorm:"column:amount_paid"  ;json:"amountPaid"`
	AmountLeft       float64   `gorm:"column:amount_left" ;json:"amountLeft"`
	Status           string    `gorm:"column:status"  ;json:"status"`
	CreatedAt        time.Time `gorm:"column:created_at";json:"-"`
	UpdatedAt        time.Time `gorm:"column:updated_at";json:"-"`
}

type PaymentItem struct {
	PaymentItemId    int64     `gorm:"primary_key ;column:id";json:"paymentItemId"`
	PaymentSessionId int64     `gorm:"column:payment_session_id";json:"paymentSessionId"`
	Amount           float64   `gorm:"column:amount";json:"amount"`
	PaymentType      string    `gorm:"column:payment_type";json:"paymentType"`
	PaymentMode      string    `gorm:"column:payment_mode";json:"paymentMode"`
	Status           string    `gorm:"column:status";json:"status"`
	CreatedAt        time.Time `gorm:"column:created_at";json:"-"`
	UpdatedAt        time.Time `gorm:"column:updated_at";json:"-"`
}

func (Order) TableName() string {
	return "order_item"
}

func (PaymentSession) TableName() string {
	return "payment_session"
}

func (order *Order) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("created_at", time.Now().UTC())
	scope.SetColumn("updated_at", time.Now().UTC())
	return nil
}

func (order *Order) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("updated_at", time.Now().UTC())
	return nil
}

func (order *PaymentSession) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("created_at", time.Now().UTC())
	scope.SetColumn("updated_at", time.Now().UTC())
	return nil
}

func (order *PaymentSession) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("updated_at", time.Now().UTC())
	return nil
}

func (order *PaymentItem) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("created_at", time.Now().UTC())
	scope.SetColumn("updated_at", time.Now().UTC())
	return nil
}

func (order *PaymentItem) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("updated_at", time.Now().UTC())
	return nil
}

func (order *CreateOrderReq) ValidateOrderRequest() *errors.RestErr {

	if (order.ProductId) <= 0 {
		if err := errors.BadRequestError("product id is empty"); err != nil {
			return err
		}
	}
	if order.Quantity == 0 {
		if err := errors.BadRequestError("quantity cannot be 0"); err != nil {
			return err
		}
	}

	return nil
}

func (order *OrderUpdateRequest) ValidateUpdateOrderRequest() *errors.RestErr {

	order.Status = strings.TrimSpace(strings.ToUpper(order.Status))
	if order.OrderId <= 0 {
		if err := errors.BadRequestError("order id cannot be empty"); err != nil {
			return err
		}
	}

	if order.Status == "" {
		if err := errors.BadRequestError("status cannot be empty"); err != nil {
			return err
		}
	}

	if !(order.Status == "FAILED" || order.Status == "COMPLETED") {
		if err := errors.BadRequestError("status value is wrong"); err != nil {
			return err
		}
	}
	return nil
}
