package main

import (
	"log"

	"github.com/LaurensVM1/slice/internal/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

const port = ":8080"

func main() {
	app := fiber.New()

	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	api := app.Group("/api")

	// Health
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Slice API is running",
		})
	})

	// Menu
	api.Get("/menu", handlers.GetMenu)

	// Order
	api.Get("/order", handlers.GetOrders)
	api.Get("/order/:id", handlers.GetOrder)
	api.Post("/order", handlers.CreateOrder)

	if err := app.Listen(port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
