package models

import (
	"database/sql/driver"
	"github.com/goccy/go-json"
	"errors"
	"time"
)

type PaymentModel struct {
	Transaction  string  `json:"transaction"`
	RequestID    string  `json:"request_id"`
	Currency     string  `json:"currency"`
	Provider     string  `json:"provider"`
	Amount       float32 `json:"amount"`
	PaymentDt    int     `json:"payment_dt"`
	Bank         string  `json:"bank"`
	DeliveryCost int     `json:"delivery_cost"`
	GoodsTotal   int     `json:"goods_total"`
	CustomFee    int     `json:"custom_fee"`
}


type DeliveryModel struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type OrderItemModel struct {
	ChrtID      int     `json:"chrt_id"`
	TrackNumber string  `json:"track_number"`
	Price       float32 `json:"price"`
	Rid         string  `json:"rid"`
	Name        string  `json:"name"`
	Sale        int     `json:"sale"`
	Size        string  `json:"size"`
	TotalPrice  float32 `json:"total_price"`
	NmID        int     `json:"nm_id"`
	Brand       string  `json:"brand"`
	Status      int     `json:"status"`
}

type OrderItems []OrderItemModel

func (s *OrderItems) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *OrderItems) Scan(src interface{}) error {
    switch v := src.(type) {
    case []byte:
        return json.Unmarshal(v, s)
    case string:
        return json.Unmarshal([]byte(v), s)
    }
    return errors.New("type assertion failed")
}

type OrderModel struct {
	OrderUID          string            `json:"order_uid" db:"order_uid"`
	TrackNumber       string            `json:"track_number" db:"track_number"`
	Entry             string            `json:"entry" db:"entry"`
	Delivery          DeliveryModel    	`json:"delivery" db:"delivery"`
	Payment           PaymentModel     	`json:"payment" db:"payment"`
	Items             OrderItems 		`json:"items" db:"items"`
	Locale            string            `json:"locale" db:"locale"`
	InternalSignature string            `json:"internal_signature" db:"internal_signature"`
	CustomerID        string            `json:"customer_id" db:"customer_id"`
	DeliveryService   string            `json:"delivery_service" db:"delivery_service"`
	Shardkey          string            `json:"shardkey" db:"shardkey"`
	SmID              int               `json:"sm_id" db:"sm_id"`
	DateCreated       time.Time         `json:"date_created" db:"date_created"`
	OofShard          string            `json:"oof_shard" db:"oof_shard"`
}
