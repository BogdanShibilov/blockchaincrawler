package transport

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/bogdanshibilov/blockchaincrawler/internal/auth/config"
	pb "github.com/bogdanshibilov/blockchaincrawler/pkg/protobuf/user/gw"
)

type User struct {
	client pb.UserServiceClient
	cfg    config.UserTransport
}

func NewUser(cfg config.UserTransport) (*User, error) {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	conn, err := grpc.Dial(cfg.Host+":"+cfg.Port, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to dial user grpc server on %v:%v error: %w", cfg.Host, cfg.Port, err)
	}

	client := pb.NewUserServiceClient(conn)

	return &User{
		client: client,
		cfg:    cfg,
	}, nil
}

func (u *User) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	res, err := u.client.CreateUser(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (u *User) GetUserByEmail(ctx context.Context, req *pb.GetUserByEmailRequest) (*pb.GetUserByEmailResponse, error) {
	res, err := u.client.GetUserByEmail(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (u *User) GetUserById(ctx context.Context, req *pb.GetUserByIdRequest) (*pb.GetUserByIdResponse, error) {
	return u.client.GetUserById(ctx, req)
}

func (u *User) GetAllUsers(ctx context.Context) (*pb.GetAllUsersResponse, error) {
	res, err := u.client.GetAllUsers(ctx, &pb.GetAllUsersRequest{})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (u *User) DeleteUserById(ctx context.Context, req *pb.DeleteUserByIdRequest) error {
	res, err := u.client.DeleteUserById(ctx, req)
	if err != nil {
		return err
	}

	_ = res
	return nil
}

func (u *User) ConfirmUser(ctx context.Context, req *pb.ConfirmUserRequest) error {
	res, err := u.client.ConfirmUser(ctx, req)
	if err != nil {
		return err
	}

	_ = res
	return nil
}
