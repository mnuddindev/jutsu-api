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

func TestCategoryHandler_GetCategory(t *testing.T) {
	app := fiber.New()
	h := handler.NewCategoryHandler()

	testCases := []struct {
		name           string
		path           string
		queryParams    string
		expectedStatus int
		validateResult func(t *testing.T, result map[string]interface{})
	}{
		{
			name:           "Success - Genre action with page 1",
			path:           "/api/genre/action",
			queryParams:    "?page=1",
			expectedStatus: http.StatusOK,
			validateResult: func(t *testing.T, result map[string]interface{}) {
				assert.True(t, result["success"].(bool))
				assert.NotNil(t, result["results"])
			},
		},
		{
			name:           "Success - Genre comedy",
			path:           "/api/genre/comedy",
			queryParams:    "?page=1",
			expectedStatus: http.StatusOK,
			validateResult: func(t *testing.T, result map[string]interface{}) {
				assert.True(t, result["success"].(bool))
			},
		},
		{
			name:           "Success - Top airing",
			path:           "/api/top-airing",
			queryParams:    "?page=1",
			expectedStatus: http.StatusOK,
			validateResult: func(t *testing.T, result map[string]interface{}) {
				assert.True(t, result["success"].(bool))
			},
		},
		{
			name:           "Success - Most popular",
			path:           "/api/most-popular",
			queryParams:    "?page=1",
			expectedStatus: http.StatusOK,
			validateResult: func(t *testing.T, result map[string]interface{}) {
				assert.True(t, result["success"].(bool))
			},
		},
		{
			name:           "Success - Default page when not specified",
			path:           "/api/genre/action",
			queryParams:    "",
			expectedStatus: http.StatusOK,
			validateResult: func(t *testing.T, result map[string]interface{}) {
				assert.True(t, result["success"].(bool))
			},
		},
		{
			name:           "Error - Invalid page number (negative)",
			path:           "/api/genre/action",
			queryParams:    "?page=-1",
			expectedStatus: http.StatusBadRequest,
			validateResult: func(t *testing.T, result map[string]interface{}) {
				assert.False(t, result["success"].(bool))
			},
		},
		{
			name:           "Error - Invalid page number (zero)",
			path:           "/api/genre/action",
			queryParams:    "?page=0",
			expectedStatus: http.StatusBadRequest,
			validateResult: func(t *testing.T, result map[string]interface{}) {
				assert.False(t, result["success"].(bool))
			},
		},
		{
			name:           "Error - Invalid page number (non-numeric)",
			path:           "/api/genre/action",
			queryParams:    "?page=abc",
			expectedStatus: http.StatusBadRequest,
			validateResult: func(t *testing.T, result map[string]interface{}) {
				assert.False(t, result["success"].(bool))
			},
		},
		{
			name:           "Success - Page 2",
			path:           "/api/genre/action",
			queryParams:    "?page=2",
			expectedStatus: http.StatusOK,
			validateResult: func(t *testing.T, result map[string]interface{}) {
				assert.True(t, result["success"].(bool))
			},
		},
		{
			name:           "Success - Martial arts (with typo fix)",
			path:           "/api/genre/martial-arts",
			queryParams:    "?page=1",
			expectedStatus: http.StatusOK,
			validateResult: func(t *testing.T, result map[string]interface{}) {
				assert.True(t, result["success"].(bool))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			app.Get(tc.path, h.GetCategory)
			req := httptest.NewRequest(http.MethodGet, tc.path+tc.queryParams, nil)
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
