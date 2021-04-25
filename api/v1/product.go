package v1

import (
	"github.com/gin-gonic/gin"
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

func (p *ProductHandler) CreateProduct(ctx *gin.Context) {
	form := &struct {
		Name     string `form:"name" json:"name" validate:"required"`
		Quantity uint32 `form:"quantity" json:"quantity" validate:"required"`
		Price    uint64 `form:"price" json:"price" validate:"required"`
	}{}

	if err := ctx.Bind(form); err != nil {
		httpBindingErrorResponse(ctx, err)
		return
	}

	if validationErrors := p.ValidatorService.ValidateForm(form); validationErrors != nil {
		httpValidationErrorResponse(ctx, validationErrors)
		return
	}

	product := model.NewProduct(form.Name, form.Quantity, form.Price)
	if err := p.ProductStore.Save(product); err !=nil {
		httpInternalServerErrorResponse(ctx, err.Error())
		return
	}

	httpOkResponse(ctx, map[string]interface{}{
		"message": "Create product success",
	})
}
