package app

import (
	"createorder/controllers/order"
	"createorder/logger"
	"github.com/gin-gonic/gin"
)

var router = gin.Default()

func Start() {
	router.POST("/Order/CreateOrder", order.CreateOrder)
	router.POST("/Order/CompleteOrder/:order_id", order.CompleteOrder)
	router.GET("/Order/GetOrder/:order_id", order.GetOrder)
	router.GET("/Payment/AddPaymentItem", order.AddPaymentItem)
	router.GET("/Payment/UpdatePaymentItem", order.UpdatePaymentItem)
	logger.Info("application to start")
	router.Run(":5551")
}
