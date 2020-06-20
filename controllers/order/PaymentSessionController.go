package order

import (
	"createorder/domain/order"
	"createorder/services"
	"createorder/utils/errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AddPaymentItem(c *gin.Context) {

	var payReq order.PaymentItem
	if err := c.ShouldBindJSON(&payReq); err != nil {
		apiReqErr := errors.BadRequestError("invalid request")
		c.JSON(apiReqErr.Code, apiReqErr)
		return
	}

	payClient := services.PaymentSessionService()
	err := payClient.CreatePaymentItem(&payReq)
	if err != nil {
		c.JSON(err.Code, err)
		return
	}
	c.JSON(http.StatusCreated, "")
}

func UpdatePaymentItem(c *gin.Context) {

	var payReq order.PaymentItemUpdateReq
	if err := c.ShouldBindJSON(&payReq); err != nil {
		apiReqErr := errors.BadRequestError("invalid request")
		c.JSON(apiReqErr.Code, apiReqErr)
		return
	}

	payClient := services.PaymentSessionService()
	err := payClient.UpdatePaymentRecord(&payReq)
	if err != nil {
		c.JSON(err.Code, err)
		return
	}
	c.JSON(http.StatusOK, "")
}
