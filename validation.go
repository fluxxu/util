package util

type ValidationErrors map[string][]string

type ValidationContext struct {
	errors ValidationErrors
}

func (v *ValidationContext) AddError(key, message string) {
	v.errors[key] = append(v.errors[key], message)
}

func (v *ValidationContext) HasError() bool {
	return len(v.errors) != 0
}

func (v *ValidationContext) Errors() ValidationErrors {
	return v.errors
}

func (v *ValidationContext) ToError() *ValidationError {
	return NewValidationError(v)
}

func NewValidationContext() *ValidationContext {
	return &ValidationContext{
		errors: make(ValidationErrors),
	}
}

type ValidationErrorInterface interface {
	ValidationErrors() ValidationErrors
}

type ValidationError struct {
	errors ValidationErrors
}

func (ve *ValidationError) Error() string {
	return "validation error"
}

func (ve *ValidationError) ValidationErrors() ValidationErrors {
	return ve.errors
}

func (ve *ValidationError) ToResponseData() map[string]interface{} {
	return map[string]interface{}{"errors": ve.errors}
}

func NewValidationError(v *ValidationContext) *ValidationError {
	return &ValidationError{errors: v.errors}
}
