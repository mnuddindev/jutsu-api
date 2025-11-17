package extractors

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

type TopTenItem struct {
	ID            string            `json:"id"`
	DataID        string            `json:"data_id"`
	Number        string            `json:"number"`
	Title         string            `json:"title"`
	JapaneseTitle string            `json:"japanese_title"`
	Poster        string            `json:"poster"`
	TVInfo        map[string]string `json:"tvInfo"`
}

type TopTenResult map[string][]TopTenItem

func ExtractTopTen(baseURL string) (TopTenResult, error) {
	c := utils.NewCollector()
	res := TopTenResult{"today": {}, "week": {}, "month": {}}
	url := fmt.Sprintf("https://%s/home", baseURL)
	labels := []string{"today", "week", "month"}
	for idx, label := range labels {
		selector := fmt.Sprintf(`#main-sidebar .block_area-realtime .block_area-content ul:eq(%d)>li`, idx)
		l := label
		c.OnHTML(selector, func(e *colly.HTMLElement) {
			item := TopTenItem{}
			item.Number = strings.TrimSpace(e.ChildText(".film-number>span"))
			item.Title = strings.TrimSpace(e.ChildText(".film-detail>.film-name>a"))
			item.Poster = e.ChildAttr(".film-poster>img", "data-src")
			item.JapaneseTitle = strings.TrimSpace(e.DOM.Find(".film-detail>.film-name>a").AttrOr("data-jname", ""))
			item.DataID = e.ChildAttr(".film-poster", "data-id")
			item.ID = lastSegment(e.ChildAttr(".film-detail>.film-name>a", "href"))
			item.TVInfo = map[string]string{}
			for _, prop := range []string{"sub", "dub", "eps"} {
				v := strings.TrimSpace(e.DOM.Find(".tick .tick-" + prop).Text())
				if v != "" {
					item.TVInfo[prop] = v
				}
			}
			res[l] = append(res[l], item)
		})
	}
	var errVisit error
	c.OnError(func(_ *colly.Response, err error) { errVisit = err })
	_ = c.Visit(url)
	c.Wait()
	if errVisit != nil {
		return nil, errVisit
	}
	return res, nil
}
