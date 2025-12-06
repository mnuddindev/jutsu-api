package helper

import (
	"context"
	"encoding/json"

	"github.com/mnuddindev/jutsu-api/internal/infrastructure/cache"
	"github.com/mnuddindev/jutsu-api/pkg/scrape"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

// GetCachedData retrieves a cached value by key into dest using the cache manager.
// If the key does not exist, dest is left untouched and an error is returned.
// Requires a background context for the operation.
func GetCachedData(ctx context.Context, manager *cache.Manager, category cache.CacheCategory, key string, dest interface{}) error {
	if !manager.IsEnabled() {
		return nil
	}

	data, err := manager.Get(ctx, category, key)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, dest)
}

// SetCachedData stores value at key with the category's TTL using the cache manager.
// Requires a background context for the operation.
func SetCachedData(ctx context.Context, manager *cache.Manager, category cache.CacheCategory, key string, value interface{}) error {
	if !manager.IsEnabled() {
		return nil
	}

	return manager.Set(ctx, category, key, value)
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
