package scrape

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

func CountPages(url string) (int, error) {
	c := utils.NewCollector()
	pages := 1
	var lastHref string
	var haveLast bool

	c.OnHTML(".tab-content .pagination .page-item:last-child a", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		if strings.TrimSpace(href) != "" {
			lastHref = href
			haveLast = true
		}
	})
	var reqErr error
	c.OnError(func(_ *colly.Response, err error) {
		reqErr = err
	})
	_ = c.Visit(url)
	c.Wait()
	if reqErr != nil {
		return 0, reqErr
	}
	if !haveLast || strings.TrimSpace(lastHref) == "" {
		return pages, nil
	}
	parts := strings.Split(lastHref, "=")
	last := parts[len(parts)-1]
	n, err := strconv.Atoi(last)
	if err != nil {
		return 0, errors.New("invalid last page number")
	}
	return n, nil
}
