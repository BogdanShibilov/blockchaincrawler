package auth

import (
	"github.com/bogdanshibilov/blockchaincrawler/internal/auth/config"
	"github.com/bogdanshibilov/blockchaincrawler/internal/auth/transport"
)

type Service struct {
	cfg           config.Auth
	userTransport *transport.User
}

func New(
	cfg config.Auth,
	userTransport *transport.User,
) UseCase {

	return &Service{
		cfg:           cfg,
		userTransport: userTransport,
	}
}
