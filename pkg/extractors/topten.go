package extractors

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
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

	c.OnHTML("#main-sidebar .block_area-realtime .block_area-content", func(e *colly.HTMLElement) {
		// Get all UL elements
		uls := e.DOM.Find("ul")

		uls.Each(func(ulIndex int, ul *goquery.Selection) {
			var label string
			switch ulIndex {
			case 0:
				label = "today"
			case 1:
				label = "week"
			case 2:
				label = "month"
			default:
				return
			}

			ul.Find("li").Each(func(liIndex int, li *goquery.Selection) {
				if liIndex >= 10 { // Only take top 10
					return
				}

				item := extractItemFromSelection(li)
				if item.Title != "" {
					res[label] = append(res[label], item)
				}
			})
		})
	})

	var errVisit error
	c.OnError(func(r *colly.Response, err error) {
		errVisit = fmt.Errorf("scraping error: %w", err)
	})

	err := c.Visit(url)
	if err != nil {
		return nil, fmt.Errorf("visit failed: %w", err)
	}

	c.Wait()

	if errVisit != nil {
		return nil, errVisit
	}

	return res, nil
}

func extractItemFromSelection(s *goquery.Selection) TopTenItem {
	item := TopTenItem{
		TVInfo: make(map[string]string),
	}

	item.Number = strings.TrimSpace(s.Find(".film-number span").Text())
	item.Title = strings.TrimSpace(s.Find(".film-detail .film-name a").Text())
	item.Poster = s.Find(".film-poster img").AttrOr("data-src", "")
	item.JapaneseTitle = strings.TrimSpace(s.Find(".film-detail .film-name a").AttrOr("data-jname", ""))
	item.DataID = s.Find(".film-poster").AttrOr("data-id", "")

	href := s.Find(".film-detail .film-name a").AttrOr("href", "")
	if href != "" {
		parts := strings.Split(strings.Trim(href, "/"), "/")
		if len(parts) > 0 {
			item.ID = parts[len(parts)-1]
		}
	}

	for _, prop := range []string{"sub", "dub", "eps"} {
		value := strings.TrimSpace(s.Find(fmt.Sprintf(".tick .tick-%s", prop)).Text())
		if value != "" {
			item.TVInfo[prop] = value
		}
	}

	return item
}
