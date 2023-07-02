package payment

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
)

type RepoPaymentPostgres struct {
	DB *sql.DB
}

func NewRepoPaymentPostgres(db *sql.DB) (*RepoPaymentPostgres, error) {
	return &RepoPaymentPostgres{
		DB: db,
	}, nil
}

func (op *RepoPaymentPostgres) AddPayment(ctx context.Context, item Payment) (*uuid.UUID, error) {
	paymentUUID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	item.PaymentUUID = paymentUUID
	_, err = op.DB.ExecContext(ctx, "INSERT INTO payment VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11);", item.PaymentUUID, item.Transaction, item.RequestID, item.Currency, item.Provider, item.Amount, item.PaymentDt, item.Bank, item.DeliveryCost, item.GoodsTotal, item.CustomFee)
	if err != nil {
		return nil, err
	}
	return &paymentUUID, nil
}
