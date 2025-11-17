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

func TestCharacterHandler_GetCharacter(t *testing.T) {
	app := fiber.New()
	h := handler.NewCharacterHandler()
	app.Get("/api/character/:id", h.GetCharacter)

	testCases := []struct {
		name           string
		path           string
		expectedStatus int
		validateResult func(t *testing.T, result map[string]interface{})
	}{
		{
			name:           "Error - Missing id parameter",
			path:           "/api/character/",
			expectedStatus: http.StatusBadRequest,
			validateResult: func(t *testing.T, result map[string]interface{}) {
				assert.False(t, result["success"].(bool))
				assert.Contains(t, result["message"].(string), "id")
			},
		},
		{
			name:           "Valid - Character ID",
			path:           "/api/character/asta-340",
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

func TestCharacterListHandler_GetVoiceActors(t *testing.T) {
	app := fiber.New()
	h := handler.NewCharacterListHandler()
	app.Get("/api/character/list/:id", h.GetVoiceActors)

	testCases := []struct {
		name           string
		path           string
		queryParams    string
		expectedStatus int
		validateResult func(t *testing.T, result map[string]interface{})
	}{
		{
			name:           "Error - Missing id parameter",
			path:           "/api/character/list/",
			queryParams:    "",
			expectedStatus: http.StatusBadRequest,
			validateResult: func(t *testing.T, result map[string]interface{}) {
				assert.False(t, result["success"].(bool))
			},
		},
		{
			name:           "Success - Valid anime ID with page 1",
			path:           "/api/character/list/frieren-beyond-journeys-end-18542",
			queryParams:    "?page=1",
			expectedStatus: http.StatusOK,
			validateResult: func(t *testing.T, result map[string]interface{}) {
				assert.True(t, result["success"].(bool))
			},
		},
		{
			name:           "Success - Default page when not specified",
			path:           "/api/character/list/frieren-beyond-journeys-end-18542",
			queryParams:    "",
			expectedStatus: http.StatusOK,
			validateResult: func(t *testing.T, result map[string]interface{}) {
				assert.True(t, result["success"].(bool))
			},
		},
		{
			name:           "Success - Page 2",
			path:           "/api/character/list/frieren-beyond-journeys-end-18542",
			queryParams:    "?page=2",
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

func TestActorsHandler_GetVoiceActor(t *testing.T) {
	app := fiber.New()
	h := handler.NewActorsHandler()
	app.Get("/api/actors/:id", h.GetVoiceActor)

	testCases := []struct {
		name           string
		path           string
		expectedStatus int
		validateResult func(t *testing.T, result map[string]interface{})
	}{
		{
			name:           "Error - Missing id parameter",
			path:           "/api/actors/",
			expectedStatus: http.StatusBadRequest,
			validateResult: func(t *testing.T, result map[string]interface{}) {
				assert.False(t, result["success"].(bool))
			},
		},
		{
			name:           "Valid - Voice actor ID",
			path:           "/api/actors/gakuto-kajiwara-534",
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
