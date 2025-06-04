package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/igntnk/stocky-2pc-controller/clients"
	"github.com/igntnk/stocky-2pc-controller/models"
	"github.com/igntnk/stocky-2pc-controller/requests"
	"github.com/igntnk/stocky-2pc-controller/services"
	"net/http"
	"sync"
)

type orderController struct {
	orders services.OrderService
	mu     sync.Mutex
}

func NewOrderController(orders services.OrderService) Controller {
	return &orderController{
		orders: orders,
		mu:     sync.Mutex{},
	}
}

func (o *orderController) Register(r *gin.Engine) {
	group := r.Group("/api/order")
	group.POST("/create", o.Create)

	tccGroup := r.Group("/api/TCC/order")
	tccGroup.POST("/create", o.TCCCreate)
}

func (o *orderController) Create(context *gin.Context) {
	o.mu.Lock()
	defer o.mu.Unlock()

	receivedOrder := requests.CreateOrder{}
	err := context.ShouldBindBodyWithJSON(&receivedOrder)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errors.Join(err, errors.New("failed to parse body")).Error()})
		return
	}

	products := []*clients.OrderProductInput{}
	for _, product := range receivedOrder.Products {
		products = append(products, &clients.OrderProductInput{
			ProductUUID: product.Uuid,
			Amount:      int32(product.Amount),
		})
	}

	order, err := o.orders.CreateOrder(context, receivedOrder.Comment, "000000000000000000000000", "000000000000000000000000", products)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"order": order})
}

func (o *orderController) TCCCreate(context *gin.Context) {
	var err error

	receivedOrder := requests.CreateOrder{}
	err = context.ShouldBindBodyWithJSON(&receivedOrder)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": errors.Join(err, errors.New("failed to parse body")).Error()})
		return
	}

	products := []models.OrderProductInput{}
	for _, product := range receivedOrder.Products {
		prUuid, err := uuid.Parse(product.Uuid)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		products = append(products, models.OrderProductInput{
			ProductID: prUuid,
			Amount:    int(product.Amount),
		})
	}

	order, err := o.orders.TccCreateOrder(context, models.OrderCreateRequest{
		UserID:   "000000000000000000000000",
		StaffID:  "000000000000000000000000",
		Comment:  receivedOrder.Comment,
		Products: products,
	})
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"order": order})
}
