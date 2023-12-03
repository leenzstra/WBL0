package stanq

import (
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"go.uber.org/zap"
)

type StanPublisher struct {
	notify chan error
	clusterId string
	clientId string
	natsUrl string
	logger *zap.Logger
	stan.Conn
}

func NewPublisher(clusterId, clientId, natsUrl string, logger *zap.Logger) *StanPublisher {
	s := &StanPublisher{
		notify: make(chan error, 1),
		clusterId:    clusterId,
		clientId: clientId,
		natsUrl: natsUrl,
		logger: logger,
	}

	return s
}

func (s *StanPublisher) Start(topic string) error {
	conn, err := stan.Connect(s.clusterId, s.clientId, stan.NatsURL(s.natsUrl))
	if err != nil {
		return err
	}
	s.Conn = conn

	s.NatsConn().SetDisconnectErrHandler(func(*nats.Conn, error) {
		s.notify <- err
	})
		
	s.logger.Info("Publisher started")

	return nil
}

func (s *StanPublisher) Notify() <-chan error {
	return s.notify
}

func (s *StanPublisher) Shutdown() error {
	return s.Close()
}
