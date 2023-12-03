package validation

import (
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/leenzstra/WBL0/internal/models"
	"github.com/stretchr/testify/suite"
)

type ValidatorSuite struct {
	suite.Suite
	v *OrderValidator
}

func (suite *ValidatorSuite) SetupSuite() {
	suite.v = NewOrderValidator()
}

func TestSubSuite(t *testing.T) {
    suite.Run(t, new(ValidatorSuite))
}

func (suite *ValidatorSuite) TestNew() {
	suite.NotNil(suite.v)
}

func (suite *ValidatorSuite) TestCorrectOrder() {
	order := models.OrderModel{
		OrderUID:          "somevalue",
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
		DateCreated:       time.Now(),
		OofShard:          "somevalue",
	}
	err := suite.v.ValidateOrder(order)
	suite.Nil(err)
}

func (suite *ValidatorSuite) TestIncorrectOrder() {
		order := models.OrderModel{
		OrderUID:          "somevalue",
		TrackNumber:       "somevalue",
		Entry:             "somevalue",
		Delivery:          models.DeliveryModel{
			Name:    "somevalue",
			Phone:   "somevalue",
			Zip:     "somevalue",
			City:    "somevalue",
			Address: "somevalue",
			Region:  "somevalue",
			Email:   "somevalue",
		},
		Payment:           models.PaymentModel{
			Transaction:  "somevalue",
			RequestID:    "somevalue",
			Currency:     "somevalue",
			Provider:     "somevalue",
			Amount:       0,
			PaymentDt:    0,
			Bank:         "somevalue",
			DeliveryCost: 0,
			GoodsTotal:   0,
			CustomFee:    0,
		},
		Items:             []models.OrderItemModel{
			{ChrtID:      0,
			TrackNumber: "somevalue",
			Price:       0,
			Rid:         "somevalue",
			Name:        "somevalue",
			Sale:        0,
			Size:        "somevalue",
			TotalPrice:  0,
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
		DateCreated:       time.Now(),
		OofShard:          "somevalue",
	}
	err := suite.v.ValidateOrder(order)
	suite.Error(err, err.(validator.ValidationErrors))
}
