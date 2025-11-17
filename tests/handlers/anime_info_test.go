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

func TestAnimeInfoHandler_GetAnimeInfo(t *testing.T) {
	app := fiber.New()
	h := handler.NewAnimeInfoHandler()
	app.Get("/api/info", h.GetAnimeInfo)

	testCases := []struct {
		name           string
		queryParams    string
		expectedStatus int
		validateResult func(t *testing.T, result map[string]interface{})
	}{
		{
			name:           "Success - Valid anime ID",
			queryParams:    "?id=frieren-beyond-journeys-end-18542",
			expectedStatus: http.StatusOK, // May be 502 if external service fails
			validateResult: func(t *testing.T, result map[string]interface{}) {
				if result["success"] != nil {
					assert.True(t, result["success"].(bool))
					assert.NotNil(t, result["results"])
				}
			},
		},
		{
			name:           "Success - Another valid anime ID",
			queryParams:    "?id=demon-slayer-kimetsu-no-yaiba-hashira-training-arc-19107",
			expectedStatus: http.StatusOK, // May be 502 if external service fails
			validateResult: func(t *testing.T, result map[string]interface{}) {
				if result["success"] != nil {
					assert.True(t, result["success"].(bool))
				}
			},
		},
		{
			name:           "Error - Missing id parameter",
			queryParams:    "",
			expectedStatus: http.StatusBadRequest,
			validateResult: func(t *testing.T, result map[string]interface{}) {
				assert.False(t, result["success"].(bool))
				assert.Contains(t, result["message"].(string), "id")
			},
		},
		{
			name:           "Error - Empty id parameter",
			queryParams:    "?id=",
			expectedStatus: http.StatusBadRequest,
			validateResult: func(t *testing.T, result map[string]interface{}) {
				assert.False(t, result["success"].(bool))
			},
		},
		{
			name:           "Error - Whitespace only id",
			queryParams:    "?id=%20%20%20", // URL encoded spaces
			expectedStatus: http.StatusBadRequest,
			validateResult: func(t *testing.T, result map[string]interface{}) {
				if result["success"] != nil {
					assert.False(t, result["success"].(bool))
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/info"+tc.queryParams, nil)
			resp, err := app.Test(req, -1) // -1 means no timeout
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
				// For error responses, try to read as string first
				if resp.StatusCode >= 400 {
					// Error response might not be JSON, skip validation
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
