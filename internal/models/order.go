package models

import "github.com/google/uuid"

type Order struct {
	ID         uuid.UUID `json:"id"`
	PizzaID    uuid.UUID `json:"pizza_id"`
	Status     string    `json:"status"`
	TotalCents int       `json:"total_cents"`
}
