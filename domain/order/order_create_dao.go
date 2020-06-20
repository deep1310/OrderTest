package order

import (
	"createorder/datarepository/mysql/order_db"
	"createorder/utils/errors"
	"github.com/jinzhu/gorm"
)

func (orderReq *CreateOrderReq) OrderSave() (*Order, *errors.RestErr) {
	db := order_db.GetSqlConn()
	createOrder := Order{}
	createOrder.UserId = orderReq.UserId
	createOrder.ProductId = orderReq.ProductId
	createOrder.Quantity = orderReq.Quantity
	createOrder.Fare = orderReq.Fare
	createOrder.Discount = orderReq.DiscountAmt
	createOrder.Status = "INITIATED"

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return nil, errors.InternalServerError("Not able to update order")
	}
	if result := tx.Model(&createOrder).Create(&createOrder); result != nil {
		if result.Error != nil {
			tx.Rollback()
			return nil, errors.InternalServerError("Not able to create order")
		}
	}
	paymentSessionData, err := PaymentSessionSave(tx, orderReq, createOrder.OrderId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if createOrder.Discount > 0 {
		paymentItemReq := PaymentItem{}
		paymentItemReq.Status = "INITIATED"
		paymentItemReq.Amount = createOrder.Discount
		paymentItemReq.PaymentSessionId = paymentSessionData.PaymentSessionId
		paymentItemReq.PaymentType = "DISCOUNT"
		paymentItemReq.PaymentMode = "DISCOUNT"
		itemErr := paymentItemReq.PaymentItemAdd(tx)
		if itemErr != nil {
			return nil, itemErr
		}

		sessionUpReq := PaymentSessionUpdateReq{}
		sessionUpReq.Status = paymentSessionData.Status
		sessionUpReq.AmountLeft = createOrder.Fare - createOrder.Discount
		sessionUpReq.AmountPaid = createOrder.Discount

		if sessionUpdateErr := sessionUpReq.PaymentSessionUpdate(tx); sessionUpdateErr != nil {
			return nil, sessionUpdateErr
		}
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, errors.InternalServerError("Not able to creat order")
	}

	return &createOrder, nil
}

func (updateReq *OrderUpdateRequest) OrderUpdate(tx *gorm.DB) *errors.RestErr {
	db := order_db.GetSqlConn()
	orderData := Order{}
	totalRowsUpdated := db.Model(&orderData).Where("orderId = ?", updateReq.OrderId).Update("status", updateReq.Status).RowsAffected
	if totalRowsUpdated == 0 {
		tx.Rollback()
		return errors.InternalServerError("Not able to update order")
	}
	return nil
}

func PaymentSessionSave(tx *gorm.DB, orderReq *CreateOrderReq, orderId int64) (*PaymentSession, *errors.RestErr) {

	paymentSessionReq := PaymentSession{}
	paymentSessionReq.Status = "INITIATED"
	paymentSessionReq.OrderAmount = orderReq.Fare
	paymentSessionReq.AmountPaid = 0.0
	paymentSessionReq.AmountLeft = orderReq.Fare
	paymentSessionReq.OrderId = orderId

	if result := tx.Model(&paymentSessionReq).Create(&paymentSessionReq); result != nil {
		if result.Error != nil {
			return nil, errors.InternalServerError("Not able to create payment session")
		}
	}
	return &paymentSessionReq, nil

}

func (paymentItemReq *PaymentItem) PaymentItemAdd(tx *gorm.DB) *errors.RestErr {
	if result := tx.Create(&paymentItemReq); result != nil {
		if result.Error != nil {
			tx.Rollback()
			return errors.InternalServerError("Not able to create payment item")
		}
	}
	return nil
}

func (req *PaymentSessionUpdateReq) PaymentSessionUpdate(tx *gorm.DB) *errors.RestErr {
	sessionModel := PaymentSession{}
	totalRowsUpdated := tx.Model(&sessionModel).Where("id = ?", req.PaymentSessionId).Updates(map[string]interface{}{
		"status":      req.Status,
		"amount_left": req.AmountLeft,
		"amount_paid": req.AmountPaid}).RowsAffected
	if totalRowsUpdated == 0 {
		tx.Rollback()
		return errors.InternalServerError("Not able to payment session")
	}

	return nil
}

func (req *PaymentItemUpdateReq) PaymentItemUpdate(tx *gorm.DB) (*PaymentItem, *errors.RestErr) {
	pt := PaymentItem{}
	totalRowsUpdated := tx.Model(&pt).Where("id = ?", req.PaymentItemId).Updates(map[string]interface{}{
		"status": req.Status}).RowsAffected
	if totalRowsUpdated == 0 {
		tx.Rollback()
		return nil, errors.InternalServerError("Not able to payment item")
	}

	return &pt, nil
}

func (o *OrderResp) GetOrder() *errors.RestErr {
	db := order_db.GetSqlConn()

	err := db.Where("id = ?", o.OrderItem.OrderId).Find(&o.OrderItem).Error
	if err != nil {
		return errors.InternalServerError("unable to get the order details")
	}
	err = db.Where("id = ?", o.OrderItem.OrderId).Find(&o.PaymentSessionData).Error
	if err != nil {
		return errors.InternalServerError("unable to get the order details")
	}
	err = db.Where("id = ?", o.PaymentSessionData.PaymentSessionId).Find(&o.PaymentItemData).Error
	if err != nil {
		return errors.InternalServerError("unable to get the order details")
	}

	return nil
}

func (pt *PaymentItemUpdateReq) UpdatePaymentItem() *errors.RestErr {
	db := order_db.GetSqlConn()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return errors.InternalServerError("Not able to update payment item")
	}

	paymentItem, err := pt.PaymentItemUpdate(tx)
	if err != nil {
		return err
	}

	if pt.Status == "COMPLETED" {
		sessionUpReq := PaymentSessionUpdateReq{}
		sessionUpReq.Status = pt.Status

		ps := &PaymentSession{
			PaymentSessionId: pt.PaymentSessionId,
		}
		if err := ps.GetPaymentSession(); err != nil {
			tx.Rollback()
			return err
		}
		sessionUpReq.AmountLeft = ps.AmountLeft - paymentItem.Amount
		sessionUpReq.AmountPaid = ps.AmountPaid + paymentItem.Amount

		if sessionUpdateErr := sessionUpReq.PaymentSessionUpdate(tx); sessionUpdateErr != nil {
			return sessionUpdateErr
		}
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return errors.InternalServerError("Not able to update payment item")
	}

	return nil
}

func (pt *PaymentItem) AddPaymentItem() *errors.RestErr {
	db := order_db.GetSqlConn()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return errors.InternalServerError("Not able to add payment item")
	}

	if err := pt.PaymentItemAdd(tx); err != nil {
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return errors.InternalServerError("Not able to add payment item")
	}

	return nil
}

func (ps *PaymentSession) GetPaymentSession() *errors.RestErr {

	db := order_db.GetSqlConn()
	err := db.Where(map[string]interface{}{"id": ps.PaymentSessionId}).Find(&ps).Error
	if err != nil {
		return errors.InternalServerError("unable to get payment session")
	}
	return nil
}

func (o *OrderResp) OrderComplete() *errors.RestErr {

	db := order_db.GetSqlConn()

	orderUpdateReq := OrderUpdateRequest{}
	orderUpdateReq.Status = "COMPLETED"
	orderUpdateReq.OrderId = o.OrderItem.OrderId

	if err := o.GetOrder(); err != nil {
		return err
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return errors.InternalServerError("Not able to update order")
	}
	if o.PaymentSessionData.AmountLeft == 0 {
		psUpdateReq := &PaymentSessionUpdateReq{
			Status:     "COMPLETED",
			AmountPaid: o.PaymentSessionData.AmountPaid,
			AmountLeft: o.PaymentSessionData.AmountLeft,
		}
		if err := psUpdateReq.PaymentSessionUpdate(tx); err != nil {
			return err
		}
		if err := orderUpdateReq.OrderUpdate(tx); err != nil {
			return errors.InternalServerError("Not able to update order")
		}
		if err := o.GetOrder(); err != nil {
			tx.Rollback()
			return errors.InternalServerError("Not able to update order")
		}
		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return errors.InternalServerError("Not able to update order")
		}
	} else {
		return errors.InternalServerError("All amount is not paid")
	}
	return nil
}
