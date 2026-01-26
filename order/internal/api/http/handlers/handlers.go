package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/fvaiiii/ordering_products/order/internal/api/http/dto"
	"github.com/fvaiiii/ordering_products/order/internal/models"
	"github.com/fvaiiii/ordering_products/order/internal/service"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	service service.OrderService
}

func NewOrderHandler(s service.OrderService) *OrderHandler {
	return &OrderHandler{service: s}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req dto.CreateOrderRequest

	log.Printf("[Handler] CreateOrder called with: %v", req)
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[Handler] JSON parse error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
		return
	}

	log.Printf("[Handler] Calling service.CreateOrder: user=%s, products=%v", req.UserUuid, req.ProductUuids)

	order, err := h.service.CreateOrder(c.Request.Context(), req.UserUuid, req.ProductUuids)
	if err != nil {
		log.Printf("[Handler] CreateOrder service error: %v", err)
		errMsg := err.Error()
		switch {
		case strings.Contains(errMsg, "user_uuid is required") ||
			strings.Contains(errMsg, "at least one product is required"):
			c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
		case strings.Contains(errMsg, "some products are not exist"):
			c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
		default:
			log.Printf("[Handler] Unexpected error in CreateOrder: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	resp := dto.CreateOrderResponse{
		OrderUUID:  order.OrderUuid,
		TotalPrice: order.TotalPrice,
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *OrderHandler) PayOrder(c *gin.Context) {
	orderUUID := c.Param("order_uuid")
	if orderUUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order_uuid is required"})
		return
	}

	var req dto.PayOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var paymentMethod models.PaymentMethod
	switch strings.ToUpper(req.PaymentMethod) {
	case "CARD":
		paymentMethod = models.PaymentMethodCard
	case "SBP":
		paymentMethod = models.PaymentMethodSBP
	case "CREDIT_CARD":
		paymentMethod = models.PaymentMethodCreditCard
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payment method"})
		return
	}

	transactionUUID, err := h.service.PayOrder(c.Request.Context(), orderUUID, paymentMethod)
	if err != nil {
		errMsg := err.Error()
		switch {
		case strings.Contains(errMsg, "order not found"):
			c.JSON(http.StatusNotFound, gin.H{"error": errMsg})
		case strings.Contains(errMsg, "order cannot be paid"):
			c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	resp := dto.PayOrderResponse{
		TransactionUUID: transactionUUID,
	}
	c.JSON(http.StatusOK, resp)
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	orderUUID := c.Param("order_uuid")
	if orderUUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order_uuid is required"})
		return
	}

	order, err := h.service.GetOrderByUUID(c.Request.Context(), orderUUID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	resp := dto.GetOrderResponse{
		OrderUUID:    order.OrderUuid,
		UserUUID:     order.UserUuid,
		ProductUUIDs: order.ProductUuids,
		TotalPrice:   order.TotalPrice,
		Status:       string(order.Status),
	}

	if order.TransactionUuid != nil {
		resp.TransactionUUID = *order.TransactionUuid
	}
	if order.PaymentMethod != nil {
		resp.PaymentMethod = *order.PaymentMethod
	}

	c.JSON(http.StatusOK, resp)
}

func (h *OrderHandler) CancelOrder(c *gin.Context) {
	orderUUID := c.Param("order_uuid")
	if orderUUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order_uuid is required"})
		return
	}

	err := h.service.CancelOrder(c.Request.Context(), orderUUID)
	if err != nil {
		errMsg := err.Error()
		switch {
		case strings.Contains(errMsg, "order not found"):
			c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		case strings.Contains(errMsg, "cannot cancel paid order"):
			c.JSON(http.StatusConflict, gin.H{"error": "cannot cancel paid order"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.Status(http.StatusNoContent)
}
