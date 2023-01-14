package config

import (
	"gorm.io/gorm"
)

// DbConfig is a configuration instance for connecting DB
type DbConfig struct {
	Username string
	Password string
	Host     string
	Port     int
	DB       *gorm.DB
}

type DbConfigApi interface {
	GenerateConnectLink() (link string)
	Connect()
}
