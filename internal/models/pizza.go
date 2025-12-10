package models

import "github.com/google/uuid"

type Pizza struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	PriceCents int       `json:"price_cents"`
}
