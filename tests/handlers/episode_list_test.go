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

func TestEpisodeListHandler_GetEpisodes(t *testing.T) {
	app := fiber.New()
	h := handler.NewEpisodeListHandler()
	app.Get("/api/episodes/:id", h.GetEpisodes)

	testCases := []struct {
		name           string
		path           string
		expectedStatus int
		validateResult func(t *testing.T, result map[string]interface{})
	}{
		{
			name:           "Error - Missing id parameter",
			path:           "/api/episodes/",
			expectedStatus: http.StatusNotFound, // 404 when path param is missing
			validateResult: func(t *testing.T, result map[string]interface{}) {
			},
		},
		{
			name:           "Success - Valid anime ID",
			path:           "/api/episodes/frieren-beyond-journeys-end-18542",
			expectedStatus: http.StatusOK,
			validateResult: func(t *testing.T, result map[string]interface{}) {
				assert.True(t, result["success"].(bool))
				assert.NotNil(t, result["results"])
			},
		},
		{
			name:           "Success - Another valid anime ID",
			path:           "/api/episodes/demon-slayer-kimetsu-no-yaiba-hashira-training-arc-19107",
			expectedStatus: http.StatusOK,
			validateResult: func(t *testing.T, result map[string]interface{}) {
				assert.True(t, result["success"].(bool))
			},
		},
		{
			name:           "Success - Anime ID with special characters",
			path:           "/api/episodes/one-piece-100",
			expectedStatus: http.StatusOK, // May be 502 if upstream fails
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

			// Allow 502 for external upstream failures when we expect 200
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
