package extractors

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

type SearchResult struct {
	TotalPage int                      `json:"totalPage"`
	Data      []map[string]interface{} `json:"data"`
}

func ExtractSearch(params map[string]string, baseURL string) (SearchResult, error) {
	q := url.Values{}
	for k, v := range params {
		if strings.TrimSpace(v) != "" {
			q.Set(k, v)
		}
	}
	if q.Get("page") == "" {
		q.Set("page", "1")
	}
	apiURL := fmt.Sprintf("https://%s/search?%s", baseURL, q.Encode())

	c := utils.NewCollector()
	res := SearchResult{}
	c.OnHTML("#main-content .film_list-wrap .flw-item", func(e *colly.HTMLElement) {
		m := map[string]interface{}{}
		id := e.ChildAttr(".film-detail .film-name .dynamic-name", "href")
		id = strings.TrimPrefix(id, "/")
		id = strings.Split(id, "?ref=search")[0]
		m["id"] = id
		m["title"] = strings.TrimSpace(e.ChildText(".film-detail .film-name .dynamic-name"))
		m["japanese_title"] = e.DOM.Find(".film-detail .film-name .dynamic-name").AttrOr("data-jname", "")
		m["poster"] = strings.TrimSpace(e.ChildAttr(".film-poster .film-poster-img", "data-src"))
		m["duration"] = strings.TrimSpace(e.ChildText(".film-detail .fd-infor .fdi-item.fdi-duration"))
		showType := strings.TrimSpace(e.ChildText(".film-detail .fd-infor .fdi-item:nth-of-type(1)"))
		if showType == "" {
			showType = "Unknown"
		}
		rating := strings.TrimSpace(e.ChildText(".film-poster .tick-rate"))
		sub := lastNumber(e.ChildText(".film-poster .tick-sub"))
		dub := lastNumber(e.ChildText(".film-poster .tick-dub"))
		eps := lastNumber(e.ChildText(".film-poster .tick-eps"))
		m["tvInfo"] = map[string]interface{}{"showType": showType, "rating": rating, "sub": sub, "dub": dub, "eps": eps}
		res.Data = append(res.Data, m)
	})
	c.OnHTML(`.pre-pagination nav .pagination > .page-item a[title="Last"]`, func(e *colly.HTMLElement) {
		res.TotalPage = atoiOr(lastSegment(e.Attr("href")), 1)
	})
	c.OnHTML(`.pre-pagination nav .pagination > .page-item a[title="Next"]`, func(e *colly.HTMLElement) {
		if res.TotalPage == 0 {
			res.TotalPage = atoiOr(lastSegment(e.Attr("href")), 1)
		}
	})
	c.OnHTML(`.pre-pagination nav .pagination > .page-item.active a`, func(e *colly.HTMLElement) {
		if res.TotalPage == 0 {
			res.TotalPage = atoiOr(strings.TrimSpace(e.Text), 1)
		}
	})
	var errVisit error
	c.OnError(func(_ *colly.Response, err error) { errVisit = err })
	_ = c.Visit(apiURL)
	c.Wait()
	if errVisit != nil {
		return SearchResult{}, errVisit
	}
	if res.TotalPage == 0 {
		res.TotalPage = 1
	}
	return res, nil
}

func lastNumber(s string) int {
	parts := strings.Fields(s)
	if len(parts) == 0 {
		return 0
	}
	n, _ := strconv.Atoi(parts[len(parts)-1])
	return n
}
