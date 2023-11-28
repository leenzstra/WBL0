package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/leenzstra/WBL0/config"
	"github.com/leenzstra/WBL0/db"
	"github.com/leenzstra/WBL0/internal/repository"
	"github.com/nats-io/stan.go"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Printf("Logger failed: %s", err)
	}

	// TODO интерфейс для подключения / ORM
	config := config.New(fmt.Sprintf("postgres://%s:%s@%s:%s/%s", os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), "db", os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_DB")))

	conn, err := pgx.Connect(context.Background(), config.ConnUrl)
	if err != nil {
		logger.Fatal("Pgx connect failed", zap.String("err", err.Error()))
	}

	repo := repository.NewOrderRepository(conn, logger)
	_ = repo

	stanConn, err := stan.Connect("stan", "server", stan.NatsURL(os.Getenv("NATS_URL")))
	if err != nil {
		logger.Fatal("Stan connect failed", zap.String("err", err.Error()))
	}

	sub, err := stanConn.Subscribe("main",
		func(m *stan.Msg) {
			var order db.OrderModel
			err := json.Unmarshal(m.Data, &order)
			if err != nil {
				logger.Error("Failed to unmarshal order", zap.String("err", err.Error()))
				return
			}

			err = repo.Insert(context.Background(), &order)

			// TODO сделать нормальную обработку ошибок
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				logger.Error("Db insert order error", zap.String("err", err.Error()), zap.String("code", pgErr.Code))
			}

			logger.Info("Inserted", zap.String("uid", order.OrderUID))
		},
		stan.StartWithLastReceived())

	if err != nil {
		logger.Fatal("Sub failed", zap.String("err", err.Error()))
	}

	logger.Info("Started successfully")

	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)

	go func() {
		for range signalChan {
			logger.Info("Gracefully shutdown")
			sub.Unsubscribe()
			stanConn.Close()
			conn.Close(context.Background())
			cleanupDone <- true
		}
	}()

	// <-cleanupDone

	// TODO а теперь сюда Fiber
	http.HandleFunc("/order", getOrder)

	if err := http.ListenAndServe(":80", nil); err != nil {
		logger.Fatal(err.Error())
	}

	<-cleanupDone

}

func getOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		uid := r.URL.Query().Get("uid")
		io.WriteString(w, uid)
	}
}