package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigNew(t *testing.T) {
	pgUrl := "postgres://localhost:5432/db"
	natsUrl := "nats://localhost:4222"
	clusterId := "cluster-1"
	topic := "topic-1"

	type args struct {
		pgUrl     string
		natsUrl   string
		clusterId string
		topic     string
		debug     string
	}

	tests := []struct {
		name string
		args args
		want Config
	}{
		{
			name: "test true debug value",
			args: args{
				pgUrl:     pgUrl,
				natsUrl:   natsUrl,
				clusterId: clusterId,
				topic:     topic,
				debug:     "deBUg",
			},
			want: Config{
				PgUrl:     pgUrl,
				NatsUrl:   natsUrl,
				ClusterId: clusterId,
				Topic:     topic,
				Debug:     true,
			},
		},
		{
			name: "test false debug value",
			args: args{
				pgUrl:     pgUrl,
				natsUrl:   natsUrl,
				clusterId: clusterId,
				topic:     topic,
				debug:     "prod",
			},
			want: Config{
				PgUrl:     pgUrl,
				NatsUrl:   natsUrl,
				ClusterId: clusterId,
				Topic:     topic,
				Debug:     false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.args.pgUrl, tt.args.natsUrl, tt.args.clusterId, tt.args.topic, tt.args.debug)
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("New() = %v, want %v", got, tt.want)
			// }
			assert.Equal(t, got, tt.want)
			
		})
	}
}
