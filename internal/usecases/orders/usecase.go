package orders

import (
	"context"
	"fmt"

	"github.com/leenzstra/WBL0/internal/models"
	"github.com/leenzstra/WBL0/pkg/cache"
	"go.uber.org/zap"
)

// TODO service or usecase?
type OrdersService struct {
	orderRepo IOrderRepo
	cache  cache.ICache[string, models.OrderModel]
	logger *zap.Logger
}

func NewService(repo IOrderRepo, ch cache.ICache[string, models.OrderModel], l *zap.Logger) *OrdersService {
	return &OrdersService{
		orderRepo: repo,
		cache: ch,
		logger: l,
	}
}

func (s *OrdersService) AddOrder(ctx context.Context, order models.OrderModel) (error) {
	err := s.orderRepo.Add(ctx, order)
	if err != nil {
		ferr := fmt.Errorf("OrdersService.AddOrder.s.orderRepo.Add: %w", err)
		s.logger.Error(ferr.Error())
		return ferr
	}

	return nil
}

func (s *OrdersService) GetOrder(ctx context.Context, uid string) (*models.OrderModel, error) {
	// берем из кэша
	cachedOrder, ok := s.cache.GetItem(uid)
	if ok {
		return &cachedOrder, nil
	}

	// берем из базы
	order, err := s.orderRepo.Get(ctx, uid)
	if err != nil {
		ferr := fmt.Errorf("OrdersService.GetOrder.s.orderRepo.Get: %w", err)
		s.logger.Error(ferr.Error())
		return nil, ferr
	}

	// записываем в кэш
	ok = s.cache.SetItem(uid, *order)
	if !ok {
		s.logger.Error("OrdersService.GetOrder.s.cache.SetItem")
	}


	return order, nil
}