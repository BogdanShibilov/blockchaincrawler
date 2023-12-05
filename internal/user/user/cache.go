package user

import (
	"context"
	"encoding/json"
	"time"

	"github.com/bogdanshibilov/blockchaincrawler/internal/user/entity"
)

func (s *Service) getUserByEmailCached(ctx context.Context, key string) (*entity.User, error) {
	value := s.cache.Client.Get(ctx, key).Val()
	if value == "" {
		return nil, nil
	}

	var user *entity.User
	err := json.Unmarshal([]byte(value), &user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) setUserByEmailCache(ctx context.Context, key string, value *entity.User) error {
	userJson, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return s.cache.Set(ctx, key, string(userJson), time.Minute*30).Err()
}
