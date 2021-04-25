package sql

import (
	"github.com/jinzhu/gorm"
	"poc-misreported-qty/model"
)

type sqlOrderStore struct {
	gormInstance *gorm.DB
}

func NewOrderStore(gormInstance *gorm.DB) (orderStore model.OrderStore) {
	gormInstance.AutoMigrate(&model.Order{})
	orderStore = &sqlOrderStore{
		gormInstance: gormInstance,
	}
	return
}

func (o *sqlOrderStore) Exist(order *model.Order) (isExist bool) {
	isExist = !o.gormInstance.NewRecord(order)
	return
}

func (o *sqlOrderStore) Save(order *model.Order) (err error) {
	err = o.gormInstance.Save(order).Error
	return
}

func (o *sqlOrderStore) FindOneBy(query map[string]interface{}) (order *model.Order, err error) {
	order = &model.Order{}
	err = o.gormInstance.First(order, query).Error
	return
}