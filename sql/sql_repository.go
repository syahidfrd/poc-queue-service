package sql

import (
	"github.com/jinzhu/gorm"
	"poc-misreported-qty/model"
)

func NewSQLDataRepository(gormInstance *gorm.DB) (dataRepository model.DataRepository) {
	productStore := NewProductStore(gormInstance)

	dataRepository = model.DataRepository{
		ProductStore: productStore,
	}
	return
}
