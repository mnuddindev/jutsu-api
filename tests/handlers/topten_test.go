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

func TestTopTenHandler_GetTopTen(t *testing.T) {
	app := fiber.New()
	h := handler.NewTopTenHandler()
	app.Get("/api/top-ten", h.GetTopTen)

	t.Run("Success - Get top ten anime", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/top-ten", nil)
		resp, err := app.Test(req, -1)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		// May succeed or fail depending on external service
		assert.Contains(t, []int{http.StatusOK, http.StatusBadGateway}, resp.StatusCode)

		if resp.StatusCode == http.StatusOK {
			assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
			var result map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&result)
			require.NoError(t, err)
			assert.True(t, result["success"].(bool))
			assert.NotNil(t, result["results"])
		}
	})
}
