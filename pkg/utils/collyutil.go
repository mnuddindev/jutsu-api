package utils

import (
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/mnuddindev/jutsu-api/pkg/httpclient"
)

// NewCollector returns a preconfigured colly Collector with sane defaults
func NewCollector() *colly.Collector {
	c := colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.Async(true),
		colly.MaxDepth(2),
	)
	c.OnRequest(func(r *colly.Request) {
		for k, v := range httpclient.DEFAULT_HEADERS {
			r.Headers.Set(k, v)
		}
	})
	_ = c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		RandomDelay: 500 * time.Millisecond,
		Parallelism: 4,
	})
	return c
}
