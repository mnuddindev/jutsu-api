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

func TestFilterHandler_Filter(t *testing.T) {
	app := fiber.New()
	h := handler.NewFilterHandler()
	app.Get("/api/filter", h.Filter)

	testCases := []struct {
		name           string
		queryParams    string
		expectedStatus int
		validateResult func(t *testing.T, result map[string]interface{})
	}{
		{
			name:           "Success - Filter by type only",
			queryParams:    "?type=tv",
			expectedStatus: http.StatusOK,
			validateResult: func(t *testing.T, result map[string]interface{}) {
				assert.True(t, result["success"].(bool))
				assert.NotNil(t, result["results"])
			},
		},
		{
			name:           "Success - Filter by type and status",
			queryParams:    "?type=tv&status=ongoing",
			expectedStatus: http.StatusOK,
			validateResult: func(t *testing.T, result map[string]interface{}) {
				assert.True(t, result["success"].(bool))
			},
		},
		{
			name:           "Success - Filter with all parameters",
			queryParams:    "?type=tv&status=ongoing&rated=pg-13&score=8&season=winter&language=sub&genres=1,2,3&sort=popularity&page=1",
			expectedStatus: http.StatusOK,
			validateResult: func(t *testing.T, result map[string]interface{}) {
				assert.True(t, result["success"].(bool))
			},
		},
		{
			name:           "Success - Filter with keyword",
			queryParams:    "?keyword=naruto&type=tv&page=1",
			expectedStatus: http.StatusOK,
			validateResult: func(t *testing.T, result map[string]interface{}) {
				assert.True(t, result["success"].(bool))
			},
		},
		{
			name:           "Success - Empty filter (all anime)",
			queryParams:    "?page=1",
			expectedStatus: http.StatusOK,
			validateResult: func(t *testing.T, result map[string]interface{}) {
				if result["success"] != nil {
					assert.True(t, result["success"].(bool))
				}
			},
		},
		{
			name:           "Success - Default page when not specified",
			queryParams:    "?type=tv",
			expectedStatus: http.StatusOK,
			validateResult: func(t *testing.T, result map[string]interface{}) {
				if result["success"] != nil {
					assert.True(t, result["success"].(bool))
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/filter"+tc.queryParams, nil)
			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			defer func() { _ = resp.Body.Close() }()

			// Allow 502 for external service failures when we expect 200
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
