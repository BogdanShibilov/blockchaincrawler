package auth

import (
	"github.com/bogdanshibilov/blockchaincrawler/internal/auth/codeproducer"
	"github.com/bogdanshibilov/blockchaincrawler/internal/auth/config"
	"github.com/bogdanshibilov/blockchaincrawler/internal/auth/transport"
	"github.com/bogdanshibilov/blockchaincrawler/pkg/redis"
)

type Service struct {
	cfg           config.Auth
	userTransport *transport.User
	codeProducer  *codeproducer.CodeProducer
	codeDb        *redis.Redis
}

func New(
	cfg config.Auth,
	userTransport *transport.User,
	codeProducer *codeproducer.CodeProducer,
	codeDb *redis.Redis,
) UseCase {

	return &Service{
		cfg:           cfg,
		userTransport: userTransport,
		codeProducer:  codeProducer,
		codeDb:        codeDb,
	}
}
