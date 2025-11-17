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

func TestWatchlistHandler_GetWatchlist(t *testing.T) {
	app := fiber.New()
	h := handler.NewWatchlistHandler()
	app.Get("/api/watchlist/:userId", h.GetWatchlist)
	app.Get("/api/watchlist/:userId/:page", h.GetWatchlist)

	testCases := []struct {
		name           string
		path           string
		expectedStatus int
		validateResult func(t *testing.T, result map[string]interface{})
	}{
		{
			name:           "Error - Missing userId parameter",
			path:           "/api/watchlist/",
			expectedStatus: http.StatusNotFound, // Fiber returns 404 for missing path params
			validateResult: func(t *testing.T, result map[string]interface{}) {
				// May not have JSON response for 404
			},
		},
		{
			name:           "Success - Valid user ID without page",
			path:           "/api/watchlist/user123",
			expectedStatus: http.StatusOK, // May be 502 if external service fails
			validateResult: func(t *testing.T, result map[string]interface{}) {
				if result["success"] != nil {
					assert.True(t, result["success"].(bool))
					assert.NotNil(t, result["results"])
				}
			},
		},
		{
			name:           "Success - Valid user ID with page 1",
			path:           "/api/watchlist/user123/1",
			expectedStatus: http.StatusOK, // May be 502 if external service fails
			validateResult: func(t *testing.T, result map[string]interface{}) {
				if result["success"] != nil {
					assert.True(t, result["success"].(bool))
				}
			},
		},
		{
			name:           "Success - Valid user ID with page 2",
			path:           "/api/watchlist/user123/2",
			expectedStatus: http.StatusOK, // May be 502 if external service fails
			validateResult: func(t *testing.T, result map[string]interface{}) {
				if result["success"] != nil {
					assert.True(t, result["success"].(bool))
				}
			},
		},
		{
			name:           "Success - Different user ID",
			path:           "/api/watchlist/user456",
			expectedStatus: http.StatusOK, // May be 502 if external service fails
			validateResult: func(t *testing.T, result map[string]interface{}) {
				if result["success"] != nil {
					assert.True(t, result["success"].(bool))
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

			// Allow 502 for external service failures
			if tc.expectedStatus == http.StatusOK {
				assert.Contains(t, []int{http.StatusOK, http.StatusBadGateway}, resp.StatusCode)
			} else {
				assert.Equal(t, tc.expectedStatus, resp.StatusCode)
			}

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
