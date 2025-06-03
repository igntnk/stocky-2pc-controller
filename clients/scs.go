package clients

import (
	"context"
	"github.com/igntnk/stocky-2pc-controller/protobufs/scs_pb"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type SCSClient interface {
	// User methods
	CreateUser(
		ctx context.Context,
		name, description, documentType, documentNumber, authID string,
	) (string, error)
	BlockUser(ctx context.Context, id string) (string, error)
	UnblockUser(ctx context.Context, id string) (string, error)
	UpdateUser(
		ctx context.Context,
		id, name, description, documentType, documentNumber string,
	) (string, error)
	GetUserByID(ctx context.Context, id string) (*scs_pb.UserModel, error)
	GetAllUsers(ctx context.Context) ([]*scs_pb.UserModel, error)
}

func NewSCSClient(conn *grpc.ClientConn) SCSClient {
	userClient := scs_pb.NewUserServiceClient(conn)

	return &scsClient{
		userClient: userClient,
	}
}

type scsClient struct {
	userClient scs_pb.UserServiceClient
}

// User methods implementation
func (c *scsClient) CreateUser(
	ctx context.Context,
	name, description, documentType, documentNumber, authID string,
) (string, error) {
	req := &scs_pb.CreateUserRequest{
		Name:           name,
		Description:    description,
		DocumentType:   documentType,
		DocumentNumber: documentNumber,
		AuthId:         authID,
	}
	resp, err := c.userClient.CreateUser(ctx, req)
	if err != nil {
		return "", err
	}
	return resp.Id, nil
}

func (c *scsClient) BlockUser(ctx context.Context, id string) (string, error) {
	req := &scs_pb.IdRequest{
		Id: id,
	}
	resp, err := c.userClient.BlockUser(ctx, req)
	if err != nil {
		return "", err
	}
	return resp.Id, nil
}

func (c *scsClient) UnblockUser(ctx context.Context, id string) (string, error) {
	req := &scs_pb.IdRequest{
		Id: id,
	}
	resp, err := c.userClient.UnblockUser(ctx, req)
	if err != nil {
		return "", err
	}
	return resp.Id, nil
}

func (c *scsClient) UpdateUser(
	ctx context.Context,
	id, name, description, documentType, documentNumber string,
) (string, error) {
	req := &scs_pb.UpdateUserRequest{
		Id:             id,
		Name:           name,
		Description:    description,
		DocumentType:   documentType,
		DocumentNumber: documentNumber,
	}
	resp, err := c.userClient.UpdateUser(ctx, req)
	if err != nil {
		return "", err
	}
	return resp.Id, nil
}

func (c *scsClient) GetUserByID(ctx context.Context, id string) (*scs_pb.UserModel, error) {
	req := &scs_pb.IdRequest{
		Id: id,
	}
	return c.userClient.GetById(ctx, req)
}

func (c *scsClient) GetAllUsers(ctx context.Context) ([]*scs_pb.UserModel, error) {
	resp, err := c.userClient.GetAllUsers(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	return resp.Users, nil
}
