package extractors

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

type SpotlightItem struct {
	ID            string                 `json:"id"`
	DataID        string                 `json:"data_id"`
	Poster        string                 `json:"poster"`
	Title         string                 `json:"title"`
	JapaneseTitle string                 `json:"japanese_title"`
	Description   string                 `json:"description"`
	TVInfo        map[string]interface{} `json:"tvInfo"`
}

func ExtractSpotlights(baseURL string) ([]SpotlightItem, error) {
	c := utils.NewCollector()
	var items []SpotlightItem
	url := fmt.Sprintf("https://%s/home", baseURL)
	c.OnHTML("div.deslide-wrap > div.container > div#slider > div.swiper-wrapper > div.swiper-slide", func(e *colly.HTMLElement) {
		var it SpotlightItem
		it.Poster = e.ChildAttr("div.deslide-item > div.deslide-cover > div.deslide-cover-img > img.film-poster-img", "data-src")
		it.Title = strings.TrimSpace(e.ChildText("div.deslide-item > div.deslide-item-content > div.desi-head-title"))
		it.JapaneseTitle = strings.TrimSpace(e.DOM.Find("div.deslide-item > div.deslide-item-content > div.desi-head-title").AttrOr("data-jname", ""))
		it.Description = strings.TrimSpace(e.ChildText("div.deslide-item > div.deslide-item-content > div.desi-description"))
		btn := e.DOM.Find(".deslide-item > .deslide-item-content > .desi-buttons > a").First()
		it.ID = lastSegment(btn.AttrOr("href", ""))
		it.DataID = lastSegment(it.ID)
		it.TVInfo = map[string]interface{}{}
		e.ForEach("div.sc-detail > div.scd-item", func(idx int, el *colly.HTMLElement) {
			key := map[int]string{0: "showType", 1: "duration", 2: "releaseDate", 3: "quality", 4: "episodeInfo"}[idx]
			if key == "" {
				return
			}
			val := strings.ReplaceAll(strings.TrimSpace(el.Text), "\n", "")
			if el.DOM.Find(".tick").Length() > 0 {
				valSub := strings.TrimSpace(el.ChildText(".tick-sub"))
				valDub := strings.TrimSpace(el.ChildText(".tick-dub"))
				it.TVInfo[key] = map[string]string{"sub": valSub, "dub": valDub}
				return
			}
			it.TVInfo[key] = val
		})
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
