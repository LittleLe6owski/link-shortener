package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/LittleLe6owski/link-shortener/config"
	"github.com/LittleLe6owski/link-shortener/internal/instance"
	"github.com/caarlos0/env/v6"
	"github.com/go-playground/validator/v10"
)

func parseConfig() (config.Config, error) {
	cfg := config.Config{}

	if err := env.Parse(&cfg); err != nil {
		return config.Config{}, fmt.Errorf(
			"failed to parse config for environment variables: %w", err,
		)
	}

	if err := validator.New().Struct(cfg); err != nil {
		return config.Config{}, fmt.Errorf("failed to validate environment config: %w", err)
	}

	return cfg, nil
}

func main() {
	cfg, err := parseConfig()
	if err != nil {
		log.Default().Fatal(err)
	}

	linkShortener, err := instance.New(cfg)
	if err != nil {
		log.Default().Fatal(err)
	}

	if err := linkShortener.Run(context.Background()); err != nil {
		os.Exit(1)
	}
}
