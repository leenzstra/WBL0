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

func HandleOrderMessage(logger *zap.Logger, service *orders.OrdersService) stan.MsgHandler {
	return func(m *stan.Msg) {
		var order models.OrderModel
		if err := json.Unmarshal(m.Data, &order); err != nil {
			logger.Error("Failed to unmarshal order", zap.String("err", err.Error()))
			return
		}

		// TODO сделать один инстанс
		v := validation.NewOrderValidator()
		if err := v.ValidateOrder(order); err != nil {
			logger.Error("Failed to validate order", zap.String("err", err.Error()))
			return
		}
		
		err := service.AddOrder(context.Background(), order)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			logger.Error("Db insert order error", zap.String("err", err.Error()), zap.String("code", pgErr.Code))
			return
		}

		if err != nil {
			logger.Error(err.Error())
			return
		}

		logger.Info("Inserted", zap.String("uid", order.OrderUID))
	}
}