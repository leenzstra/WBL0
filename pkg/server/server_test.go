package server

import (
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestServerStarts(t *testing.T) {
    app := fiber.New()
    logger := zap.NewNop()

    server := New(app, logger)

    assert.NotNil(t, server)
}

func TestServerListen(t *testing.T) {
    app := fiber.New()
    logger := zap.NewNop()

    app.Get("/", func (c *fiber.Ctx) error  {
        return c.SendString("ok")
    })
    server := New(app, logger)

    server.Listen(":5000")

    time.Sleep(100 * time.Millisecond)

    resp, err := http.Get("http://localhost:5000/")
    assert.NoError(t, err)
    assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestServerErrorAddr(t *testing.T) {
    app := fiber.New()
    logger := zap.NewNop()

    server := New(app, logger)

    server.Listen("invalid_address")
    err := <- server.Notify()
    assert.Error(t, err)
}

func TestServerShutdown(t *testing.T) {
    app := fiber.New()
    logger := zap.NewNop()

    server := New(app, logger)

    err := server.Shutdown()

    assert.Nil(t, err)
}