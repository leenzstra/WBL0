package config

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	pgUrl := "postgres://localhost:5432/db"
    natsUrl := "nats://localhost:4222"
    clusterId := "cluster-1"
    topic := "topic-1"

	type Args struct {
		pgUrl     string
		natsUrl   string
		clusterId string
		topic     string
		debug     string
	}

	tests := []struct {
		name string
		args Args
		want Config
	}{
		{name: "test debug value",
		args: Args{pgUrl, natsUrl, clusterId, topic, "deBUg"},
		want: Config{pgUrl, natsUrl, clusterId, topic, true},},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.pgUrl, tt.args.natsUrl, tt.args.clusterId, tt.args.topic, tt.args.debug); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
