package extractors

import (
	"github.com/mnuddindev/jutsu-api/pkg/scrape"
)

type CategoryResult struct {
	Data       []scrape.ExtractedItem `json:"data"`
	TotalPages int                    `json:"totalPages"`
}

func ExtractCategory(path string, page int, baseURL string) (CategoryResult, error) {
	data, total, err := scrape.ExtractPage(page, path, baseURL)
	if err != nil {
		return CategoryResult{}, err
	}
	return CategoryResult{Data: data, TotalPages: total}, nil
}
