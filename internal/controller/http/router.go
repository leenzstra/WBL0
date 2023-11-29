package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/leenzstra/WBL0/internal/controller/http/httpv1"
	"github.com/leenzstra/WBL0/internal/usecases/orders"
	"go.uber.org/zap"
)

func SetupRouter(app *fiber.App, logger *zap.Logger, service *orders.OrdersService) {
	app.Use(recover.New())
	app.Get("/metrics", monitor.New())

	apiv1 := app.Group("/v1")
	ordersv1 := apiv1.Group("/orders")

	httpv1.SetupOrdersRoutes(ordersv1,service,logger )
}