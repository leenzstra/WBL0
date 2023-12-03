package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/leenzstra/WBL0/config"
	"github.com/leenzstra/WBL0/internal/controller/http"
	"github.com/leenzstra/WBL0/internal/controller/stan"
	"github.com/leenzstra/WBL0/internal/models"
	"github.com/leenzstra/WBL0/internal/services/orders"
	"github.com/leenzstra/WBL0/pkg/cache"
	"github.com/leenzstra/WBL0/pkg/database"
	"github.com/leenzstra/WBL0/pkg/logger"
	"github.com/leenzstra/WBL0/pkg/server"
	"github.com/leenzstra/WBL0/pkg/stanq"
	"go.uber.org/zap"
)

const (
	logsFile = "/wb.log"
	addr = ":80"
	consumerId = "consumer-1"
)

func setupDatabase(url string, l *zap.Logger) (*database.DB) {
	// создаем pgxpool (интерфейс IPool)
	pool, err := database.NewPgxPool(url)
	if err != nil {
		l.Fatal(err.Error())
	}

	// создаем pgxscan (интерфейс IScanner)
	scanner, err := database.NewPgxScanner()
	if err != nil {
		l.Fatal(err.Error())
	}

	// база данных заказов
	db, err := database.New(pool, database.PgxBuilder, scanner, l)
	if err != nil {
		l.Fatal(err.Error())
	}

	return db
}

func Run(config *config.Config) {
	logger, err := logger.New(logsFile, config.Debug)
	if err != nil {
		panic(err)
	}

	db := setupDatabase(config.PgUrl, logger)
	defer db.Close()

	// сервис работы с заказами
	ordersService := orders.NewService(
		orders.NewRepo(db),
		cache.New[string, models.OrderModel](24 * time.Hour),
		logger,
	)

	// восстановление кэша из БД
	err = ordersService.RestoreCache(context.Background())
	if err != nil {
		logger.Fatal(err.Error())
	}

	// http сервер
	app := fiber.New()
	http.SetupRouter(app, logger, ordersService)
	server := server.New(app, logger)
	server.Listen(addr)

	// stan consumer
	msgHandler := stan.HandleOrderMessage(logger, ordersService)
	stanServer := stanq.NewConsumer(config.ClusterId, consumerId, config.NatsUrl, logger)
	err = stanServer.Start(config.Topic, msgHandler)
	if err != nil {
		logger.Fatal(err.Error(), zap.String("cluster", config.ClusterId))
	}

	logger.Info("All services ready")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	select {
	case s := <-interrupt:
		logger.Info("interrupted" + s.String())
	case err = <-server.Notify():
		logger.Error("err interrupted" + err.Error())
	case err = <-stanServer.Notify():
		logger.Error("err interrupted" + err.Error())
	}

	err = server.Shutdown()
	if err != nil {
		logger.Error(err.Error())
	}

	err = stanServer.Shutdown()
	if err != nil {
		logger.Error(err.Error())
	}
}