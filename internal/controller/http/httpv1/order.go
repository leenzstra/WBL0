package httpv1

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/leenzstra/WBL0/internal/models"
	"github.com/leenzstra/WBL0/internal/models/response"
	"github.com/leenzstra/WBL0/internal/services/orders"
	"go.uber.org/zap"
)

type ordersHandler struct {
	ordersService *orders.OrdersService
	logger *zap.Logger
}

func SetupOrdersRoutes(group fiber.Router, service *orders.OrdersService, logger *zap.Logger) {
	h := &ordersHandler{service, logger}

	// routes
	group.Get("/:uid", h.order)
}

type orderResponse struct {
	Order models.OrderModel `json:"order"`
}

func (h *ordersHandler) order(c *fiber.Ctx) error {
	uid := c.Params("uid")
	if uid == "" {
		return c.JSON(response.Error[any]("no uid"))
	}

	order, err := h.ordersService.GetOrder(context.Background(), uid)
	if err != nil {
		return c.JSON(response.Error[any]("get order error"))
	}

	r := orderResponse{*order}

	return c.JSON(response.Ok[orderResponse](r))
}

