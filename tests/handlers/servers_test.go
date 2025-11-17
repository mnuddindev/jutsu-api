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

func TestServersHandler_GetServers(t *testing.T) {
	app := fiber.New()
	h := handler.NewServersHandler()
	app.Get("/api/servers/:id", h.GetServers)

	testCases := []struct {
		name           string
		path           string
		queryParams    string
		expectedStatus int
		validateResult func(t *testing.T, result map[string]interface{})
	}{
		{
			name:           "Error - Missing ep parameter",
			path:           "/api/servers/demon-slayer-kimetsu-no-yaiba-hashira-training-arc-19107",
			queryParams:    "",
			expectedStatus: http.StatusBadRequest,
			validateResult: func(t *testing.T, result map[string]interface{}) {
				assert.False(t, result["success"].(bool))
				assert.Contains(t, result["message"].(string), "ep")
			},
		},
		{
			name:           "Error - Empty ep parameter",
			path:           "/api/servers/demon-slayer-kimetsu-no-yaiba-hashira-training-arc-19107",
			queryParams:    "?ep=",
			expectedStatus: http.StatusBadRequest,
			validateResult: func(t *testing.T, result map[string]interface{}) {
				assert.False(t, result["success"].(bool))
			},
		},
		{
			name:           "Success - Valid episode ID",
			path:           "/api/servers/demon-slayer-kimetsu-no-yaiba-hashira-training-arc-19107",
			queryParams:    "?ep=124260",
			expectedStatus: http.StatusOK,
			validateResult: func(t *testing.T, result map[string]interface{}) {
				assert.True(t, result["success"].(bool))
				assert.NotNil(t, result["results"])
			},
		},
		{
			name:           "Success - Another valid episode ID",
			path:           "/api/servers/frieren-beyond-journeys-end-18542",
			queryParams:    "?ep=107257",
			expectedStatus: http.StatusOK,
			validateResult: func(t *testing.T, result map[string]interface{}) {
				assert.True(t, result["success"].(bool))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
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
