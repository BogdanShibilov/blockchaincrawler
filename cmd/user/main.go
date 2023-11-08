package main

import (
	"blockchaincrawler/internal/user/app"
	"blockchaincrawler/internal/user/config"
	"fmt"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	sugaredLogger := logger.Sugar()
	sugaredLogger = sugaredLogger.With(zap.String("app", "user-service"))

	cfg, err := loadConfig("config/user")
	if err != nil {
		sugaredLogger.Fatalf("failed to load config error: %w", err)
	}

	app := app.NewApp(sugaredLogger, &cfg)
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
