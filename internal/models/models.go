package models

import (
	"database/sql/driver"
	"errors"
	"time"

	"github.com/goccy/go-json"
)

type PaymentModel struct {
	Transaction  string  `json:"transaction" validate:"required"`
	RequestID    string  `json:"request_id" validate:"required"`
	Currency     string  `json:"currency" validate:"required"`
	Provider     string  `json:"provider" validate:"required"`
	Amount       float32 `json:"amount" validate:"required"`
	PaymentDt    int     `json:"payment_dt" validate:"required"`
	Bank         string  `json:"bank" validate:"required"`
	DeliveryCost float32     `json:"delivery_cost" validate:"required"`
	GoodsTotal   float32     `json:"goods_total" validate:"required"`
	CustomFee    float32     `json:"custom_fee"`
}


type DeliveryModel struct {
	Name    string `json:"name" validate:"required"`
	Phone   string `json:"phone" validate:"required"`
	Zip     string `json:"zip" validate:"required"`
	City    string `json:"city" validate:"required"`
	Address string `json:"address" validate:"required"`
	Region  string `json:"region" validate:"required"`
	Email   string `json:"email" validate:"required,email"`
}

type OrderItemModel struct {
	ChrtID      int     `json:"chrt_id" validate:"required"`
	TrackNumber string  `json:"track_number" validate:"required"`
	Price       float32 `json:"price" validate:"required"`
	Rid         string  `json:"rid" validate:"required"`
	Name        string  `json:"name" validate:"required"`
	Sale        int     `json:"sale" validate:"required"`
	Size        string  `json:"size" validate:"required"`
	TotalPrice  float32 `json:"total_price" validate:"required"`
	NmID        int     `json:"nm_id" validate:"required"`
	Brand       string  `json:"brand" validate:"required"`
	Status      int     `json:"status" validate:"required"`
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
	OrderUID          string            `json:"order_uid" db:"order_uid" validate:"required"`
	TrackNumber       string            `json:"track_number" db:"track_number" validate:"required"`
	Entry             string            `json:"entry" db:"entry" validate:"required"`
	Delivery          DeliveryModel    	`json:"delivery" db:"delivery" validate:"required"`
	Payment           PaymentModel     	`json:"payment" db:"payment" validate:"required"`
	Items             OrderItems 		`json:"items" db:"items" validate:"required"`
	Locale            string            `json:"locale" db:"locale" validate:"required"`
	InternalSignature string            `json:"internal_signature" db:"internal_signature"`
	CustomerID        string            `json:"customer_id" db:"customer_id" validate:"required"`
	DeliveryService   string            `json:"delivery_service" db:"delivery_service" validate:"required"`
	Shardkey          string            `json:"shardkey" db:"shardkey" validate:"required"`
	SmID              int               `json:"sm_id" db:"sm_id"`
	DateCreated       time.Time         `json:"date_created" db:"date_created" validate:"required"`
	OofShard          string            `json:"oof_shard" db:"oof_shard" validate:"required"`
}
