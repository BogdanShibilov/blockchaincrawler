package transport

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/bogdanshibilov/blockchaincrawler/internal/apigateway/config"
	"github.com/bogdanshibilov/blockchaincrawler/internal/apigateway/controller/http/v1/dto"
	pb "github.com/bogdanshibilov/blockchaincrawler/pkg/protobuf/user/gw"
)

type User struct {
	client pb.UserServiceClient
	cfg    *config.UserTransport
}

func NewUser(cfg *config.UserTransport) (*User, error) {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	conn, err := grpc.Dial(cfg.Host+":"+cfg.Port, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to dial auth grpc server on %v:%v error: %w", cfg.Host, cfg.Port, err)
	}

	client := pb.NewUserServiceClient(conn)

	return &User{
		client: client,
		cfg:    cfg,
	}, nil
}

func (u *User) GetAllUsers(ctx context.Context) (*pb.GetAllUsersResponse, error) {
	req := &pb.GetAllUsersRequest{}

	res, err := u.client.GetAllUsers(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (u *User) DeleteUserById(ctx context.Context, id string) error {
	req := &pb.DeleteUserByIdRequest{
		Id: id,
	}

	_, err := u.client.DeleteUserById(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) UpdateProfile(ctx context.Context, id string, profile *dto.UserProfileDto) error {
	req := &pb.UpdateProfileRequest{
		UserId: id,
		Profile: &pb.Profile{
			Name:    profile.Name,
			Surname: profile.Surname,
			AboutMe: profile.AboutMe,
		},
	}

	_, err := u.client.UpdateProfile(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) GetProfileById(ctx context.Context, id string) (*dto.UserProfileDto, error) {
	req := &pb.GetProfileByIdRequest{
		Id: id,
	}

	res, err := u.client.GetProfileById(ctx, req)
	if err != nil {
		return nil, err
	}

	return &dto.UserProfileDto{
		Name:    res.Profile.Name,
		Surname: res.Profile.Surname,
		AboutMe: res.Profile.AboutMe,
	}, nil
}
