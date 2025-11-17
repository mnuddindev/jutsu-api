package scrape

import (
	"github.com/gocolly/colly/v2"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

func FetchScript(url string) (string, error) {
	c := utils.NewCollector()
	var body string
	var visitErr error
	c.OnResponse(func(r *colly.Response) {
		body = string(r.Body)
	})
	c.OnError(func(_ *colly.Response, err error) { visitErr = err })
	_ = c.Visit(url)
	c.Wait()
	if visitErr != nil {
		return "", visitErr
	}
	return body, nil
}
