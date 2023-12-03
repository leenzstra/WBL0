package orders

import (
	"context"
	"fmt"

	"github.com/leenzstra/WBL0/internal/models"
	"github.com/leenzstra/WBL0/pkg/cache"
	"go.uber.org/zap"
)

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

	// записываем в кэш
	ok := s.cache.SetItem(order.OrderUID, order)
	if !ok {
		s.logger.Error("OrdersService.AddOrder.s.cache.SetItem")
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

func (s *OrdersService) GetAllOrders(ctx context.Context) ([]models.OrderModel, error) {
	orders, err := s.orderRepo.GetAll(ctx)
	if err != nil {
		ferr := fmt.Errorf("OrdersService.GetAllOrders: %w", err)
		s.logger.Error(ferr.Error())
		return nil, ferr
	}

	return orders, nil
}

func (s *OrdersService) RestoreCache(ctx context.Context) error {
	orders, err := s.GetAllOrders(ctx)
	if err != nil {
		ferr := fmt.Errorf("RestoreCache.GetAllOrders: %w", err)
		s.logger.Error(ferr.Error())
		return ferr
	}

	r := true
	for _, o := range orders {
		r = r && s.cache.SetItem(o.OrderUID, o)
	}

	if !r {
		ferr := fmt.Errorf("RestoreCache.SetItem: !ok")
		s.logger.Error(ferr.Error())
		return ferr
	}

	s.logger.Debug("Cache restored: ", zap.Int("orders", s.cache.Len()))

	return nil
}

