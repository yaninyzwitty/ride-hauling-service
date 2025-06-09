package pkg

import (
	"fmt"
	"log/slog"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration structure
type Config struct {
	APIGateway    APIGateway    `yaml:"api-gateway"`
	DriverService DriverService `yaml:"driver-service"`
	TripService   TripService   `yaml:"trip-service"`
}

type APIGateway struct {
	Port int `yaml:"port"`
}

type DriverService struct {
	Port int `yaml:"port"`
}

type TripService struct {
	Port int `yaml:"port"`
}

// LoadConfig reads and parses the configuration file
func (c *Config) LoadConfig(configPath string) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		slog.Error("error reading config file", "error", err)
		return fmt.Errorf("error reading config file: %w", err)
	}
	err = yaml.Unmarshal(data, c)
	if err != nil {
		slog.Error("error parsing config file", "error", err)
		return fmt.Errorf("error parsing config file: %w", err)
	}
	return nil
}
