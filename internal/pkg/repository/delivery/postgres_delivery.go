package delivery

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
)

type RepoDeliveryPostgres struct {
	DB *sql.DB
}

func NewRepoDeliveryPostgres(db *sql.DB) (*RepoDeliveryPostgres, error) {
	return &RepoDeliveryPostgres{
		DB: db,
	}, nil
}

func (op *RepoDeliveryPostgres) AddDelivery(item Delivery, ctx context.Context) (*uuid.UUID, error) {
	deliveryUUID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	item.DeliveryUUID = deliveryUUID
	_, err = op.DB.ExecContext(ctx, "INSERT INTO delivery VALUES ($1,$2,$3,$4,$5,$6,$7,$8);", item.DeliveryUUID, item.Name, item.Phone, item.Zip, item.City, item.Address, item.Region, item.Email)
	if err != nil {
		return err
	}
	return nil
}
