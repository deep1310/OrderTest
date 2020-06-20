package order

import (
	"createorder/domain/order"
	"createorder/services"
	"createorder/utils/errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

func CreateOrder(c *gin.Context) {

	var orderReq order.CreateOrderReq
	if err := c.ShouldBindJSON(&orderReq); err != nil {
		apiReqErr := errors.BadRequestError("invalid request")
		c.JSON(apiReqErr.Code, apiReqErr)
		return
	}

	orderClient := services.OrderService()
	result, err := orderClient.CreateOrder(&orderReq)
	if err != nil {
		c.JSON(err.Code, err)
		return
	}
	c.JSON(http.StatusCreated, result)
}

func CompleteOrder(c *gin.Context) {

	orderId := c.Param("order_id")
	orderId = strings.TrimSpace(orderId)
	if orderId == "" {
		apiReqErr := errors.BadRequestError("order id is empty")
		c.JSON(apiReqErr.Code, apiReqErr)
		return
	}

	orderIdInt, userErr := strconv.ParseInt(orderId, 10, 64)
	if userErr != nil {
		apiReqErr := errors.BadRequestError("order id is not int")
		c.JSON(apiReqErr.Code, apiReqErr)
		return
	}

	orderClient := services.OrderService()
	orderResp, err := orderClient.CompleteOrder(orderIdInt)
	if err != nil {
		c.JSON(err.Code, err)
		return
	}
	c.JSON(http.StatusOK, orderResp)
}

func GetOrder(c *gin.Context) {
	orderId := c.Param("order_id")
	orderId = strings.TrimSpace(orderId)
	if orderId == "" {
		apiReqErr := errors.BadRequestError("order id is empty")
		c.JSON(apiReqErr.Code, apiReqErr)
		return
	}

	orderIdInt, userErr := strconv.ParseInt(orderId, 10, 64)
	if userErr != nil {
		apiReqErr := errors.BadRequestError("order id is not int")
		c.JSON(apiReqErr.Code, apiReqErr)
		return
	}

	orderClient := services.OrderService()
	orderResp, err := orderClient.GetOrder(orderIdInt)
	if err != nil {
		c.JSON(err.Code, err)
		return
	}
	c.JSON(http.StatusOK, orderResp)
}
