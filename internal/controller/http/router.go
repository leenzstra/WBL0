package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/leenzstra/WBL0/internal/controller/http/httpv1"
	"github.com/leenzstra/WBL0/internal/services/orders"
	"go.uber.org/zap"
)

func SetupRouter(app *fiber.App, logger *zap.Logger, service *orders.OrdersService) {
	app.Use(recover.New())
	app.Get("/metrics", monitor.New())

	app.Static("/ui", "./static")

	api := app.Group("/api")

	apiv1 := api.Group("/v1")

	ordersv1 := apiv1.Group("/orders")

	httpv1.SetupOrderRoutes(ordersv1,service,logger)
}