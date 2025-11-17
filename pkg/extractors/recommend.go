package extractors

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

type RecommendItem struct {
	DataID        string            `json:"data_id"`
	ID            string            `json:"id"`
	Title         string            `json:"title"`
	JapaneseTitle string            `json:"japanese_title"`
	Poster        string            `json:"poster"`
	TVInfo        map[string]string `json:"tvInfo"`
	AdultContent  bool              `json:"adultContent"`
}

func ExtractRecommendedData(id, baseURL string) ([]RecommendItem, error) {
	url := fmt.Sprintf("https://%s/%s", baseURL, id)
	c := utils.NewCollector()
	selector := `#main-content .block_area_category .tab-content .block_area-content .film_list-wrap .flw-item`
	var items []RecommendItem
	c.OnHTML(selector, func(e *colly.HTMLElement) {
		var it RecommendItem
		it.ID = lastSegment(e.ChildAttr(".film-detail .film-name a", "href"))
		it.DataID = e.ChildAttr(".film-poster a", "data-id")
		it.Title = strings.TrimSpace(e.ChildText(".film-detail .film-name a"))
		it.JapaneseTitle = strings.TrimSpace(e.DOM.Find(".film-detail .film-name a").AttrOr("data-jname", ""))
		it.Poster = e.ChildAttr(".film-poster img", "data-src")
		// showType and duration
		it.TVInfo = map[string]string{}
		e.DOM.Find(".film-detail .fd-infor .fdi-item").Each(func(_ int, s *goquery.Selection) {
			t := strings.ToLower(strings.TrimSpace(s.Text()))
			for _, typ := range []string{"tv", "ona", "movie", "ova", "special"} {
				if strings.Contains(t, typ) && it.TVInfo["showType"] == "" {
					it.TVInfo["showType"] = strings.TrimSpace(s.Text())
				}
			}
		})
		if it.TVInfo["showType"] == "" {
			it.TVInfo["showType"] = "Unknown"
		}
		it.TVInfo["duration"] = strings.TrimSpace(e.ChildText(".film-detail .fd-infor .fdi-duration"))
		for _, k := range []string{"sub", "dub", "eps"} {
			v := strings.TrimSpace(e.DOM.Find(".tick .tick-" + k).Text())
			if v != "" {
				it.TVInfo[k] = v
			}
		}
		it.AdultContent = strings.Contains(strings.TrimSpace(e.ChildText(".film-poster>.tick-rate")), "18+")
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
