package helper

import (
	"time"

	"github.com/mnuddindev/jutsu-api/internal/infrastructure/cache"
	"github.com/mnuddindev/jutsu-api/pkg/scrape"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

// Cache helpers using the app Redis cache

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

// CountPages counts the number of pages for a given URL.
func CountPages(url string) (int, error) {
	return scrape.CountPages(url)
}

// ExtractPage extracts page data.
func ExtractPage(page int, params string, baseURL string) (interface{}, int, error) {
	return scrape.ExtractPage(page, params, baseURL)
}

// FetchScript fetches script content from a URL.
func FetchScript(url string) (string, error) {
	return scrape.FetchScript(url)
}

// FormatTitle formats anime title.
func FormatTitle(title, dataID string) string {
	return utils.FormatTitle(title, dataID)
}

// ExtractToken extracts token from URL.
func ExtractToken(url string, baseURL string) (string, error) {
	return scrape.ExtractToken(url, baseURL)
}
