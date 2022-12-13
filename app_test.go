package main

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/sivalabs-bookstore/payment-service-go/config"
	"github.com/sivalabs-bookstore/payment-service-go/payments"
	pgtc "github.com/sivalabs-bookstore/payment-service-go/test"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var cfg config.AppConfig
var app *App
var router *mux.Router

func TestMain(m *testing.M) {
	//Common Setup
	ctx := context.Background()
	pgContainer, err := pgtc.SetupPostgres(ctx)
	if err != nil {
		log.Fatalf("failed to setup Postgres container")
	}
	defer pgContainer.CloseFn()
	overrideEnv(pgContainer)

	cfg = config.GetConfig()
	app = NewApp(cfg)
	router = app.Router

	code := m.Run()

	//Common Teardown
	os.Exit(code)
}

func overrideEnv(pgContainer *pgtc.PostgresContainer) {
	os.Setenv("APP_DB_HOST", pgContainer.Host)
	os.Setenv("APP_DB_PORT", pgContainer.Port)
	os.Setenv("APP_DB_USERNAME", pgContainer.Username)
	os.Setenv("APP_DB_PASSWORD", pgContainer.Password)
	os.Setenv("APP_DB_NAME", pgContainer.Database)
	os.Setenv("APP_DB_RUN_MIGRATIONS", "true")
}

func TestValidatePaymentAccepted(t *testing.T) {
	body := strings.NewReader(`
			{
				"cardNumber": "1111222233334444",
				"cvv": "123",
				"expiryMonth": 2,
				"expiryYear": 2030
			}
		`)
	req, _ := http.NewRequest(http.MethodPost, "/api/payments/validate", body)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	actualResponseJson := response.Body
	var paymentResponse payments.PaymentResponse
	err := json.NewDecoder(actualResponseJson).Decode(&paymentResponse)
	assert.NoError(t, err)
	assert.Equal(t, "ACCEPTED", paymentResponse.Status, "Expected status: ACCEPTED. Got %s", paymentResponse.Status)
}

func TestValidatePaymentRejected(t *testing.T) {
	body := strings.NewReader(`
			{
				"cardNumber": "1111222233334444",
				"cvv": "456",
				"expiryMonth": 2,
				"expiryYear": 2030
			}
		`)
	req, _ := http.NewRequest(http.MethodPost, "/api/payments/validate", body)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	actualResponseJson := response.Body
	var paymentResponse payments.PaymentResponse
	err := json.NewDecoder(actualResponseJson).Decode(&paymentResponse)
	assert.NoError(t, err)
	assert.Equal(t, "REJECTED", paymentResponse.Status, "Expected status: REJECTED. Got %s", paymentResponse.Status)
}

func TestValidatePaymentBadRequest(t *testing.T) {
	body := strings.NewReader(`
			{
				"expiryMonth": 2,
				"expiryYear": 2030
			}
		`)
	req, _ := http.NewRequest(http.MethodPost, "/api/payments/validate", body)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
