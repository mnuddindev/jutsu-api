package helper

import (
	"github.com/mnuddindev/jutsu-api/pkg/scrape"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

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
