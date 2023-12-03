package stanq

import (
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"go.uber.org/zap"
)

type StanConsumer struct {
	notify chan error
	clusterId string
	clientId string
	natsUrl string
	logger *zap.Logger
	conn stan.Conn
	sub stan.Subscription
}

func NewConsumer(clusterId, clientId, natsUrl string, logger *zap.Logger) *StanConsumer {
	s := &StanConsumer{
		notify: make(chan error, 1),
		clusterId:    clusterId,
		clientId: clientId,
		natsUrl: natsUrl,
		logger: logger,
	}

	return s
}

func (s *StanConsumer) Start(topic string, handler stan.MsgHandler) error {
	conn, err := stan.Connect(s.clusterId, s.clientId, stan.NatsURL(s.natsUrl))
	if err != nil {
		return err
	}
	s.conn = conn

	sub, err := conn.Subscribe(topic,handler,stan.StartWithLastReceived(), stan.DurableName("durable."+topic))
	if err != nil {
		return err
	}
	s.sub = sub

	s.conn.NatsConn().SetDisconnectErrHandler(func(*nats.Conn, error) {
		s.notify <- err
	})
	
	s.logger.Info("Stan started")
	
	return nil
}

func (s *StanConsumer) Notify() <-chan error {
	return s.notify
}

func (s *StanConsumer) Shutdown() error {
	return s.conn.Close()
}
