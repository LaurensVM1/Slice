package handlers

import (
	"github.com/LaurensVM1/slice/internal/shop"
	"github.com/gofiber/fiber/v2"
)

func GetMenu(c *fiber.Ctx) error {
	shop.S.Mu.RLock()
	defer shop.S.Mu.RUnlock()

	return c.JSON(fiber.Map{
		"menu": shop.S.Pizzas,
	})
}
