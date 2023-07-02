package orders

import (
	"context"
	"database/sql"
)

type RepoOrderPostgres struct {
	DB *sql.DB
}

func NewRepoOrderPostgres(db *sql.DB) (*RepoOrderPostgres, error) {
	return &RepoOrderPostgres{
		DB: db,
	}, nil
}

func (op *RepoOrderPostgres) AddOrder(item Order, ctx context.Context) error {
	_, err := op.DB.ExecContext(ctx, "INSERT INTO orders VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15);", item.OrderUID, item.TrackNumber, item.Pass, item.Entry, item.Delivery, item.Payment, item.Items, item.Locale, item.InternalSignature, item.CustomerID, item.DeliveryService, item.Shardkey, item.SmID, item.DateCreated, item.OofShard)
	if err != nil {
		return err
	}
	return nil
}
