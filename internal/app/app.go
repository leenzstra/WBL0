package app

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/leenzstra/WBL0/config"
	"github.com/leenzstra/WBL0/internal/controller/http"
	"github.com/leenzstra/WBL0/internal/controller/stan"
	"github.com/leenzstra/WBL0/internal/models"
	"github.com/leenzstra/WBL0/internal/usecases/orders"
	"github.com/leenzstra/WBL0/pkg/cache"
	"github.com/leenzstra/WBL0/pkg/database"
	"github.com/leenzstra/WBL0/pkg/logger"
	"github.com/leenzstra/WBL0/pkg/server"
	"github.com/leenzstra/WBL0/pkg/stanq"
)

func Run(config *config.Config) {
	// TODO все брать из конфига
	logger, err := logger.New("/wb.log", true)
	if err != nil {
		// TODO ???
		panic(err)
	}

	db, err := database.New(config.PgUrl, logger)
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer db.Close()

	ordersService := orders.NewService(
		orders.NewRepo(db),
		cache.New[string, models.OrderModel](30 * time.Minute),
		logger,
	)

	app := fiber.New()
	http.SetupRouter(app, logger, ordersService)
	server := server.New(app, logger)
	server.Listen(":80")

	msgHandler := stan.HandleOrderMessage(logger, ordersService)
	stanServer := stanq.New("stan", "server", config.NatsUrl, logger)
	stanServer.Start("main", msgHandler)

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