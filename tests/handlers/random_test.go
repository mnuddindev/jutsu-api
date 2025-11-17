package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mnuddindev/jutsu-api/internal/interface/http/handler"
)

func TestRandomHandler_GetRandom(t *testing.T) {
	app := fiber.New()
	h := handler.NewRandomHandler()
	app.Get("/api/random", h.GetRandom)

	t.Run("Success - Get random anime", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/random", nil)
		resp, err := app.Test(req, -1)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		// May succeed or fail depending on external service
		assert.Contains(t, []int{http.StatusOK, http.StatusBadGateway, http.StatusServiceUnavailable}, resp.StatusCode)

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			return
		}

		if resp.StatusCode == http.StatusOK {
			assert.True(t, result["success"].(bool))
			assert.NotNil(t, result["results"])
		}
	})
}

func TestRandomHandler_GetRandomID(t *testing.T) {
	app := fiber.New()
	h := handler.NewRandomHandler()
	app.Get("/api/random/id", h.GetRandomID)

	t.Run("Success - Get random anime ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/random/id", nil)
		resp, err := app.Test(req, -1)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		// May succeed or fail depending on external service
		assert.Contains(t, []int{http.StatusOK, http.StatusBadGateway, http.StatusServiceUnavailable}, resp.StatusCode)

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			return
		}

		if resp.StatusCode == http.StatusOK {
			assert.True(t, result["success"].(bool))
			assert.NotNil(t, result["results"])
			results := result["results"].(map[string]interface{})
			assert.NotEmpty(t, results["id"])
		}
	})
}
