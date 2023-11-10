package transport

import (
	"blockchaincrawler/internal/auth/config"
	pb "blockchaincrawler/pkg/protobuf/userservice/gw"
	"context"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserGrpcTransport struct {
	cfg    config.UserGrpcTransport
	client pb.UserServiceClient
}

func NewUserGrpcTransport(config config.UserGrpcTransport) (*UserGrpcTransport, error) {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	conn, err := grpc.Dial(config.Host+config.Port, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to dial user grpc server on %v%v error: %w", config.Host, config.Port, err)
	}

	client := pb.NewUserServiceClient(conn)

	return &UserGrpcTransport{
		cfg:    config,
		client: client,
	}, nil
}

func (t *UserGrpcTransport) GetUserByEmail(ctx context.Context, email string) (*pb.User, error) {
	res, err := t.client.GetUserByEmail(ctx, &pb.GetUserByEmailRequest{
		Email: email,
	})
	if err != nil {
		return nil, fmt.Errorf("grpcClient.GetUserByEmail error: %w", err)
	}

	if res == nil {
		return nil, fmt.Errorf("not found")
	}

	return res.Result, nil
}

func (t *UserGrpcTransport) CreateUser(ctx context.Context, user *pb.User) (*uuid.UUID, error) {
	res, err := t.client.CreateUser(ctx, &pb.CreateUserRequest{
		User: user,
	})
	if err != nil {
		return &uuid.Nil, fmt.Errorf("grpcClient.CreateUser error: %w", err)
	}

	id, err := uuid.Parse(res.Result)
	if err != nil {
		return &uuid.Nil, fmt.Errorf("failed to parse uuid error: %w", err)
	}

	return &id, nil
}

func (t *UserGrpcTransport) DeleteUserById(ctx context.Context, id uuid.UUID) error {
	res, err := t.client.DeleteUserById(ctx, &pb.DeleteUserByIdRequest{
		Id: id.String(),
	})
	if err != nil {
		return fmt.Errorf("grpcClient.DeleteUserById error: %w", err)
	}

	_ = res
	return nil
}

func (t *UserGrpcTransport) IsValidLogin(ctx context.Context, email string, password string) (bool, error) {
	res, err := t.client.IsValidLogin(ctx, &pb.IsValidLoginRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return false, fmt.Errorf("grpcClient.IsValidLogin error: %w", err)
	}

	return res.IsValid, nil
}
