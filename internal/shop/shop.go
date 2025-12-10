package shop

import (
	"sync"

	"github.com/LaurensVM1/slice/internal/models"
	"github.com/google/uuid"
)

type Shop struct {
	Mu     sync.RWMutex
	Pizzas []models.Pizza `json:"pizzas"`
	Orders []models.Order `json:"orders"`
}

var S = &Shop{
	Pizzas: []models.Pizza{
		{ID: uuid.MustParse("de1ad48f-e135-4d86-99dd-1292d3255d90"), Name: "Margherita", PriceCents: 1000},
		{ID: uuid.MustParse("80920e03-cf32-4825-833d-4a5c61fe10d0"), Name: "Pepperoni", PriceCents: 1200},
		{ID: uuid.MustParse("351c07e0-cf2c-461d-b2ee-c41adeaa40d6"), Name: "Veggie", PriceCents: 1100},
		{ID: uuid.MustParse("9879fc94-6e7a-4cbb-a979-4dcaa90946e5"), Name: "Hawaiian", PriceCents: 1300},
	},
	Orders: []models.Order{},
}
