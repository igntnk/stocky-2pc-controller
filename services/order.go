package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/igntnk/stocky-2pc-controller/clients"
	"github.com/igntnk/stocky-2pc-controller/models"
	"github.com/igntnk/stocky-2pc-controller/protobufs/iims_pb"
	"github.com/igntnk/stocky-2pc-controller/protobufs/oms_pb"
)

type OrderService interface {
	CreateOrder(ctx context.Context, comment, userId, staffId string, products []*clients.OrderProductInput) (models.Order, error)
	TccCreateOrder(ctx context.Context, req models.OrderCreateRequest) (*models.OrderResponse, error)
}

func NewOrderService(oms clients.OMSClient,
	sms clients.SMSClient,
	scs clients.SCSClient,
	iims clients.IIMSClient) OrderService {

	return &orderService{
		OMSClient:  oms,
		SMSClient:  sms,
		SCSClient:  scs,
		IIMSClient: iims,
	}
}

type orderService struct {
	OMSClient  clients.OMSClient
	SMSClient  clients.SMSClient
	SCSClient  clients.SCSClient
	IIMSClient clients.IIMSClient
}

func (o orderService) CreateOrder(ctx context.Context, comment, userId, staffId string, products []*clients.OrderProductInput) (res models.Order, err error) {
	startProductAmount := make(map[string]float32)

	defer func() {
		if err != nil {
			for uuid, amount := range startProductAmount {
				o.SMSClient.SetProductAmount(ctx, uuid, amount)
			}
		}
	}()

	for _, product := range products {
		amount, err := o.SMSClient.GetProductAmount(ctx, product.ProductUUID)
		if err != nil {
			return models.Order{}, err
		}

		if amount < float32(product.Amount) {
			return models.Order{}, errors.Join(ErrProductOut, errors.New(fmt.Sprintf("product - %s", product.ProductUUID)))
		}
		startProductAmount[product.ProductUUID] = amount

		_, err = o.SMSClient.SetProductAmount(ctx, product.ProductUUID, amount-float32(product.Amount))
		if err != nil {
			return models.Order{}, err
		}
	}

	var order *oms_pb.Order
	order, err = o.OMSClient.CreateOrder(ctx, comment, userId, staffId, products)
	if err != nil {
		return models.Order{}, err
	}

	defer func() {
		if err != nil {
			_ = o.OMSClient.DeleteOrder(ctx, order.Uuid)
		}
	}()

	resultProducts := make([]models.OrderProduct, len(order.Products))
	for i, product := range order.Products {
		var productData *iims_pb.GetProductMessage

		productData, err = o.IIMSClient.GetByProductCode(ctx, product.ProductCode)
		if err != nil {
			return models.Order{}, err
		}

		resultProducts[i] = models.OrderProduct{
			Uuid:        product.ProductUuid,
			Name:        productData.Name,
			Description: productData.Description,
			Price:       productData.Price,
		}
	}

	return models.Order{
		Uuid:     order.Uuid,
		Comment:  order.Comment,
		UserId:   order.UserId,
		StaffId:  order.StaffId,
		Price:    order.OrderCost,
		Products: resultProducts,
	}, nil
}

func (s *orderService) TccCreateOrder(ctx context.Context, req models.OrderCreateRequest) (*models.OrderResponse, error) {
	products := make([]*oms_pb.OrderProductInput, len(req.Products))
	for i, product := range req.Products {
		products[i] = &oms_pb.OrderProductInput{
			ProductUuid: product.ProductID.String(),
			Amount:      int32(product.Amount),
		}
	}

	res, err := s.OMSClient.TCCOrderCreation(ctx, &oms_pb.CreateOrderRequest{
		Comment:  req.Comment,
		UserId:   req.UserID,
		StaffId:  req.StaffID,
		Products: products,
	})
	if err != nil {
		return nil, err
	}

	resProducts := make([]models.ProductDetail, len(res.Products))
	for i, product := range res.Products {
		resProducts[i] = models.ProductDetail{
			Price:       product.ResultPrice,
			ProductCode: product.ProductCode,
			Amount:      int(product.Amount),
			TotalPrice:  product.ResultPrice * float64(product.Amount),
		}
	}

	return &models.OrderResponse{
		ID:           res.Uuid,
		Comment:      res.Comment,
		UserID:       res.UserId,
		StaffID:      res.StaffId,
		OrderCost:    res.OrderCost,
		Status:       models.OrderStatusNew,
		CreationDate: res.CreationDate.String(),
		Products:     resProducts,
	}, nil
}
