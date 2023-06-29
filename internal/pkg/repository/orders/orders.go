package orders

import "github.com/google/uuid"

type Orders struct {
	OrderUID          string      `json:"order_uid"`
	TrackNumber       string      `json:"track_number"`
	Pass              string      `json:"pass"`
	Entry             string      `json:"entry"`
	Delivery          uuid.UUID   `json:"delivery"`
	Payment           uuid.UUID   `json:"payment"`
	Items             []uuid.UUID `json:"items"`
	Locale            string      `json:"locale"`
	InternalSignature string      `json:"internal_signature"`
	CustomerID        string      `json:"customer_id"`
	DeliveryService   string      `json:"delivery_service"`
	Shardkey          string      `json:"shardkey"`
	SmID              int         `json:"sm_id"`
	DateCreated       time.Time   `json:"date_created"`
	OofShard          string      `json:"oof_shard"`
}

type OrderRepo interface {
	AddOrder(item *Links) error
	GetOrderByID(url string) (*Links, error)
	GetShortenLink(url string) (*Links, error)
}
