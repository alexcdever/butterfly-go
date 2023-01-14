package dao

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Dao struct {
	DB *gorm.DB
}

type CreateTableApi interface {
	HaveTable() bool
	CreateTable() (err error)
}

type CrudApi interface {
	Create() (err error)
	Delete() (err error)
	Update() (err error)
	Read() (err error)
}

func (d Dao) CreateTable(instance interface{}) (err error) {
	if !d.DB.Migrator().HasTable(instance) {
		err = d.DB.Migrator().CreateTable(instance)
	}
	if err != nil {
		err = errors.Wrap(err, "failed to create table for butterfly")
	}
	return err
}

func (d Dao) HasTable(instance interface{}) bool {
	return d.DB.Migrator().HasTable(instance)
}
