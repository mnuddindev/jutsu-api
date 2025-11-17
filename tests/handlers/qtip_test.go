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

func TestQtipHandler_GetQtip(t *testing.T) {
	app := fiber.New()
	h := handler.NewQtipHandler()
	app.Get("/api/qtip/:id", h.GetQtip)

	testCases := []struct {
		name           string
		path           string
		expectedStatus int
		validateResult func(t *testing.T, result map[string]interface{})
	}{
		{
			name:           "Error - Missing id parameter",
			path:           "/api/qtip/",
			expectedStatus: http.StatusBadRequest,
			validateResult: func(t *testing.T, result map[string]interface{}) {
				assert.False(t, result["success"].(bool))
				assert.Contains(t, result["message"].(string), "id")
			},
		},
		{
			name:           "Valid - Anime ID",
			path:           "/api/qtip/frieren-beyond-journeys-end-18542",
			expectedStatus: http.StatusOK,
			validateResult: func(t *testing.T, result map[string]interface{}) {
				// May succeed or fail depending on external service
				if result["success"] != nil {
					assert.IsType(t, true, result["success"])
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			defer func() { _ = resp.Body.Close() }()

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			var result map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&result)
			if err != nil {
				if resp.StatusCode >= 400 {
					return
				}
				require.NoError(t, err)
			}

			if tc.validateResult != nil && err == nil {
				tc.validateResult(t, result)
			}
		})
	}
}
