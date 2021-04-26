package v1

import (
	"poc-misreported-qty/model"
	"time"
)

func formatProduct(product *model.Product) (result map[string]interface{}) {
	result = map[string]interface{}{
		"id": product.GetID(),
		"name": product.GetName(),
		"price": product.GetPrice(),
		"quantity": product.GetQuantity(),
		"created_at": product.GetCreatedAt().Format(time.RFC3339),
	}
	return
}

func formatProducts(products []*model.Product) (result []map[string]interface{}) {
	result = []map[string]interface{}{}
	for _, product := range products {
		result = append(result, formatProduct(product))
	}
	return
}
