package extractors

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

type TopSearchItem struct {
	Title string `json:"title"`
	Link  string `json:"link"`
}

func ExtractTopSearch(baseURL string) ([]TopSearchItem, error) {
	c := utils.NewCollector()
	var items []TopSearchItem
	url := fmt.Sprintf("https://%s", baseURL)
	c.OnHTML(".xhashtag a.item", func(e *colly.HTMLElement) {
		items = append(items, TopSearchItem{Title: strings.TrimSpace(e.Text), Link: e.Attr("href")})
	})
	var errVisit error
	c.OnError(func(_ *colly.Response, err error) { errVisit = err })
	_ = c.Visit(url)
	c.Wait()
	if errVisit != nil {
		return nil, errVisit
	}
	return items, nil
}
