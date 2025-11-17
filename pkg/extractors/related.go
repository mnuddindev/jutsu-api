package extractors

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

func ExtractRelatedData(id, baseURL string) ([]SidebarItem, error) {
	url := fmt.Sprintf("https://%s/%s", baseURL, id)
	c := utils.NewCollector()
	var items []SidebarItem
	selector := `#main-sidebar .block_area:has(.cat-heading:contains("Related Anime")) .anif-block-ul .ulclear li`
	c.OnHTML(selector, func(e *colly.HTMLElement) {
		var it SidebarItem
		it.ID = lastSegment(e.ChildAttr(".film-detail .film-name a", "href"))
		it.DataID = e.ChildAttr(".film-poster", "data-id")
		it.Title = strings.TrimSpace(e.ChildText(".film-detail .film-name a"))
		it.JapaneseTitle = strings.TrimSpace(e.DOM.Find(".film-detail .film-name a").AttrOr("data-jname", ""))
		it.Poster = firstNonEmpty(e.ChildAttr(".film-poster img", "data-src"), e.ChildAttr(".film-poster img", "src"))
		// showType
		showTypeText := strings.ToLower(strings.TrimSpace(e.DOM.Find(".tick").Text()))
		it.TVInfo = map[string]string{}
		for _, typ := range []string{"TV", "ONA", "Movie", "OVA", "Special"} {
			if strings.Contains(showTypeText, strings.ToLower(typ)) {
				it.TVInfo["showType"] = typ
				break
			}
		}
		if it.TVInfo["showType"] == "" {
			it.TVInfo["showType"] = "Unknown"
		}
		for _, k := range []string{"sub", "dub", "eps"} {
			v := strings.TrimSpace(e.DOM.Find(".tick-item.tick-" + k).Text())
			if v != "" {
				it.TVInfo[k] = v
			}
		}
		it.AdultContent = strings.Contains(strings.TrimSpace(e.ChildText(".film-poster > .tick-rate")), "18+")
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
