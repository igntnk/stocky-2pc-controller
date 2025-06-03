package clients

import (
	"context"
	"github.com/igntnk/stocky-2pc-controller/protobufs/iims_pb"
	"google.golang.org/grpc"
)

type IIMSClient interface {
	// Product methods
	InsertProduct(ctx context.Context, name, description, creationDate string, price float32) (string, error)
	GetProducts(ctx context.Context, limit, offset int64) ([]*iims_pb.GetProductMessage, error)
	GetById(ctx context.Context, productId string) (*iims_pb.GetProductMessage, error)
	GetByProductCode(ctx context.Context, productCode string) (*iims_pb.GetProductMessage, error)
	DeleteProduct(ctx context.Context, id string) error
	UpdateProduct(ctx context.Context, id, name, description, creationDate string, price float32) error
	BlockProduct(ctx context.Context, id string) error
	UnblockProduct(ctx context.Context, id string) error

	// Sale methods
	InsertSale(ctx context.Context, name, description string, saleSize int32, product string) (string, error)
	GetSales(ctx context.Context, limit, offset int64) ([]*iims_pb.GetSaleMessage, error)
	DeleteSale(ctx context.Context, id string) error
	UpdateSale(ctx context.Context, id, name, description string, saleSize int32) error
	BlockSale(ctx context.Context, id string) error
	UnblockSale(ctx context.Context, id string) error
}

func NewIIMSClient(conn *grpc.ClientConn) IIMSClient {
	productClient := iims_pb.NewProductServiceClient(conn)
	saleClient := iims_pb.NewSaleServiceClient(conn)

	return &iimsClient{
		productClient: productClient,
		saleClient:    saleClient,
	}
}

type iimsClient struct {
	productClient iims_pb.ProductServiceClient
	saleClient    iims_pb.SaleServiceClient
}

// Product methods implementation
func (c *iimsClient) InsertProduct(ctx context.Context, name, description, creationDate string, price float32) (string, error) {
	req := &iims_pb.InsertProductRequest{
		Name:         name,
		Description:  description,
		CreationDate: creationDate,
		Price:        price,
	}
	resp, err := c.productClient.InsertOne(ctx, req)
	if err != nil {
		return "", err
	}
	return resp.Id, nil
}

func (c *iimsClient) GetProducts(ctx context.Context, limit, offset int64) ([]*iims_pb.GetProductMessage, error) {
	req := &iims_pb.GetProductsRequest{
		Limit:  limit,
		Offset: offset,
	}
	resp, err := c.productClient.Get(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Products, nil
}

func (c *iimsClient) DeleteProduct(ctx context.Context, id string) error {
	req := &iims_pb.DeleteProductRequest{
		Id: id,
	}
	_, err := c.productClient.Delete(ctx, req)
	return err
}

func (c *iimsClient) GetById(ctx context.Context, productId string) (*iims_pb.GetProductMessage, error) {
	return c.productClient.GetById(ctx, &iims_pb.GetByIdProductRequest{
		Id: productId,
	})
}

func (c *iimsClient) GetByProductCode(ctx context.Context, code string) (*iims_pb.GetProductMessage, error) {
	return c.productClient.GetByProductCode(ctx, &iims_pb.GetByProductCodeRequest{
		Code: code,
	})
}

func (c *iimsClient) UpdateProduct(ctx context.Context, id, name, description, creationDate string, price float32) error {
	req := &iims_pb.UpdateProductRequest{
		Id:           id,
		Name:         name,
		Description:  description,
		CreationDate: creationDate,
		Price:        price,
	}
	_, err := c.productClient.Update(ctx, req)
	return err
}

func (c *iimsClient) BlockProduct(ctx context.Context, id string) error {
	req := &iims_pb.BlockProductOperationMessage{
		Id: id,
	}
	_, err := c.productClient.BlockProduct(ctx, req)
	return err
}

func (c *iimsClient) UnblockProduct(ctx context.Context, id string) error {
	req := &iims_pb.BlockProductOperationMessage{
		Id: id,
	}
	_, err := c.productClient.UnblockProduct(ctx, req)
	return err
}

// Sale methods implementation
func (c *iimsClient) InsertSale(ctx context.Context, name, description string, saleSize int32, product string) (string, error) {
	req := &iims_pb.InsertSaleRequest{
		Name:        name,
		Description: description,
		SaleSize:    saleSize,
		Product:     product,
	}
	resp, err := c.saleClient.InsertOne(ctx, req)
	if err != nil {
		return "", err
	}
	return resp.Id, nil
}

func (c *iimsClient) GetSales(ctx context.Context, limit, offset int64) ([]*iims_pb.GetSaleMessage, error) {
	req := &iims_pb.GetSalesRequest{
		Limit:  limit,
		Offset: offset,
	}
	resp, err := c.saleClient.Get(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Sales, nil
}

func (c *iimsClient) DeleteSale(ctx context.Context, id string) error {
	req := &iims_pb.DeleteSaleRequest{
		Id: id,
	}
	_, err := c.saleClient.Delete(ctx, req)
	return err
}

func (c *iimsClient) UpdateSale(ctx context.Context, id, name, description string, saleSize int32) error {
	req := &iims_pb.UpdateSaleRequest{
		Id:          id,
		Name:        name,
		Description: description,
		SaleSize:    saleSize,
	}
	_, err := c.saleClient.Update(ctx, req)
	return err
}

func (c *iimsClient) BlockSale(ctx context.Context, id string) error {
	req := &iims_pb.BlockSaleOperationMessage{
		Id: id,
	}
	_, err := c.saleClient.BlockSale(ctx, req)
	return err
}

func (c *iimsClient) UnblockSale(ctx context.Context, id string) error {
	req := &iims_pb.BlockSaleOperationMessage{
		Id: id,
	}
	_, err := c.saleClient.UnblockSale(ctx, req)
	return err
}
