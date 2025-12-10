package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LaurensVM1/slice/internal/handlers"
	"github.com/LaurensVM1/slice/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

// Get menu
func TestGetMenu_Success(t *testing.T) {
	app := fiber.New()
	app.Get("/menu", handlers.GetMenu)

	req := httptest.NewRequest(http.MethodGet, "/menu", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)

	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200 OK")

	var result struct {
		Menu []models.Pizza `json:"menu"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)

	assert.Equal(t, 4, len(result.Menu))
}
