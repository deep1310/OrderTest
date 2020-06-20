package services

import (
	"createorder/domain/order"
	"createorder/utils/errors"
)

type PaymentSessionInterface interface {
	CreatePaymentItem(*order.PaymentItem) *errors.RestErr
	UpdatePaymentRecord(req *order.PaymentItemUpdateReq) *errors.RestErr
}

type createPaymentSessionRepo struct{}

func PaymentSessionService() PaymentSessionInterface {
	return &createPaymentSessionRepo{}
}

func (oService *createPaymentSessionRepo) CreatePaymentItem(paymentReq *order.PaymentItem) *errors.RestErr {
	return paymentReq.AddPaymentItem()
}

func (oService *createPaymentSessionRepo) UpdatePaymentRecord(paymentReq *order.PaymentItemUpdateReq) *errors.RestErr {
	return paymentReq.UpdatePaymentItem()
}
