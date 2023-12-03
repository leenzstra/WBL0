package httpv1

import (
	"context"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/leenzstra/WBL0/internal/models"
	"github.com/leenzstra/WBL0/internal/services/orders"
	"github.com/leenzstra/WBL0/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type MockCache struct {
	mocks.ICache[string, models.OrderModel]
}

func prepareMocks() (*mocks.IOrderRepo, *MockCache, *zap.Logger) {
	return &mocks.IOrderRepo{}, &MockCache{}, zap.NewNop()
}

func TestOrderRouteOk(t *testing.T) {
	r, c, l := prepareMocks()
	service := orders.NewService(r, c, l)
	uid := "some_uid"

	c.On("GetItem", uid).
		Times(1).
		Return(models.OrderModel{}, true)

	app := fiber.New()

	SetupOrderRoutes(app, service, l)

	req := httptest.NewRequest("GET", "/"+uid, nil)
	resp, _ := app.Test(req, 1)
	assert.Equal(t, 200, resp.StatusCode, "ok")
}

func TestOrderRouterNoSuchRoute(t *testing.T) {
	r, c, l := prepareMocks()
	service := orders.NewService(r, c, l)

	app := fiber.New()

	SetupOrderRoutes(app, service, l)

	req := httptest.NewRequest("GET", "/", nil)
	resp, _ := app.Test(req, 1)
	assert.Equal(t, 404, resp.StatusCode, "no such route")
}

func TestOrderRouterOtherError(t *testing.T) {
	r, c, l := prepareMocks()
	service := orders.NewService(r, c, l)
	uid := "nosuchuid"

	c.On("GetItem", uid).
		Times(1).
		Return(models.OrderModel{}, false)

	r.On("Get", context.Background(), uid).
		Times(1).
		Return(nil, errors.New("no such uid"))

	app := fiber.New()

	SetupOrderRoutes(app, service, l)

	req := httptest.NewRequest("GET", "/"+uid, nil)
	resp, _ := app.Test(req, 1)
	assert.Equal(t, 400, resp.StatusCode, "no such uid")
}