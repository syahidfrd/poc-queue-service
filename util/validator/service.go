package validator

// Service validator interface
type Service interface {
	ValidateForm(form interface{}) (validationErrors map[string][]string)
	AddValidationError(validationErrors map[string][]string, field, errorMessage string) (result map[string][]string)
}
