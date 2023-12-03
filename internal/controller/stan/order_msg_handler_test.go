package stan

import (
	"context"
	"testing"
	"time"

	"github.com/goccy/go-json"
	"github.com/leenzstra/WBL0/internal/models"
	"github.com/leenzstra/WBL0/internal/services/orders"
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

func TestAddOrderFuncOk(t *testing.T) {
	r, c, l := prepareMocks()

	data, _ := json.Marshal(order)

	r.On("Add", context.Background(), order).
		Times(1).
		Return(nil)

	c.On("SetItem", order.OrderUID, order).
		Times(1).
		Return(true)

	service := orders.NewService(r, c, l)

	err := handleOrderFunc(data, l, service)
	
	assert.NoError(t, err)
}

func TestAddOrderFuncErrorOnJson(t *testing.T) {
	r, c, l := prepareMocks()
	service := orders.NewService(r, c, l)

	data := []byte("test")

	err := handleOrderFunc(data, l, service)
	
	assert.Error(t, err)
}

func TestAddOrderFuncErrorOnValidation(t *testing.T) {
	r, c, l := prepareMocks()
	service := orders.NewService(r, c, l)

	o := models.OrderModel{}
	data, _ := json.Marshal(o)

	err := handleOrderFunc(data, l, service)
	
	assert.Error(t, err)
}