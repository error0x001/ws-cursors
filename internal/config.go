package internal

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Address      string
	IsSSLUsing   bool
	Port         string
	TemplatePath string
	ShutDownTime time.Duration
}

func (c *Config) getAddressPort() string {
	return fmt.Sprintf("%s:%s", c.Address, c.Port)
}

func (c *Config) GetOnlyPort() string {
	return fmt.Sprintf(":%s", c.Port)
}

func (c *Config) GetWSPath() string {
	if c.IsSSLUsing {
		return fmt.Sprintf("wss://%s/ws", c.Address)
	}
	return fmt.Sprintf("ws://%s/ws", c.getAddressPort())
}

func NewConfig() Config {
	rawShutdownTime := os.Getenv("SHUTDOWN_TIME")
	shutdownTime, err := strconv.Atoi(rawShutdownTime)
	if err != nil {
		log.Fatalln("invalid SHUTDOWN_TIME format")
	}
	return Config{
		Address:      os.Getenv("ADDRESS"),
		Port:         os.Getenv("PORT"),
		TemplatePath: os.Getenv("TEMPLATE_PATH"),
		ShutDownTime: time.Duration(shutdownTime) * time.Second,
		IsSSLUsing:   os.Getenv("IS_SSL_USING") == "1",
	}
}
