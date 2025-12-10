package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LaurensVM1/slice/internal/handlers"
	"github.com/LaurensVM1/slice/internal/models"
	"github.com/LaurensVM1/slice/internal/shop"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// Get orders
func TestGetOrders_Success(t *testing.T) {
	app := fiber.New()
	app.Get("/order", handlers.GetOrders)

	orderID := uuid.New()
	testOrder1 := models.Order{
		ID:         orderID,
		PizzaID:    shop.S.Pizzas[0].ID,
		Status:     "pending",
		TotalCents: shop.S.Pizzas[0].PriceCents,
	}

	testOrder2 := models.Order{
		ID:         orderID,
		PizzaID:    shop.S.Pizzas[1].ID,
		Status:     "done",
		TotalCents: shop.S.Pizzas[1].PriceCents,
	}

	shop.S.Mu.Lock()
	shop.S.Orders = []models.Order{testOrder1, testOrder2}
	shop.S.Mu.Unlock()

	t.Cleanup(func() {
		shop.S.Mu.Lock()
		shop.S.Orders = []models.Order{}
		shop.S.Mu.Unlock()
	})

	req := httptest.NewRequest(http.MethodGet, "/order", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)

	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200 OK")

	var result struct {
		Orders []models.Order `json:"orders"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)

	assert.Equal(t, 2, len(result.Orders))
}

// Get order
func TestGetOrder_Success(t *testing.T) {
	app := fiber.New()
	app.Get("/order/:id", handlers.GetOrder)

	orderID := uuid.New()
	testOrder := models.Order{
		ID:         orderID,
		PizzaID:    shop.S.Pizzas[0].ID,
		Status:     "pending",
		TotalCents: shop.S.Pizzas[0].PriceCents,
	}

	shop.S.Mu.Lock()
	shop.S.Orders = []models.Order{testOrder}
	shop.S.Mu.Unlock()

	t.Cleanup(func() {
		shop.S.Mu.Lock()
		shop.S.Orders = []models.Order{}
		shop.S.Mu.Unlock()
	})

	req := httptest.NewRequest(http.MethodGet, "/order/"+orderID.String(), nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)

	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200 OK")

	var result struct {
		Order models.Order `json:"order"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)

	assert.Equal(t, testOrder.ID, result.Order.ID)
	assert.Equal(t, testOrder.PizzaID, result.Order.PizzaID)
	assert.Equal(t, testOrder.Status, result.Order.Status)
	assert.Equal(t, testOrder.TotalCents, result.Order.TotalCents)
}

func TestGetOrder_NotFound(t *testing.T) {
	app := fiber.New()
	app.Get("/order/:id", handlers.GetOrder)

	orderID := uuid.New()
	testOrder := models.Order{
		ID:         orderID,
		PizzaID:    shop.S.Pizzas[0].ID,
		Status:     "pending",
		TotalCents: shop.S.Pizzas[0].PriceCents,
	}

	shop.S.Mu.Lock()
	shop.S.Orders = []models.Order{testOrder}
	shop.S.Mu.Unlock()

	t.Cleanup(func() {
		shop.S.Mu.Lock()
		shop.S.Orders = []models.Order{}
		shop.S.Mu.Unlock()
	})

	req := httptest.NewRequest(http.MethodGet, "/order/"+uuid.New().String(), nil)

	resp, _ := app.Test(req)

	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode, "Expected status code 404 Not found")
}

func TestGetOrder_InvalidID(t *testing.T) {
	app := fiber.New()
	app.Get("/order/:id", handlers.GetOrder)

	invalidID := "not-a-uuid-string"
	req := httptest.NewRequest(http.MethodGet, "/order/"+invalidID, nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Expected status code 400 Bad Request")
}

// Create order
func TestCreateOrder_Success(t *testing.T) {
	app := fiber.New()
	app.Post("/order", handlers.CreateOrder)

	orderID := uuid.New()
	testOrder := models.Order{
		ID:         orderID,
		PizzaID:    shop.S.Pizzas[0].ID,
		Status:     "pending",
		TotalCents: shop.S.Pizzas[0].PriceCents,
	}

	shop.S.Mu.Lock()
	shop.S.Orders = []models.Order{testOrder}
	shop.S.Mu.Unlock()

	t.Cleanup(func() {
		shop.S.Mu.Lock()
		shop.S.Orders = []models.Order{}
		shop.S.Mu.Unlock()
	})

	testPizza := shop.S.Pizzas[2]

	requestPayload := map[string]any{
		"pizza_id": testPizza.ID.String(),
	}

	payload, err := json.Marshal(requestPayload)
	assert.NoError(t, err, "Failed to marshal test data into JSON")

	requestBody := bytes.NewReader(payload)
	req := httptest.NewRequest(http.MethodPost, "/order", requestBody)

	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)

	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode, "Expected status code 201 CREATED")

	var result struct {
		Order models.Order `json:"order"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)

	assert.Equal(t, testPizza.ID, result.Order.PizzaID)
	assert.Equal(t, "pending", result.Order.Status)
	assert.Equal(t, testPizza.PriceCents, result.Order.TotalCents)
}

func TestCreateOrder_InvalidPizzaId(t *testing.T) {
	app := fiber.New()
	app.Post("/order", handlers.CreateOrder)

	orderID := uuid.New()
	testOrder := models.Order{
		ID:         orderID,
		PizzaID:    shop.S.Pizzas[0].ID,
		Status:     "pending",
		TotalCents: shop.S.Pizzas[0].PriceCents,
	}

	shop.S.Mu.Lock()
	shop.S.Orders = []models.Order{testOrder}
	shop.S.Mu.Unlock()

	t.Cleanup(func() {
		shop.S.Mu.Lock()
		shop.S.Orders = []models.Order{}
		shop.S.Mu.Unlock()
	})

	// Non existing pizza id
	requestPayload := map[string]any{
		"pizza_id": uuid.New().String(),
	}

	payload, err := json.Marshal(requestPayload)
	assert.NoError(t, err, "Failed to marshal test data into JSON")

	requestBody := bytes.NewReader(payload)
	req := httptest.NewRequest(http.MethodPost, "/order", requestBody)

	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)

	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Expected status code 400 Bad request")
}
