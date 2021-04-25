package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/serenize/snaker"
)

// DefaultValidatorService structure
type DefaultValidatorService struct {
	validatorInstance *validator.Validate
}

// NewDefaultValidatorService func
func NewDefaultValidatorService() (defaultValidatorService Service) {
	defaultValidatorService = &DefaultValidatorService{
		validatorInstance: validator.New(),
	}

	return
}

// ValidateForm func
func (v *DefaultValidatorService) ValidateForm(form interface{}) (validationErrors map[string][]string) {

	err := v.validatorInstance.Struct(form)
	if err == nil {
		return nil
	}

	validationErrors = map[string][]string{}
	for _, err := range err.(validator.ValidationErrors) {
		fieldName := snaker.CamelToSnake(err.Field())
		errorMessage := err.ActualTag()
		fmt.Println(errorMessage)
		validationErrors[fieldName] = append(validationErrors[fieldName], errorMessage)
	}

	if len(validationErrors) == 0 {
		return nil
	}

	return
}

// AddValidationError func
func (v *DefaultValidatorService) AddValidationError(validationErrors map[string][]string, field, errorMessage string) (result map[string][]string) {
	if validationErrors == nil {
		validationErrors = map[string][]string{}
	}

	_, exist := validationErrors[field]
	if !exist {
		validationErrors[field] = []string{}
	}
	validationErrors[field] = append(validationErrors[field], errorMessage)

	result = validationErrors
	return
}
