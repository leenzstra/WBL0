package stanq

import (
	"errors"
	"fmt"
	"testing"

	"github.com/nats-io/nats-streaming-server/server"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

const testPort = 4222
const clusterName = "test"
const clientName = "tester"

func runServer(ID string) (*server.StanServer, error) {
	return server.RunServer(ID)
}

type PubSuite struct {
    suite.Suite
    ss *server.StanServer
}

func (suite *PubSuite) SetupTest() {
	var s *server.StanServer

	s, err := runServer(clusterName)
	suite.Nil(err)

	suite.ss = s
}

func (suite *PubSuite) TearDownTest() {
	suite.ss.Shutdown()
}

func TestPubSuite(t *testing.T) {
    suite.Run(t, new(PubSuite))
}

func (suite *PubSuite) TestStart() {
	logger := zap.NewNop()
	suite.NotNil(suite.ss)
	
	pub := NewPublisher(clusterName, clientName, fmt.Sprintf("nats://localhost:%d", testPort), logger)
	err := pub.Start("test_topic")

	suite.Nil(err)
}

func (suite *PubSuite) TestStartError() {
	logger := zap.NewNop()
	
	pub := NewPublisher("random_cluster", clientName, fmt.Sprintf("nats://localhost:%d", testPort), logger)
	err := pub.Start("test_topic")

	suite.Error(err)
}

func (suite *PubSuite) TestShutdown() {
	logger := zap.NewNop()
	
	pub := NewPublisher(clusterName, clientName, fmt.Sprintf("nats://localhost:%d", testPort), logger)
	err := pub.Start("test_topic")
	suite.Nil(err)

	err = pub.Shutdown()
	suite.Nil(err)
}

func (suite *PubSuite) TestNotifyError() {
	logger := zap.NewNop()
	
	pub := NewPublisher(clusterName, clientName, fmt.Sprintf("nats://localhost:%d", testPort), logger)
	err := pub.Start("test_topic")
	suite.Nil(err)

	pub.notify <- errors.New("test error")
	err = <- pub.Notify()
	suite.Error(err)
}
