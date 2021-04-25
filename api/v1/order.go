package v1

import (
	"poc-misreported-qty/model"
	"poc-misreported-qty/util/queue"
	"poc-misreported-qty/util/validator"
)

type OrderHandler struct {
	ValidatorService validator.Service
	QueueService     queue.Service
	ProductStore     model.ProductStore
	OrderStore       model.OrderStore
}

func NewOrderHandler(validatorService validator.Service, queueService queue.Service, productStore model.ProductStore, orderStore model.OrderStore) (orderHandler *OrderHandler) {
	orderHandler = &OrderHandler{
		ValidatorService: validatorService,
		QueueService: queueService,
		ProductStore: productStore,
		OrderStore: orderStore,
	}
	return
}
