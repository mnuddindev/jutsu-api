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

func TestSuggestionHandler_GetSuggestions(t *testing.T) {
	app := fiber.New()
	h := handler.NewSuggestionHandler()
	app.Get("/api/search/suggest", h.GetSuggestions)

	testCases := []struct {
		name           string
		queryParams    string
		expectedStatus int
		validateResult func(t *testing.T, result map[string]interface{})
	}{
		{
			name:           "Error - Missing keyword parameter",
			queryParams:    "",
			expectedStatus: http.StatusBadRequest,
			validateResult: func(t *testing.T, result map[string]interface{}) {
				assert.False(t, result["success"].(bool))
				assert.Contains(t, result["message"].(string), "keyword")
			},
		},
		{
			name:           "Error - Empty keyword parameter",
			queryParams:    "?keyword=",
			expectedStatus: http.StatusBadRequest,
			validateResult: func(t *testing.T, result map[string]interface{}) {
				assert.False(t, result["success"].(bool))
			},
		},
		{
			name:           "Success - Valid keyword",
			queryParams:    "?keyword=naru",
			expectedStatus: http.StatusOK, // May be 502 if external service fails
			validateResult: func(t *testing.T, result map[string]interface{}) {
				if result != nil && result["success"] != nil {
					assert.True(t, result["success"].(bool))
					if result["results"] != nil {
						assert.NotNil(t, result["results"])
					}
				}
			},
		},
		{
			name:           "Success - Full keyword",
			queryParams:    "?keyword=naruto",
			expectedStatus: http.StatusOK,
			validateResult: func(t *testing.T, result map[string]interface{}) {
				assert.True(t, result["success"].(bool))
			},
		},
		{
			name:           "Success - Keyword with special characters",
			queryParams:    "?keyword=one-piece",
			expectedStatus: http.StatusOK,
			validateResult: func(t *testing.T, result map[string]interface{}) {
				assert.True(t, result["success"].(bool))
			},
		},
		{
			name:           "Success - Short keyword",
			queryParams:    "?keyword=op",
			expectedStatus: http.StatusOK,
			validateResult: func(t *testing.T, result map[string]interface{}) {
				assert.True(t, result["success"].(bool))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/search/suggest"+tc.queryParams, nil)
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
				// Only validate if we got a successful response or if result is not nil
				if resp.StatusCode == http.StatusOK || result != nil {
					tc.validateResult(t, result)
				}
			}
		})
	}
}
