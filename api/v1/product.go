package v1

import (
	"poc-misreported-qty/model"
	"poc-misreported-qty/util/validator"
)

type ProductHandler struct {
	ValidatorService validator.Service
	ProductStore     model.ProductStore
}

func NewProductHandler(validatorService validator.Service, productStore model.ProductStore) (productHandler *ProductHandler) {
	productHandler = &ProductHandler{
		ValidatorService: validatorService,
		ProductStore:     productStore,
	}
	return
}
