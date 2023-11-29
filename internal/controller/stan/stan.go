package stan

import (
	"context"
	"errors"
	"fmt"

	"github.com/goccy/go-json"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/leenzstra/WBL0/internal/models"
	"github.com/leenzstra/WBL0/internal/usecases/orders"
	"github.com/nats-io/stan.go"
	"go.uber.org/zap"
)

func validateMessage() {
	
}

func HandleOrderMessage(logger *zap.Logger, service *orders.OrdersService) stan.MsgHandler {
	return func(m *stan.Msg) {
		var order models.OrderModel
		err := json.Unmarshal(m.Data, &order)
		if err != nil {
			logger.Error("Failed to unmarshal order", zap.String("err", err.Error()))
			return
		}
		fmt.Printf("%+v\n", order)
		
		err = service.AddOrder(context.Background(), order)

		// TODO сделать нормальную обработку ошибок
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			logger.Error("Db insert order error", zap.String("err", err.Error()), zap.String("code", pgErr.Code))
		}

		logger.Info("Inserted", zap.String("uid", order.OrderUID))
	}
}