package stanq

import (
	"errors"
	"fmt"
	"testing"

	"github.com/nats-io/nats-streaming-server/server"
	"github.com/nats-io/stan.go"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type SubSuite struct {
    suite.Suite
    ss *server.StanServer
}

func (suite *SubSuite) SetupTest() {
	var s *server.StanServer

	s, err := runServer(clusterName)
	suite.Nil(err)

	suite.ss = s
}

func (suite *SubSuite) TearDownTest() {
	suite.ss.Shutdown()
}

func TestSubSuite(t *testing.T) {
    suite.Run(t, new(SubSuite))
}

func (suite *SubSuite) TestStart() {
	logger := zap.NewNop()
	suite.NotNil(suite.ss)
	
	sub := NewConsumer(clusterName, clientName, fmt.Sprintf("nats://localhost:%d", testPort), logger)
	err := sub.Start("test_topic", func(msg *stan.Msg) {})

	suite.Nil(err)
}

func (suite *SubSuite) TestStartError() {
	logger := zap.NewNop()
	
	sub := NewConsumer("random cluster", clientName, fmt.Sprintf("nats://localhost:%d", testPort), logger)
	err := sub.Start("test_topic", func(msg *stan.Msg) {})

	suite.Error(err)
}

func (suite *SubSuite) TestShutdown() {
	logger := zap.NewNop()
	
	sub := NewConsumer(clusterName, clientName, fmt.Sprintf("nats://localhost:%d", testPort), logger)
	err := sub.Start("test_topic", func(msg *stan.Msg) {})
	suite.Nil(err)

	err = sub.Shutdown()
	suite.Nil(err)
}

func (suite *SubSuite) TestNotifyError() {
	logger := zap.NewNop()
	
	pub := NewPublisher(clusterName, clientName, fmt.Sprintf("nats://localhost:%d", testPort), logger)
	err := pub.Start("test_topic")
	suite.Nil(err)

	pub.notify <- errors.New("test error")
	err = <- pub.Notify()
	suite.Error(err)
}
