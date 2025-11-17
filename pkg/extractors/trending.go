package extractors

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

type TrendingItem struct {
	ID            string `json:"id"`
	DataID        string `json:"data_id"`
	Number        string `json:"number"`
	Poster        string `json:"poster"`
	Title         string `json:"title"`
	JapaneseTitle string `json:"japanese_title"`
}

func ExtractTrending(baseURL string) ([]TrendingItem, error) {
	c := utils.NewCollector()
	var items []TrendingItem
	url := fmt.Sprintf("https://%s/home", baseURL)

	c.OnHTML("#anime-trending #trending-home .swiper-slide", func(e *colly.HTMLElement) {
		var it TrendingItem
		it.DataID = e.Attr("data-id")
		it.Number = strings.TrimSpace(e.ChildText(".number > span"))
		it.Poster = e.ChildAttr("img", "data-src")
		it.Title = strings.TrimSpace(e.ChildText(".film-title"))
		it.JapaneseTitle = strings.TrimSpace(e.DOM.Find(".film-title").AttrOr("data-jname", ""))
		it.ID = lastSegment(e.ChildAttr("a", "href"))
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
