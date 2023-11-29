package orders

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/leenzstra/WBL0/internal/models"
	"github.com/leenzstra/WBL0/pkg/database"
)

type IOrderRepo interface {
	Add(ctx context.Context, order models.OrderModel) error
	Get(ctx context.Context, uid string) (*models.OrderModel, error)
}

type OrderRepo struct {
	*database.DB
}

func NewRepo(db *database.DB) *OrderRepo {
	return &OrderRepo{
		DB: db,
	}
}

func (r *OrderRepo) Add(ctx context.Context, order models.OrderModel) error {
	sql, args, err := r.Builder.
		Insert("orders").
		Columns("order_uid", "track_number", "entry", "delivery", "payment","items","locale",
			"internal_signature","customer_id", "delivery_service", "shardkey","sm_id",
			"date_created","oof_shard").
		Values(order.OrderUID, order.TrackNumber, order.Entry, order.Delivery, order.Payment, order.Items, order.Locale,
			order.InternalSignature, order.CustomerID, order.DeliveryService, order.Shardkey, order.SmID, 
			order.DateCreated, order.OofShard).
		ToSql()

	fmt.Printf("%s %+v %+v\n", sql, args, order)
	
	if err != nil {
		return fmt.Errorf("OrderRepo.Add.r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("OrderRepo.Add.r.Pool.Exec: %w", err)
	}

	return nil
}

func (r *OrderRepo) Get(ctx context.Context, uid string) (*models.OrderModel, error) {
	sql, args, err := r.Builder.
		Select("*").
		From("orders").
		Where(squirrel.Eq{"order_uid": uid}).
		ToSql()
	// fmt.Println(sql)
	// fmt.Println(args)
	if err != nil {
		return nil, fmt.Errorf("OrderRepo.Get.r.Builder: %w", err)
	}

	o := make([]models.OrderModel, 1)
	err = pgxscan.Select(ctx, r.Pool, &o, sql, args...)
	// row := r.Pool.QueryRow(ctx, sql, args...)

	// err = row.Scan(o)
	if err != nil {
		return nil, fmt.Errorf("OrderRepo.Get.row.Scan: %w", err)
	}

	if len(o) != 1 {
		return nil, fmt.Errorf("OrderRepo.Get.row.len: len != 1")
	}

	return &o[0], nil
}


