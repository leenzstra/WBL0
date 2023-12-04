package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/goccy/go-json"
	"github.com/leenzstra/WBL0/internal/models"
	"github.com/leenzstra/WBL0/pkg/logger"
	"github.com/leenzstra/WBL0/pkg/stanq"
	"go.uber.org/zap"
)

const (
	modelsFile      = "./_task/wb_l0_data_1000.json"
	waitDurationSec = 5
	clusterId = "stan"
	topic = "orders"
	publisherId = "aboba_publisher"
	natsUrl = "nats://localhost:4222"
	logFile = "pub.log"
	debug = true
)

func readOrdersFromFile(file string) ([]*models.OrderModel, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, errors.New("Error opening file: " + err.Error())
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		return nil, errors.New("Error reading file: " + err.Error())
	}

	var data []*models.OrderModel
	err = json.Unmarshal(content, &data)
	if err != nil {
		return nil, errors.New("Error Unmarshal() " + err.Error())
	}

	return data, nil
}

func countUnique(orders []*models.OrderModel) int {
	counter := make(map[string]int )    
	for _, o := range orders {
		counter[o.OrderUID]++
	} 

	return len(counter)
}

func sendOrders(data []*models.OrderModel, pub *stanq.StanPublisher, logger *zap.Logger, ch chan bool) {
	for i, order := range data {
		orderJson, err := json.Marshal(order)
		if err != nil {
			logger.Error("json marshal error" + err.Error())
			break
		}
		err = pub.Publish(topic, orderJson)
		if err != nil {
			logger.Error("publish error" + err.Error())
			break
		}
		logger.Debug("published", zap.String("uid", order.OrderUID), zap.Int("idx", i))
		time.Sleep(waitDurationSec * time.Microsecond)
	}
	ch <- true
}

func main() {
	logger, err := logger.New(logFile, debug)
	if err != nil {
		panic(err)
	}

	data, err := readOrdersFromFile(modelsFile)
	if err != nil {
		logger.Fatal(err.Error())
	}

	pub := stanq.NewPublisher(clusterId, publisherId, natsUrl, logger)
	err = pub.Start(topic)
	if err != nil {
		logger.Fatal(err.Error())
	}

	uc := countUnique(data)
	logger.Info("unique orders = " + fmt.Sprint(uc))

	interrupt := make(chan os.Signal, 1)
	finished := make(chan bool, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	go sendOrders(data, pub, logger, finished)

	select {
	case result := <-finished:
		logger.Info("process finished with status: " + fmt.Sprint(result))
	case s := <-interrupt:
		logger.Info("interrupted" + s.String())
	case err = <-pub.Notify():
		logger.Error("err interrupted" + err.Error())
	}

	if err := pub.Close(); err != nil {
		logger.Error(err.Error())
	}

}
