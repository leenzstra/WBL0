package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"time"

	"github.com/leenzstra/WBL0/db"
	"github.com/nats-io/stan.go"
	"go.uber.org/zap"
)

const (
	modelsFile      = "wb_l0_data.json"
	waitDurationSec = 5
)

func readOrdersFromFile(file string) []*db.OrderModel {
	f, err := os.Open(file)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	var data []*db.OrderModel
	err = json.Unmarshal(content, &data)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	return data
}

func main() {
	data := readOrdersFromFile(modelsFile)

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	sc, err := stan.Connect("stan", "aboba_publisher", stan.NatsURL("nats://localhost:4222"))
	if err != nil {
		panic(err)
	}

	for _, order := range data {
		orderJson, err := json.Marshal(order)
		if err != nil {
			logger.Error("Json marshal error", zap.String("err", err.Error()))
		}
		err = sc.Publish("main", orderJson)
		if err != nil {
			panic(err)
		}
		break
		time.Sleep(waitDurationSec * time.Second)
	}

}
