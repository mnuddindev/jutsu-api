package extractors

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

type FilterResult struct {
	TotalPage int                      `json:"totalPage"`
	Data      []map[string]interface{} `json:"data"`
	Page      int                      `json:"page"`
	HasNext   bool                     `json:"hasNext"`
}

func ExtractFilter(params map[string]string, baseURL string) (FilterResult, error) {
	q := url.Values{}
	for k, v := range params {
		if strings.TrimSpace(v) != "" {
			q.Set(k, v)
		}
	}
	if q.Get("page") == "" {
		q.Set("page", "1")
	}
	apiURL := fmt.Sprintf("https://%s/filter?%s", baseURL, q.Encode())
	if q.Get("keyword") != "" {
		apiURL = fmt.Sprintf("https://%s/search?%s", baseURL, q.Encode())
	}
	c := utils.NewCollector()
	res := FilterResult{Page: atoiOr(q.Get("page"), 1)}
	c.OnHTML(".flw-item", func(e *colly.HTMLElement) {
		m := map[string]interface{}{}
		href := e.ChildAttr(".film-poster-ahref", "href")
		dataID := e.ChildAttr(".film-poster-ahref", "data-id")
		m["id"] = trimPrefixOrNil(href)
		m["data_id"] = dataID
		m["poster"] = firstNonEmpty(e.ChildAttr(".film-poster .film-poster-img", "data-src"), e.ChildAttr(".film-poster .film-poster-img", "src"))
		m["title"] = strings.TrimSpace(e.ChildText(".film-name .dynamic-name"))
		m["japanese_title"] = e.DOM.Find(".film-name .dynamic-name").AttrOr("data-jname", "")
		showType := strings.TrimSpace(e.ChildText(".fd-infor .fdi-item:first-child"))
		if showType == "" {
			showType = "Unknown"
		}
		duration := strings.TrimSpace(e.ChildText(".fd-infor .fdi-duration"))
		sub := digits(e.ChildText(".tick-sub"))
		dub := digits(e.ChildText(".tick-dub"))
		eps := digits(e.ChildText(".tick-eps"))
		m["tvInfo"] = map[string]interface{}{"showType": showType, "duration": duration, "sub": sub, "dub": dub, "eps": eps}
		m["adultContent"] = strings.TrimSpace(e.ChildText(".tick-rate"))
		res.Data = append(res.Data, m)
	})
	c.OnHTML(`.pre-pagination nav .pagination > .page-item a[title="Last"]`, func(e *colly.HTMLElement) {
		res.TotalPage = atoiOr(strings.Split(lastSegment(e.Attr("href")), "page=")[1], 1)
	})
	c.OnHTML(`.pre-pagination nav .pagination > .page-item a[title="Next"]`, func(e *colly.HTMLElement) {
		if res.TotalPage == 0 {
			res.TotalPage = atoiOr(strings.Split(lastSegment(e.Attr("href")), "page=")[1], 1)
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
		return FilterResult{}, errVisit
	}
	if res.TotalPage == 0 {
		res.TotalPage = 1
	}
	res.HasNext = res.Page < res.TotalPage
	return res, nil
}

func trimPrefixOrNil(href string) string {
	if href == "" {
		return ""
	}
	return strings.TrimPrefix(href, "/")
}
func atoiOr(s string, d int) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		return d
	}
	return n
}
func digits(s string) int {
	s = strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, s)
	n, _ := strconv.Atoi(s)
	return n
}
