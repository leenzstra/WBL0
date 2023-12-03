package stan

import (
	"context"
	"errors"

	"github.com/goccy/go-json"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/leenzstra/WBL0/internal/models"
	"github.com/leenzstra/WBL0/internal/services/orders"
	"github.com/leenzstra/WBL0/internal/validation"
	"github.com/nats-io/stan.go"
	"go.uber.org/zap"
)

func handleOrderFunc(data []byte, logger *zap.Logger, service *orders.OrdersService) error {
	var order models.OrderModel
	if err := json.Unmarshal(data, &order); err != nil {
		logger.Error("Failed to unmarshal order", zap.String("err", err.Error()))
		return err
	}

	v := validation.NewOrderValidator()
	if err := v.ValidateOrder(order); err != nil {
		logger.Error("Failed to validate order", zap.String("err", err.Error()))
		return err
	}

	err := service.AddOrder(context.Background(), order)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		logger.Error("Db insert order error", zap.String("err", err.Error()), zap.String("code", pgErr.Code))
		return err
	}

	if err != nil {
		logger.Error(err.Error())
		return err
	}

	logger.Info("Inserted", zap.String("uid", order.OrderUID))

	return nil
}

func HandleOrderMessage(logger *zap.Logger, service *orders.OrdersService) stan.MsgHandler {
	return func(m *stan.Msg) {
		handleOrderFunc(m.Data, logger, service)
	}
}
