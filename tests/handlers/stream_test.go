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

func TestStreamHandler_GetStream(t *testing.T) {
	app := fiber.New()
	h := handler.NewStreamHandler()
	app.Get("/api/stream", h.GetStream)

	testCases := []struct {
		name           string
		queryParams    string
		expectedStatus int
		validateResult func(t *testing.T, result map[string]interface{})
	}{
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
			name:           "Valid - With id, server, and type",
			queryParams:    "?id=frieren-beyond-journeys-end-18542?ep=107257&server=hd-1&type=sub",
			expectedStatus: http.StatusOK,
			validateResult: func(t *testing.T, result map[string]interface{}) {
				// May succeed or fail depending on external service
				if result["success"] != nil {
					assert.IsType(t, true, result["success"])
				}
			},
		},
		{
			name:           "Valid - With id only",
			queryParams:    "?id=frieren-beyond-journeys-end-18542?ep=107257",
			expectedStatus: http.StatusOK,
			validateResult: func(t *testing.T, result map[string]interface{}) {
				// May succeed or fail depending on external service
				if result["success"] != nil {
					assert.IsType(t, true, result["success"])
				}
			},
		},
		{
			name:           "Valid - With id and server",
			queryParams:    "?id=frieren-beyond-journeys-end-18542?ep=107257&server=vidcloud",
			expectedStatus: http.StatusOK,
			validateResult: func(t *testing.T, result map[string]interface{}) {
				// May succeed or fail depending on external service
				if result["success"] != nil {
					assert.IsType(t, true, result["success"])
				}
			},
		},
		{
			name:           "Valid - With id, server, and type=dub",
			queryParams:    "?id=frieren-beyond-journeys-end-18542?ep=107257&server=hd-1&type=dub",
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
			req := httptest.NewRequest(http.MethodGet, "/api/stream"+tc.queryParams, nil)
			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			defer func() { _ = resp.Body.Close() }()

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			var result map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&result)
			if err != nil {
				if resp.StatusCode >= 400 {
					return // Error response might not be JSON
				}
				require.NoError(t, err)
			}

			if tc.validateResult != nil && err == nil {
				tc.validateResult(t, result)
			}
		})
	}
}

func TestStreamHandler_GetStreamFallback(t *testing.T) {
	app := fiber.New()
	h := handler.NewStreamHandler()
	app.Get("/api/stream/fallback", h.GetStreamFallback)

	testCases := []struct {
		name           string
		queryParams    string
		expectedStatus int
	}{
		{
			name:           "Error - Missing id parameter",
			queryParams:    "",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Valid - With all parameters",
			queryParams:    "?id=frieren-beyond-journeys-end-18542?ep=107257&server=hd-1&type=sub",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/stream/fallback"+tc.queryParams, nil)
			resp, err := app.Test(req, -1)
			require.NoError(t, err)
			defer func() { _ = resp.Body.Close() }()

			assert.Equal(t, tc.expectedStatus, resp.StatusCode)
		})
	}
}
