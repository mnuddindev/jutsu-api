package helper

import (
	"time"

	"github.com/mnuddindev/jutsu-api/internal/infrastructure/cache"
	"github.com/mnuddindev/jutsu-api/pkg/scrape"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

// Cache helpers (mirror cache.helper.js using the app Redis cache)

// GetCachedData retrieves a cached value by key into dest.
// If the key does not exist, dest is left untouched and nil is returned.
func GetCachedData(key string, dest interface{}) error {
	if cache.Client == nil {
		return nil
	}
	return cache.Get(key, dest)
}

// SetCachedData stores value at key with the provided ttl.
// If ttl is zero, a default of 24h is used.
func SetCachedData(key string, value interface{}, ttl time.Duration) error {
	if cache.Client == nil {
		return nil
	}
	if ttl <= 0 {
		ttl = 24 * time.Hour
	}
	return cache.Set(key, value, ttl)
}

// Scraping helpers – thin wrappers around existing Go implementations.

// CountPages mirrors countPages.helper.js.
func CountPages(url string) (int, error) {
	return scrape.CountPages(url)
}

// ExtractPage mirrors extractPages.helper.js.
func ExtractPage(page int, params string, baseURL string) (interface{}, int, error) {
	return scrape.ExtractPage(page, params, baseURL)
}

// FetchScript mirrors fetchScript.helper.js.
func FetchScript(url string) (string, error) {
	return scrape.FetchScript(url)
}

// FormatTitle mirrors formatTitle.helper.js.
func FormatTitle(title, dataID string) string {
	return utils.FormatTitle(title, dataID)
}

// GetKeys mirrors getKey.helper.js.
func GetKeys(script string) [][]int {
	return utils.GetKeys(script)
}

// ExtractToken mirrors token.helper.js.
func ExtractToken(url string, baseURL string) (string, error) {
	return scrape.ExtractToken(url, baseURL)
}
