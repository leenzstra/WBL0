package orders

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/leenzstra/WBL0/internal/models"
	"github.com/leenzstra/WBL0/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var (
	order = models.OrderModel{
		OrderUID:          "some_uid",
		TrackNumber:       "somevalue",
		Entry:             "somevalue",
		Delivery:          models.DeliveryModel{
			Name:    "somevalue",
			Phone:   "somevalue",
			Zip:     "somevalue",
			City:    "somevalue",
			Address: "somevalue",
			Region:  "somevalue",
			Email:   "somevalue@somevalue.com",
		},
		Payment:           models.PaymentModel{
			Transaction:  "somevalue",
			RequestID:    "somevalue",
			Currency:     "somevalue",
			Provider:     "somevalue",
			Amount:       1,
			PaymentDt:    1,
			Bank:         "somevalue",
			DeliveryCost: 1,
			GoodsTotal:   1,
			CustomFee:    0,
		},
		Items:             []models.OrderItemModel{
			{ChrtID:      1,
			TrackNumber: "somevalue",
			Price:       1,
			Rid:         "somevalue",
			Name:        "somevalue",
			Sale:        1,
			Size:        "somevalue",
			TotalPrice:  1,
			NmID:        0,
			Brand:       "somevalue",
			Status:      0,},
		},
		Locale:            "somevalue",
		InternalSignature: "somevalue",
		CustomerID:        "somevalue",
		DeliveryService:   "somevalue",
		Shardkey:          "somevalue",
		SmID:              0,
		DateCreated:       time.Date(2020,1,1,1,1,1,1,time.UTC),
		OofShard:          "somevalue",
	}
)

type MockCache struct {
	mocks.ICache[string, models.OrderModel]
}

func prepareMocks() (*mocks.IOrderRepo, *MockCache, *zap.Logger) {
	return &mocks.IOrderRepo{}, &MockCache{}, zap.NewNop()
}

func TestNewOrderService(t *testing.T) {
	repo, cache, logger := prepareMocks()
	
	s := NewService(repo, cache, logger)

	assert.NotNil(t, s)
}

func TestAddOrder(t *testing.T) {
	repo, cache, logger := prepareMocks()

	repo.
		On("Add", context.Background(), order).
		Return(nil)

	cache.
		On("SetItem", order.OrderUID, order).
		Return(true)

	s := NewService(repo, cache, logger)
    err := s.AddOrder(context.Background(), order)

    repo.AssertCalled(t, "Add", context.Background(), order)
    cache.AssertCalled(t, "SetItem", order.OrderUID, order)
    assert.NoError(t, err)
}

func TestAddOrderErrorOnRepoAdd(t *testing.T) {
	repo, cache, logger := prepareMocks()
	repo.
		On("Add", context.Background(), order).
		Return(errors.New("some error"))

	cache.
		On("SetItem")

	s := NewService(repo, cache, logger)
    err := s.AddOrder(context.Background(), order)

    repo.AssertExpectations(t)
    cache.AssertNotCalled(t, "SetItem")
    assert.Error(t, err)
}

func TestAddOrderErrorOnSetItem(t *testing.T) {
	repo, cache, logger := prepareMocks()
	repo.
		On("Add", context.Background(), order).
		Return(nil)

	cache.
		On("SetItem", order.OrderUID, order).
		Return(false)

	s := NewService(repo, cache, logger)
    err := s.AddOrder(context.Background(), order)

    repo.AssertExpectations(t)
    cache.AssertExpectations(t)
    assert.Nil(t, err)
}

func TestGetOrderNoCache(t *testing.T) {
	repo, cache, logger := prepareMocks()
	repo.
		On("Get", context.Background(), order.OrderUID).
		Times(1).
		Return(&order, nil)

	// в кэше пусто
	cache.
		On("GetItem", order.OrderUID).
		Times(1).
		Return(models.OrderModel{}, false)

	cache.
		On("SetItem", order.OrderUID, order).
		Times(1).
		Return(true)

	s := NewService(repo, cache, logger)
    o, err := s.GetOrder(context.Background(), order.OrderUID)

    repo.AssertExpectations(t)
    cache.AssertExpectations(t)
    assert.NoError(t, err)
	assert.Equal(t, o, &order)
}

func TestGetOrderWithCache(t *testing.T) {
	repo, cache, logger := prepareMocks()
	repo.
		On("Get", context.Background(), order.OrderUID).
		Times(0)

	cache.
		On("GetItem", order.OrderUID).
		Times(1).
		Return(order, true)

	s := NewService(repo, cache, logger)
    o, err := s.GetOrder(context.Background(), order.OrderUID)
	
    repo.AssertNotCalled(t,"Get", context.Background(), order.OrderUID)
    cache.AssertExpectations(t)

    assert.NoError(t, err)
	assert.Equal(t, o, &order)
}

func TestGetOrderErrorOnRepoGet(t *testing.T) {
	repo, cache, logger := prepareMocks()
	repo.
		On("Get", context.Background(), order.OrderUID).
		Times(1).
		Return(nil, errors.New("some error"))

	cache.
		On("GetItem", order.OrderUID).
		Times(1).
		Return(models.OrderModel{}, false)

	s := NewService(repo, cache, logger)
    o, err := s.GetOrder(context.Background(), order.OrderUID)
	
    repo.AssertExpectations(t)
    cache.AssertExpectations(t)
    assert.Error(t, err)
	assert.Nil(t, o)
}

func TestGetOrderErrorOnSetItem(t *testing.T) {
	repo, cache, logger := prepareMocks()
	repo.
		On("Get", context.Background(), order.OrderUID).
		Times(1).
		Return(&order, nil)

	cache.
		On("GetItem", order.OrderUID).
		Times(1).
		Return(models.OrderModel{}, false)
	
	cache.
		On("SetItem", order.OrderUID, order).
		Times(1).
		Return(false)

	s := NewService(repo, cache, logger)
    o, err := s.GetOrder(context.Background(), order.OrderUID)
	
    repo.AssertExpectations(t)
    cache.AssertExpectations(t)
    assert.NoError(t, err)
	assert.Equal(t, o, &order)
}

func TestGetAllOrders(t *testing.T) {
	repo, cache, logger := prepareMocks()
	orders := make([]models.OrderModel, 0)
	orders = append(orders, order)

	repo.
		On("GetAll", context.Background()).
		Times(1).
		Return(orders, nil)

	s := NewService(repo, cache, logger)
    o, err := s.GetAllOrders(context.Background())
	
    repo.AssertExpectations(t)
    assert.NoError(t, err)
	assert.Equal(t, 1, len(o))
}

func TestGetAllOrdersError(t *testing.T) {
	repo, cache, logger := prepareMocks()

	repo.
		On("GetAll", context.Background()).
		Times(1).
		Return(nil, errors.New("some error"))

	s := NewService(repo, cache, logger)
    o, err := s.GetAllOrders(context.Background())
	
    repo.AssertExpectations(t)
    assert.Error(t, err)
	assert.Nil(t, o)
}

func TestRestoreCacheOk(t *testing.T) {
	repo, cache, logger := prepareMocks()
	orders := make([]models.OrderModel, 0)
	orders = append(orders, order)

	repo.
		On("GetAll", context.Background()).
		Times(1).
		Return(orders, nil)

	cache.
		On("SetItem", order.OrderUID, order).
		Times(1).
		Return(true)

	cache.
		On("Len").
		Times(1).
		Return(1)

	s := NewService(repo, cache, logger)
    err := s.RestoreCache(context.Background())
	
    repo.AssertExpectations(t)
	cache.AssertExpectations(t)
    assert.NoError(t, err)
}

func TestRestoreCacheErrorOnGetAll(t *testing.T) {
	repo, cache, logger := prepareMocks()

	repo.
		On("GetAll", context.Background()).
		Times(1).
		Return(nil, errors.New("some error"))

	s := NewService(repo, cache, logger)
    err := s.RestoreCache(context.Background())
	
    repo.AssertExpectations(t)
    assert.Error(t, err)
}

func TestRestoreCacheErrorOnSetItem(t *testing.T) {
	repo, cache, logger := prepareMocks()
	orders := make([]models.OrderModel, 0)
	orders = append(orders, order)

	repo.
		On("GetAll", context.Background()).
		Times(1).
		Return(orders, nil)

	cache.
		On("SetItem", order.OrderUID, order).
		Times(1).
		Return(false)

	s := NewService(repo, cache, logger)
    err := s.RestoreCache(context.Background())
	
    repo.AssertExpectations(t)
	cache.AssertExpectations(t)
    assert.Error(t, err)
}


