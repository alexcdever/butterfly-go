package config

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type PgConfig struct {
	DbConfig
	Postgres struct {
		Database string
		SslMode  string
		TimeZone string
	}
}

func (config *PgConfig) GenerateConnectLink() (link string) {
	link = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%v sslmode=%s TimeZone=%s",
		config.Host, config.Username, config.Password, config.Postgres.Database, config.Port,
		config.Postgres.SslMode, config.Postgres.TimeZone)
	return link
}

func (config *PgConfig) Connect() {
	db, err := gorm.Open(postgres.Open(config.GenerateConnectLink()), &gorm.Config{})
	if err != nil {
		log.Panicf("failed to connect to DB: %s", err)
	}
	config.DB = db
}
