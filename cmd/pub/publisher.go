package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/goccy/go-json"
	"github.com/leenzstra/WBL0/internal/models"

	"github.com/nats-io/stan.go"
	"go.uber.org/zap"
)

const (
	modelsFile      = "./_task/wb_l0_data.json"
	waitDurationSec = 5
	clusterId = "stan"
	clientId = "aboba_publisher"
	natsUrl = "nats://localhost:4222"
)

func readOrdersFromFile(file string) []*models.OrderModel {
	f, err := os.Open(file)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	var data []*models.OrderModel
	err = json.Unmarshal(content, &data)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	return data
}

func countUnique(orders []*models.OrderModel) int {
	counter := make(map[string]int )    
	for _, o := range orders {
		counter[o.OrderUID]++
	} 

	return len(counter)
}

func main() {
	data := readOrdersFromFile(modelsFile)

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	sc, err := stan.Connect(clusterId, clientId, stan.NatsURL(natsUrl))
	if err != nil {
		panic(err)
	}

	uc := countUnique(data)
	logger.Info("unique orders = " + fmt.Sprint(uc))

	for _, order := range data {
		orderJson, err := json.Marshal(order)
		if err != nil {
			logger.Error("Json marshal error", zap.String("err", err.Error()))
		}
		err = sc.Publish("main", orderJson)
		if err != nil {
			panic(err)
		}
		// break
		fmt.Println(order.OrderUID)
		time.Sleep(waitDurationSec * time.Microsecond)
	}

}
