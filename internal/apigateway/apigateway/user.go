package apigateway

import (
	"context"

	"github.com/bogdanshibilov/blockchaincrawler/internal/apigateway/controller/http/v1/dto"
)

func (a *ApiGateway) GetAllUsers(ctx context.Context) ([]*dto.UserDto, error) {
	res, err := a.users.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	usersDto := make([]*dto.UserDto, 0)
	for _, u := range res.Users {
		user := &dto.UserDto{
			Id:          u.Id,
			Role:        u.Role,
			Email:       u.Email,
			IsConfirmed: u.IsConfirmed,
		}
		usersDto = append(usersDto, user)
	}

	return usersDto, nil
}

func (a *ApiGateway) DeleteUserById(ctx context.Context, id string) error {
	return a.users.DeleteUserById(ctx, id)
}

func (a *ApiGateway) UpdateProfile(ctx context.Context, id string, p *dto.UserProfileDto) error {
	return a.users.UpdateProfile(ctx, id, p)
}

func (a *ApiGateway) GetProfileById(ctx context.Context, id string) (*dto.UserProfileDto, error) {
	return a.users.GetProfileById(ctx, id)
}
