package payments

import (
	"strings"
)

type PaymentRequest struct {
	CardNumber  *string `json:"cardNumber"`
	Cvv         *string `json:"cvv"`
	ExpiryMonth *int    `json:"expiryMonth"`
	ExpiryYear  *int    `json:"expiryYear"`
}

func (p PaymentRequest) Validate() error {
	validationErrors := ValidationErrors{Errors: make(map[string]string)}
	if p.CardNumber == nil || strings.TrimSpace(*p.CardNumber) == "" {
		validationErrors.AddError("CardNumber", "CardNumber is required")
	}
	if p.Cvv == nil || strings.TrimSpace(*p.Cvv) == "" {
		validationErrors.AddError("CVV", "CVV is required")
	}
	if p.ExpiryMonth == nil {
		validationErrors.AddError("ExpiryMonth", "ExpiryMonth is required")
	} else if *p.ExpiryMonth < 1 || *p.ExpiryMonth > 12 {
		validationErrors.AddError("ExpiryMonth", "Invalid ExpiryMonth value")
	}
	if p.ExpiryYear == nil {
		validationErrors.AddError("ExpiryYear", "ExpiryYear is required")
	}
	if len(validationErrors.Errors) == 0 {
		return nil
	}
	return validationErrors
}

type PaymentResponse struct {
	Status string `json:"status"`
}

type ValidationErrors struct {
	Errors map[string]string
}

func (v ValidationErrors) AddError(property, message string) {
	v.Errors[property] = message
}

func (v ValidationErrors) Error() string {
	return "Validation error"
}
