package dao

import (
	"butterfly/db/model"
	"log"
)

type ButterflyDao struct {
	Dao
}

func Delete() {}
func Update() {}
func (d ButterflyDao) ReadLast() (record *model.Butterfly) {
	d.DB.Last(record)
	return record
}

func (d ButterflyDao) Create(
	instance *model.Butterfly) {
	if err := d.DB.Model(&model.Butterfly{}).Create(instance).Error; err != nil {
		log.Panicf("failed to insert the data for butterfly: %v", err)
	}
}

func (d ButterflyDao) BatchCreate(dataList []model.Butterfly) {
	if err := d.DB.Model(&model.Butterfly{}).Create(&dataList).Error; err != nil {
		log.Panicf("failed to batch insert the data for butterfly: %v", err)
	}
}
