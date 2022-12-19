package payments

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type PaymentController struct {
	repo CreditCardRepository
}

func NewPaymentController(repo CreditCardRepository) *PaymentController {
	return &PaymentController{repo}
}

func (b *PaymentController) ValidatePayment(w http.ResponseWriter, r *http.Request) {
	var paymentRequest PaymentRequest
	err := json.NewDecoder(r.Body).Decode(&paymentRequest)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Unable to parse request body. Error: "+err.Error())
		return
	}

	paymentResponse, err := b.Validate(paymentRequest)
	if err != nil {
		if e, ok := err.(ValidationErrors); ok {
			log.Errorf("Invalid payment request %v", e.Errors)
			RespondWithJSON(w, http.StatusBadRequest, e.Errors)
			return
		}
		log.Errorf("Error while validating payment %v", err)
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusOK, paymentResponse)
}

func (b *PaymentController) Validate(paymentRequest PaymentRequest) (*PaymentResponse, error) {
	err := paymentRequest.Validate()
	if err != nil {
		return nil, err
	}
	creditCard, err := b.repo.GetCreditCardByNumber(*paymentRequest.CardNumber)
	if err != nil {
		return nil, err
	}
	if creditCard.Cvv == *paymentRequest.Cvv &&
		creditCard.ExpiryMonth == *paymentRequest.ExpiryMonth &&
		creditCard.ExpiryYear == *paymentRequest.ExpiryYear {
		return &PaymentResponse{Status: "ACCEPTED"}, nil
	}
	return &PaymentResponse{Status: "REJECTED"}, nil
}

func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
