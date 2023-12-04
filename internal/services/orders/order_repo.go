package orders

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/leenzstra/WBL0/internal/models"
	"github.com/leenzstra/WBL0/pkg/database"
)

type IOrderRepo interface {
	Add(ctx context.Context, order models.OrderModel) error
	Get(ctx context.Context, uid string) (*models.OrderModel, error)
	GetAll(ctx context.Context) ([]models.OrderModel, error)
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
	if err != nil {
		return nil, fmt.Errorf("OrderRepo.Get.r.Builder: %w", err)
	}

	var o []models.OrderModel
	err = r.DB.Scanner.Select(ctx, r.Pool, &o, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("OrderRepo.Get.Select: %w", err)
	}

	if len(o) != 1 {
		return nil, fmt.Errorf("OrderRepo.Get.row.len: len != 1")
	}

	return &o[0], nil
}

func (r *OrderRepo) GetAll(ctx context.Context) ([]models.OrderModel, error) {
	sql, _, err := r.Builder.
		Select("*").
		From("orders").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("OrderRepo.GetAll.r.Builder: %w", err)
	}

	var o []models.OrderModel
	err = r.DB.Scanner.Select(ctx, r.Pool, &o, sql)
	if err != nil {
		return nil, fmt.Errorf("OrderRepo.GetAll.Select: %w", err)
	}

	return o, nil
}

