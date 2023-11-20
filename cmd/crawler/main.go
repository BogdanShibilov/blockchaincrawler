package main

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/bogdanshibilov/blockchaincrawler/internal/crawler/app"
	"github.com/bogdanshibilov/blockchaincrawler/internal/crawler/config"
	"github.com/bogdanshibilov/blockchaincrawler/pkg/logger"
)

func main() {
	logger := logger.NewZap()
	defer logger.Sync()

	cfg, err := loadConfig("./../../config/crawler")
	if err != nil {
		logger.Panicf("failed to load config error: %v", err)
	}

	app := app.New(logger, &cfg)
	app.Run()
}

func loadConfig(path string) (config config.Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return config, fmt.Errorf("failed to ReadInConfig err: %w", err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return config, fmt.Errorf("failed to Unmarshal config err: %w", err)
	}
	return config, nil
}
