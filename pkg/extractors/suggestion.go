package extractors

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

type SuggestionItem struct {
	ID            string `json:"id"`
	DataID        string `json:"data_id"`
	Poster        string `json:"poster"`
	Title         string `json:"title"`
	JapaneseTitle string `json:"japanese_title"`
	ReleaseDate   string `json:"releaseDate"`
	ShowType      string `json:"showType"`
	Duration      string `json:"duration"`
}

func ExtractSuggestions(keyword, baseURL string) ([]SuggestionItem, error) {
	c := utils.NewCollector()
	var items []SuggestionItem
	url := fmt.Sprintf("https://%s/ajax/search/suggest?keyword=%s", baseURL, keyword)
	c.OnResponse(func(r *colly.Response) {
		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(r.Body))
		if err != nil {
			return
		}
		doc.Find(".nav-item:not(.nav-bottom)").Each(func(_ int, sel *goquery.Selection) {
			var it SuggestionItem
			href := sel.AttrOr("href", "")
			id := strings.Split(strings.TrimPrefix(href, "/"), "?")[0]
			it.ID = id
			it.DataID = lastSegment(id)
			it.Poster = sel.Find(".film-poster-img").AttrOr("data-src", "")
			it.Title = strings.TrimSpace(sel.Find(".film-name").Text())
			it.JapaneseTitle = strings.TrimSpace(sel.Find(".film-name").AttrOr("data-jname", ""))
			info := sel.Find(".film-infor span")
			if info.Length() > 0 {
				it.ReleaseDate = strings.TrimSpace(info.First().Text())
				it.Duration = strings.TrimSpace(info.Last().Text())
			}
			if html, err := sel.Find(".film-infor").Html(); err == nil {
				if idx := strings.Index(html, "<i class=\"dot\"></i>"); idx >= 0 {
					parts := strings.Split(html, "<i class=\"dot\"></i>")
					if len(parts) >= 3 {
						it.ShowType = strings.TrimSpace(stripTags(parts[1]))
					}
				}
			}
			items = append(items, it)
		})
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

func stripTags(s string) string {
	res := s
	for {
		start := strings.Index(res, "<")
		if start < 0 {
			break
		}
		end := strings.Index(res[start:], ">")
		if end < 0 {
			break
		}
		res = res[:start] + res[start+end+1:]
	}
	return res
}
