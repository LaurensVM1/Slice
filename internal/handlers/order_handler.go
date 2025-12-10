package handlers

import (
	"time"

	"github.com/LaurensVM1/slice/internal/models"
	"github.com/LaurensVM1/slice/internal/shop"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetOrders(c *fiber.Ctx) error {
	shop.S.Mu.RLock()
	defer shop.S.Mu.RUnlock()

	return c.JSON(fiber.Map{
		"orders": shop.S.Orders,
	})
}

func GetOrder(c *fiber.Ctx) error {
	id := c.Params("id")

	orderID, err := uuid.Parse(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid order id"})
	}

	shop.S.Mu.RLock()
	defer shop.S.Mu.RUnlock()

	for _, o := range shop.S.Orders {
		if o.ID == orderID {
			return c.JSON(fiber.Map{
				"order": o,
			})
		}
	}

	return c.Status(404).JSON(fiber.Map{"error": "order not found"})
}

func CreateOrder(c *fiber.Ctx) error {
	var req struct {
		PizzaID uuid.UUID `json:"pizza_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	shop.S.Mu.RLock()

	var pizza *models.Pizza
	for _, p := range shop.S.Pizzas {
		if p.ID == req.PizzaID {
			pizza = &p
		}
	}

	shop.S.Mu.RUnlock()

	if pizza == nil {
		return c.Status(400).JSON(fiber.Map{"error": "pizza not found"})
	}

	var order = models.Order{
		ID:         uuid.New(),
		PizzaID:    pizza.ID,
		Status:     "pending",
		TotalCents: pizza.PriceCents,
	}

	shop.S.Mu.RLock()
	shop.S.Orders = append(shop.S.Orders, order)
	shop.S.Mu.RUnlock()

	go func(orderID uuid.UUID) {
		time.Sleep(5 * time.Second)

		shop.S.Mu.Lock()
		defer shop.S.Mu.Unlock()

		for i := range shop.S.Orders {
			if shop.S.Orders[i].ID == orderID {
				shop.S.Orders[i].Status = "done"
				break
			}
		}
	}(order.ID)

	return c.Status(201).JSON(fiber.Map{
		"order": order,
	})
}
