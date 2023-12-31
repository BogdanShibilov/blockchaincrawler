package postgres

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/bogdanshibilov/blockchaincrawler/internal/user/config"
	"github.com/bogdanshibilov/blockchaincrawler/internal/user/entity"
)

type Pg struct {
	*gorm.DB
}

type Config config.DbNode

func (c Config) dsn() string {
	return fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=%s",
		c.Host, c.Port, c.Name, c.User, c.Password, c.SslMode)
}

func NewWithGorm(cfg config.DbNode) (*Pg, error) {
	conf := Config(cfg)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  conf.dsn(),
		PreferSimpleProtocol: true,
	}))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	err = db.AutoMigrate(&entity.User{}, &entity.Profile{})
	if err != nil {
		return nil, fmt.Errorf("failed to auto migrate: %w", err)
	}

	return &Pg{db}, nil
}

func (pg *Pg) Close() error {
	db, err := pg.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to retrieve DB: %w", err)
	}

	return db.Close()
}
