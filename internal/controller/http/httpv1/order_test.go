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
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type MockCache struct {
	mocks.ICache[string, models.OrderModel]
}

func prepareMocks() (*mocks.IOrderRepo, *MockCache, *zap.Logger) {
	return &mocks.IOrderRepo{}, &MockCache{}, zap.NewNop()
}

type OrderRouteSuite struct {
	suite.Suite
	app *fiber.App
	cache *MockCache
	repo *mocks.IOrderRepo
	service *orders.OrdersService
}

func (suite *OrderRouteSuite) SetupSuite() {
	r, c, l := prepareMocks()
	suite.cache = c
	suite.repo = r
	suite.service = orders.NewService(r, c, l)
	suite.app = fiber.New()
	SetupOrderRoutes(suite.app, suite.service, l)
}


func TestSuite(t *testing.T) {
	suite.Run(t, new(OrderRouteSuite))
}

func (suite *OrderRouteSuite) TestOrderRouteOk() {
	uid := "some_uid"

	suite.cache.On("GetItem", uid).
		Times(1).
		Return(models.OrderModel{}, true)

	req := httptest.NewRequest("GET", "/"+uid, nil)
	resp, err := suite.app.Test(req, 1)
	suite.Equal(200, resp.StatusCode, "ok")
	suite.NoError(err)
}

func (suite *OrderRouteSuite) TestOrderRouterNoSuchRoute() {
	req := httptest.NewRequest("GET", "/", nil)
	resp, err := suite.app.Test(req, 1)
	suite.Equal(404, resp.StatusCode, "no such route")
	suite.NoError(err)
}

func (suite *OrderRouteSuite) TestOrderRouterOtherError() {
	uid := "nosuchuid"

	suite.cache.On("GetItem", uid).
		Times(1).
		Return(models.OrderModel{}, false)

	suite.repo.On("Get", context.Background(), uid).
		Times(1).
		Return(nil, errors.New("no such uid"))

	req := httptest.NewRequest("GET", "/"+uid, nil)
	resp, err := suite.app.Test(req, 1)
	suite.Equal(400, resp.StatusCode, "no such uid")
	suite.NoError(err)
}