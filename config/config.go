package config

import (
	"strings"
)

type Config struct {
	PgUrl     string
	NatsUrl   string
	ClusterId string
	Topic     string
	Debug     bool
}

func New(pgUrl, natsUrl, clusterId, topic, debug string) Config {
	return Config{
		PgUrl:     pgUrl,
		NatsUrl:   natsUrl,
		ClusterId: clusterId,
		Topic:     topic,
		Debug:     strings.ToLower(debug) == "debug",
	}
}