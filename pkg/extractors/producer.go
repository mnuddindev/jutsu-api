package extractors

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mnuddindev/jutsu-api/pkg/scrape"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

var extractPageFn = scrape.ExtractPage

// ErrCreatorPageOutOfRange is returned when the requested page is larger than the
// number of available pages for a producer/studio listing.
var ErrCreatorPageOutOfRange = errors.New("creator page exceeds available pages")

// ErrCreatorIDRequired is returned when no producer/studio identifier is provided.
var ErrCreatorIDRequired = errors.New("creator id is required")

// ExtractProducer extracts producer information, returning the paginated
// list of anime for a given producer slug.
func ExtractProducer(id string, page int, baseURL string) (CategoryResult, error) {
	return extractCreatorList("producer", id, page, baseURL)
}

// ExtractStudio extracts studio information (same layout as producer pages).
func ExtractStudio(id string, page int, baseURL string) (CategoryResult, error) {
	return extractCreatorList("studio", id, page, baseURL)
}

func extractCreatorList(prefix, id string, page int, baseURL string) (CategoryResult, error) {
	slug := strings.Trim(strings.TrimSpace(id), "/")
	if slug == "" {
		return CategoryResult{}, ErrCreatorIDRequired
	}
	if page <= 0 {
		page = 1
	}
	route := fmt.Sprintf("%s/%s", prefix, utils.ExtractDataID(slug))
	data, totalPages, err := extractPageFn(page, route, baseURL)
	if err != nil {
		return CategoryResult{}, err
	}
	if totalPages > 0 && page > totalPages {
		return CategoryResult{}, ErrCreatorPageOutOfRange
	}
	return CategoryResult{
		Data:       data,
		TotalPages: totalPages,
	}, nil
}

// SetExtractPageFuncForTest allows tests to replace the scraping function and returns
// a restore callback to reinstate the default behavior.
func SetExtractPageFuncForTest(fn func(int, string, string) ([]scrape.ExtractedItem, int, error)) func() {
	prev := extractPageFn
	if fn != nil {
		extractPageFn = fn
	}
	return func() {
		extractPageFn = prev
	}
}
