package handlers

import (
	"net/http"
	"strings"

	"github.com/fvaiiii/ordering_products/order/internal/api/http/dto"
	"github.com/fvaiiii/ordering_products/order/internal/models"
	"github.com/fvaiiii/ordering_products/order/internal/service"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	service *service.OrderService
}

func NewOrderService(s *service.OrderService) *OrderHandler {
	return &OrderHandler{service: s}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req dto.CreateOrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	order, err := h.service.CreateOrder(c.Request.Context(), req.UserUuid, req.ProductUuids)
	if err != nil {
		errMsg := err.Error()
		switch {
		case strings.Contains(errMsg, "user_uuid is required") || strings.Contains(errMsg, "at least one product is required"):
			c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})

		case strings.Contains(errMsg, "some products are not exist"):
			c.JSON(http.StatusNotFound, gin.H{"error": errMsg})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": errMsg})
		}
		return
	}

	resp := dto.CreateOrderResponse{
		OrderUUID:  order.OrderUuid,
		TotalPrice: order.TotalPrice,
	}

	c.JSON(http.StatusOK, resp)
}

func (h *OrderHandler) PayOrder(c *gin.Context) {
	orderUuid := c.Param("order_uuid")
	if orderUuid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order_uuid is required"})
		return
	}

	var req dto.CreatePaymentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	transactionUuid, err := h.service.PayOrder(
		c.Request.Context(),
		orderUuid,
		models.PaymentMethod(req.PaymentMethod),
	)
	if err != nil {
		errMsg := err.Error()
		switch {
		case strings.Contains(errMsg, "order not found"):
			c.JSON(http.StatusNotFound, gin.H{"error": errMsg})
		case strings.Contains(errMsg, "order cannot be paid"):
			c.JSON(http.StatusConflict, gin.H{"error": errMsg})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	resp := dto.CreatePaymentResponse{
		TransactionUuid: transactionUuid,
	}

	c.JSON(http.StatusOK, resp)
}

func (h *OrderHandler) GetByUUID(c *gin.Context) {
	orderUuid := c.Param("order_uuid")
	if orderUuid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order_uuid is required"})
		return
	}

	order, err := h.service.GetOrderByUUID(c.Request.Context(), orderUuid)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "order not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
		}
		return

	}

	resp := dto.GetPaymentRequestResponse{
		OrderUuid:    order.OrderUuid,
		UserUuid:     order.UserUuid,
		ProductUuids: order.ProductUuids,
		TotalPrice:   order.TotalPrice,
		Status:       string(order.Status),
	}

	if order.TransactionUuid != nil {
		resp.TransactionUuid = *order.TransactionUuid
	}
	if order.PaymentMethod != nil {
		resp.PaymentMethod = *order.PaymentMethod
	}
	c.JSON(http.StatusOK, resp)
}

func (h *OrderHandler) CancelOrder(c *gin.Context) {
	orderUuid := c.Param("order_uuid")
	if orderUuid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order_uuid is required"})
		return
	}

	err := h.service.CancelOrder(c.Request.Context(), orderUuid)
	if err != nil {
		errMsg := err.Error()
		switch {
		case strings.Contains(errMsg, "order not found"):
			c.JSON(http.StatusNotFound, gin.H{
				"error": "order not found",
			})
		case strings.Contains(errMsg, "cannot cancel paid order"):
			c.JSON(http.StatusConflict, gin.H{
				"error": "the order has been paid and cannot be cancelled",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
		}
		return
	}

	c.Status(http.StatusNoContent)
}
