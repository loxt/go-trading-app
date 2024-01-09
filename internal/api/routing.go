package api

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func AddRoutes(app *fiber.App) {
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}

		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/trades/:ticker", websocket.New(ListenTicker))

	app.Route("/api/v1", func(router fiber.Router) {
		router.Get("/tickers", GetAllTickers)
	})
}
