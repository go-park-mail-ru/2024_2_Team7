package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

const (
	DefaultPort = "8080"
)

type Config struct {
	Port                    string
	AuthServiceAddr         string
	UserServiceAddr         string
	EventServiceAddr        string
	ImageServiceAddr        string
	CSATServiceAddr         string
	NotificationServiceAddr string
}

func LoadConfig() (Config, error) {
	var conf Config
	err := godotenv.Load()
	if err != nil {
		return conf, err
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = DefaultPort
	}
	conf.Port = port

	conf.AuthServiceAddr = os.Getenv("AUTH_SERVICE_ADDR")
	if conf.AuthServiceAddr == "" {
		return Config{}, errors.New("Failed to get auth service address")
	}

	conf.UserServiceAddr = os.Getenv("USER_SERVICE_ADDR")
	if conf.UserServiceAddr == "" {
		return Config{}, errors.New("Failed to get user service address")
	}

	conf.EventServiceAddr = os.Getenv("EVENT_SERVICE_ADDR")
	if conf.EventServiceAddr == "" {
		return Config{}, errors.New("Failed to get event service address")
	}

	conf.ImageServiceAddr = os.Getenv("IMAGE_SERVICE_ADDR")
	if conf.ImageServiceAddr == "" {
		return Config{}, errors.New("Failed to get image service address")
	}

	conf.CSATServiceAddr = os.Getenv("CSAT_SERVICE_ADDR")
	if conf.CSATServiceAddr == "" {
		return Config{}, errors.New("Failed to get csat service address")
	}

	conf.NotificationServiceAddr = os.Getenv("NOTIFICATION_SERVICE_ADDR")
	if conf.NotificationServiceAddr == "" {
		return Config{}, errors.New("Failed to get notification service address")
	}

	return conf, nil
}
