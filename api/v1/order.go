package v1

import (
	"github.com/gin-gonic/gin"
	"poc-queue-service/model"
	"poc-queue-service/util/queue"
	"poc-queue-service/util/validator"
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

func (o *OrderHandler) CreateOrder(ctx *gin.Context) {
	form := &struct {
		ProductID uint64 `form:"product_id" json:"product_id" validate:"required"`
		Quantity uint32 `form:"quantity" json:"quantity" validate:"required"`
	}{}

	if err := ctx.Bind(form); err != nil {
		httpBindingErrorResponse(ctx, err)
		return
	}

	if validationErrors := o.ValidatorService.ValidateForm(form); validationErrors != nil {
		httpValidationErrorResponse(ctx, validationErrors)
		return
	}

	product, _ := o.ProductStore.FindOneBy(map[string]interface{}{
		"id": form.ProductID,
	})

	if !o.ProductStore.Exist(product) {
		httpValidationErrorResponse(ctx, map[string][]string{
			"product_id": {"not found"},
		})
		return
	}

	if product.GetQuantity() < form.Quantity {
		httpValidationErrorResponse(ctx, map[string][]string{
			"quantity": {"invalid"},
		})
		return
	}

	// Publish order product queue
	o.QueueService.PublishOrderQueue(product.GetID(), form.Quantity)

	httpOkResponse(ctx, map[string]interface{}{
		"message": "Create order success",
	})
}

func (o *OrderHandler) GetAllOrder(ctx *gin.Context) {
	orders, err := o.OrderStore.FindAll()
	if err != nil {
		httpInternalServerErrorResponse(ctx, err.Error())
		return
	}

	httpOkResponse(ctx, map[string]interface{}{
		"orders": formatOrders(orders),
	})
}
