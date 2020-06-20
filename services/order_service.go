package services

import (
	"createorder/domain/order"
	"createorder/utils/errors"
)

type OrderInterface interface {
	CreateOrder(*order.CreateOrderReq) (*order.Order, *errors.RestErr)
	GetOrder(int64) (*order.OrderResp, *errors.RestErr)
	CompleteOrder(int64) (*order.OrderResp, *errors.RestErr)
}

type createOrderRepo struct{}

func OrderService() OrderInterface {
	return &createOrderRepo{}
}

func (oService *createOrderRepo) CreateOrder(orderReq *order.CreateOrderReq) (*order.Order, *errors.RestErr) {

	if err := orderReq.ValidateOrderRequest(); err != nil {
		return nil, err
	}
	orderData, restErr := orderReq.OrderSave()
	if restErr != nil {
		return nil, restErr
	}

	return orderData, nil
}

func (oService *createOrderRepo) GetOrder(orderId int64) (*order.OrderResp, *errors.RestErr) {

	orderResp := &order.OrderResp{}
	orderResp.OrderItem.OrderId = orderId

	if err := orderResp.GetOrder(); err != nil {
		return nil, err
	}
	return orderResp, nil
}

func (oService *createOrderRepo) CompleteOrder(orderId int64) (*order.OrderResp, *errors.RestErr) {
	orderComplete := &order.OrderResp{}
	orderComplete.OrderItem.OrderId = orderId

	if err := orderComplete.OrderComplete(); err != nil {
		return nil, err
	}

	return orderComplete, nil
}
