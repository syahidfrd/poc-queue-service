package sql

import (
	"github.com/jinzhu/gorm"
	"poc-misreported-qty/model"
)

type sqlProductStore struct {
	gormInstance *gorm.DB
}

func NewProductStore(gormInstance *gorm.DB) (productStore model.ProductStore) {
	gormInstance.AutoMigrate(&model.Product{})
	productStore = &sqlProductStore{
		gormInstance: gormInstance,
	}
	return
}

func (p *sqlProductStore) Exist(product *model.Product) (isExist bool) {
	isExist = !p.gormInstance.NewRecord(product)
	return
}

func (p *sqlProductStore) Save(product *model.Product) (err error) {
	err = p.gormInstance.Save(product).Error
	return
}

func (p *sqlProductStore) FindOneBy(query map[string]interface{}) (product *model.Product, err error) {
	product = &model.Product{}
	err = p.gormInstance.First(product, query).Error
	return
}