package validation

import "github.com/leenzstra/WBL0/internal/models"
import "github.com/go-playground/validator/v10"

type OrderValidator struct {
	*validator.Validate
}

func NewOrderValidator(options ...validator.Option) *OrderValidator {
	return &OrderValidator{
		validator.New(append(options, validator.WithRequiredStructEnabled())...),
	}
}

func (v *OrderValidator) ValidateOrder(order models.OrderModel) error {
	err :=v.Struct(order)
	return err
}