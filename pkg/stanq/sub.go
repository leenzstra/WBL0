package stanq

import (
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"go.uber.org/zap"
)

type StanServer struct {
	notify chan error
	clustedId string
	clientId string
	natsUrl string
	logger *zap.Logger
	conn stan.Conn
	sub stan.Subscription
}

func New(clustedId, clientId, natsUrl string, logger *zap.Logger) *StanServer {
	s := &StanServer{
		notify: make(chan error, 1),
		clustedId:    clustedId,
		clientId: clientId,
		natsUrl: natsUrl,
		logger: logger,
	}

	return s
}

func (s *StanServer) Start(topic string, handler stan.MsgHandler) {
	conn, err := stan.Connect(s.clustedId, s.clientId, stan.NatsURL(s.natsUrl))
	s.conn = conn
	if err != nil {
		s.notify <- err
		return
	}

	sub, err := conn.Subscribe(topic,handler,stan.StartWithLastReceived())
	s.sub = sub
	if err != nil {
		s.notify <- err
		return
	}

	s.conn.NatsConn().SetDisconnectErrHandler(func(*nats.Conn, error) {
		s.notify <- err
	})
	
	s.logger.Info("Stan started")
}

func (s *StanServer) Notify() <-chan error {
	return s.notify
}

func (s *StanServer) Shutdown() error {
	err := s.conn.Close()
	if err != nil {
		return err
	}
	return s.conn.Close()
}
