package items

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
)

type RepoItemsPostgres struct {
	DB *sql.DB
}

func NewRepoItemsPostgres(db *sql.DB) (*RepoItemsPostgres, error) {
	return &RepoItemsPostgres{
		DB: db,
	}, nil
}

func (op *RepoItemsPostgres) AddPayment(items []Item, ctx context.Context) error {
	for _, item := range items {
		itemUUID, err := uuid.NewUUID()
		if err != nil {
			return err
		}
		item.ItemUUID = itemUUID
		_, err = op.DB.ExecContext(ctx, "INSERT INTO item VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12);", item.ItemUUID, item.ChrtID, item.TrackNumber, item.Price, item.RID, item.Name, item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status)
		if err != nil {
			return err
		}
	}
	return nil
}
