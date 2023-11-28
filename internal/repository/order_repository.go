package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/leenzstra/WBL0/db"
	"go.uber.org/zap"
)

type OrderRepository struct {
	db *pgx.Conn
	logger *zap.Logger
}

func NewOrderRepository(conn *pgx.Conn, logger *zap.Logger) OrderRepository {
	return OrderRepository{
		db: conn,
	}
}

func (repo *OrderRepository) Insert(ctx context.Context, order *db.OrderModel) error {
	q := `
		INSERT INTO orders
			(order_uid, track_number, entry, locale, internal_signature, customer_id, 
				delivery_service, shardkey, sm_id, date_created, oof_shard, items, delivery, payment) 
		VALUES 
			(@order_uid, @track_number, @entry, @locale, @internal_signature, @customer_id, 
				@delivery_service, @shardkey, @sm_id, @date_created, @oof_shard, @items, @delivery, @payment) 
	`

	args := pgx.NamedArgs{
		"order_uid": order.OrderUID,
		"track_number": order.TrackNumber,
		"entry": order.Entry,
		"locale": order.Locale,
		"internal_signature": order.InternalSignature,
		"customer_id": order.CustomerID,
		"delivery_service": order.DeliveryService,
		"shardkey": order.Shardkey,
		"sm_id": order.SmID,
		"date_created": order.DateCreated,
		"oof_shard": order.OofShard,
		"items": order.Items,
		"payment":order.Payment,
		"delivery":order.Delivery,
	}

	_, err := repo.db.Exec(ctx, q, args)
	if err != nil {
		return err
	}

	return nil

}

func (repo *OrderRepository) Get(ctx context.Context, orderId string) (*db.OrderModel, error) {
	q := `SELECT * FROM orders LIMIT 1`

	rows, err := repo.db.Query(ctx, q)
	if err != nil {
		repo.logger.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	orders := []*db.OrderModel{}
	for rows.Next() {
		order := &db.OrderModel{}
		err := rows.Scan(order)
		if err != nil {
			return nil, fmt.Errorf("unable to scan row: %w", err)
		}
		orders = append(orders, order)
	}

	if len(orders) != 1 {
		return nil, fmt.Errorf("returned orders count = %d", len(orders))
	}

  	return orders[0], nil
}

