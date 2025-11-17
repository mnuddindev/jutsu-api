package extractors

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

type SeasonItem struct {
	ID           string `json:"id"`
	DataNumber   int    `json:"data_number"`
	DataID       int    `json:"data_id"`
	Season       string `json:"season"`
	Title        string `json:"title"`
	SeasonPoster string `json:"season_poster"`
}

func ExtractSeasons(id, baseURL string) ([]SeasonItem, error) {
	c := utils.NewCollector()
	var items []SeasonItem
	url := fmt.Sprintf("https://%s/watch/%s", baseURL, id)
	c.OnHTML(".anis-watch>.other-season>.inner>.os-list>a", func(e *colly.HTMLElement) {
		var it SeasonItem
		it.DataNumber = e.Index
		href := e.Attr("href")
		if href != "" {
			parts := strings.Split(href, "-")
			if len(parts) > 0 {
				last := parts[len(parts)-1]
				if n, err := strconvAtoi(last); err == nil {
					it.DataID = n
				}
			}
			it.ID = strings.TrimPrefix(href, "/")
		}
		it.Season = strings.TrimSpace(e.ChildText(".title"))
		it.Title = strings.TrimSpace(e.Attr("title"))
		style := e.ChildAttr(".season-poster", "style")
		re := regexp.MustCompile(`url\((.*?)\)`)
		m := re.FindStringSubmatch(style)
		if len(m) == 2 {
			it.SeasonPoster = m[1]
		}
		items = append(items, it)
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

func strconvAtoi(s string) (int, error) {
	n := 0
	for _, r := range s {
		if r < '0' || r > '9' {
			return 0, fmt.Errorf("nan")
		}
		n = n*10 + int(r-'0')
	}
	return n, nil
}
