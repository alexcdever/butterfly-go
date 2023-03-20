package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Config 是数据库连接配置
type Config struct {
	Dialect      string // 数据库类型，mysql 或 postgres
	Host         string // 数据库地址
	Port         int    // 数据库端口号
	User         string // 数据库用户名
	Password     string // 数据库密码
	DatabaseName string // 数据库名
}

// NewConnection 返回一个新的 GORM 数据库连接
func NewConnection(cfg Config) (*gorm.DB, error) {
	var dsn string
	switch cfg.Dialect {
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DatabaseName)
		return gorm.Open(mysql.Open(dsn), &gorm.Config{})
	case "postgres":
		dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DatabaseName)
		return gorm.Open(postgres.Open(dsn), &gorm.Config{})
	default:
		return nil, fmt.Errorf("unsupported dialect: %s", cfg.Dialect)
	}
}
