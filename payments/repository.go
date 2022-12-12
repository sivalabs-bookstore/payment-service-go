package payments

import (
	"context"
	"database/sql"
)

type CreditCard struct {
	CardNumber  string `json:"cardNumber"`
	Cvv         string `json:"cvv"`
	ExpiryMonth int    `json:"expiryMonth"`
	ExpiryYear  int    `json:"expiryYear"`
}

type CreditCardRepository interface {
	GetCreditCardByNumber(cardNumber string) (CreditCard, error)
}

type creditCardRepo struct {
	db *sql.DB
}

func NewCreditCardRepo(db *sql.DB) CreditCardRepository {
	var repo CreditCardRepository = creditCardRepo{db}
	return repo
}

func (b creditCardRepo) GetCreditCardByNumber(cardNumber string) (CreditCard, error) {
	// log.Infof("Fetching cardNumber with number=%s", cardNumber)
	ctx := context.Background()
	tx, err := b.db.BeginTx(ctx, nil)
	if err != nil {
		return CreditCard{}, err
	}
	var creditCard = CreditCard{}
	query := "select card_number, cvv, expiry_month, expiry_year FROM credit_cards where card_number=$1"
	err = tx.QueryRowContext(ctx, query, cardNumber).Scan(
		&creditCard.CardNumber, &creditCard.Cvv, &creditCard.ExpiryMonth, &creditCard.ExpiryYear)
	if err != nil {
		tx.Rollback()
		return CreditCard{}, err
	}
	err = tx.Commit()
	if err != nil {
		return CreditCard{}, err
	}
	return creditCard, nil
}
