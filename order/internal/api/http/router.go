package http

import (
	"github.com/fvaiiii/ordering_products/order/internal/api/http/handlers"
	"github.com/fvaiiii/ordering_products/order/internal/api/http/middleware"
	"github.com/gin-gonic/gin"
)

func NewRouter(order *handlers.OrderHandler) *gin.Engine {
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(middleware.RequestID())

	api := router.Group("/api/v1")
	{
		api.POST("/orders", order.CreateOrder)
		api.POST("/orders/:order_uuid/pay", order.PayOrder)
		api.GET("/orders/:order_uuid", order.GetOrder)
		api.POST("orders/:order_uuid/cancel", order.CancelOrder)
	}

	return router
}
