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

type orderResponse struct {
	Order models.OrderModel `json:"order"`
}

func SetupOrderRoutes(group fiber.Router, service *orders.OrdersService, logger *zap.Logger) {
	h := &ordersHandler{service, logger}

	group.Get("/:uid", h.order)
}

func (h *ordersHandler) order(c *fiber.Ctx) error {
	uid := c.Params("uid")
	if uid == "" {
		c.JSON(response.Error[any]("no uid"))
		return c.SendStatus(400)
	}

	order, err := h.ordersService.GetOrder(context.Background(), uid)
	if err != nil {
		c.JSON(response.Error[any]("get order error"))
		return c.SendStatus(400)
	}

	r := orderResponse{*order}
	return c.JSON(response.Ok[orderResponse](r))
}

