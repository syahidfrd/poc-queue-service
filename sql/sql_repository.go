package sql

import (
	"github.com/jinzhu/gorm"
	"poc-queue-service/model"
)

func NewSQLDataRepository(gormInstance *gorm.DB) (dataRepository model.DataRepository) {
	productStore := NewProductStore(gormInstance)
	orderStore := NewOrderStore(gormInstance)

	dataRepository = model.DataRepository{
		ProductStore: productStore,
		OrderStore: orderStore,
	}
	return
}
