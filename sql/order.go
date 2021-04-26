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
	tx := o.gormInstance.Begin()
	err = tx.Save(order).Error
	if err != nil {
		tx.Rollback()
		return
	}

	sql := "UPDATE products SET quantity = quantity - ? WHERE id = ?"
	err = tx.Exec(sql, order.GetQuantity(), order.GetProductID()).Error
	if err != nil {
		tx.Rollback()
		return
	}

	tx.Commit()
	return
}

func (o *sqlOrderStore) FindOneBy(query map[string]interface{}) (order *model.Order, err error) {
	order = &model.Order{}
	err = o.gormInstance.First(order, query).Error
	return
}

func (o *sqlOrderStore) FindAll() (orders []*model.Order, err error) {
	err = o.gormInstance.Preload("Product").Order("created_at DESC").Find(&orders).Error
	return
}