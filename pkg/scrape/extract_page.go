package scrape

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

type TVInfo struct {
	ShowType string `json:"showType"`
	Duration string `json:"duration"`
	Sub      string `json:"sub,omitempty"`
	Dub      string `json:"dub,omitempty"`
	Eps      string `json:"eps,omitempty"`
}

type ExtractedItem struct {
	ID            string `json:"id"`
	DataID        string `json:"data_id"`
	Poster        string `json:"poster"`
	Title         string `json:"title"`
	JapaneseTitle string `json:"japanese_title"`
	Description   string `json:"description"`
	TVInfo        TVInfo `json:"tvInfo"`
	AdultContent  bool   `json:"adultContent"`
}

func ExtractPage(page int, params, baseURL string) ([]ExtractedItem, int, error) {
	url := fmt.Sprintf("https://%s/%s?page=%d", baseURL, params, page)
	c := utils.NewCollector()

	items := make([]ExtractedItem, 0, 50)
	totalPages := 1

	c.OnHTML(`.pre-pagination nav .pagination > .page-item a[title="Last"]`, func(e *colly.HTMLElement) {
		if href := e.Attr("href"); href != "" {
			parts := strings.Split(href, "=")
			last := parts[len(parts)-1]
			if n, err := strconv.Atoi(last); err == nil && n > 0 {
				totalPages = n
			}
		}
	})
	c.OnHTML(`.pre-pagination nav .pagination > .page-item a[title="Next"]`, func(e *colly.HTMLElement) {
		if totalPages != 1 {
			return
		}
		if href := e.Attr("href"); href != "" {
			parts := strings.Split(href, "=")
			last := parts[len(parts)-1]
			if n, err := strconv.Atoi(last); err == nil && n > 0 {
				totalPages = n
			}
		}
	})
	c.OnHTML(`.pre-pagination nav .pagination > .page-item.active a`, func(e *colly.HTMLElement) {
		if totalPages != 1 {
			return
		}
		if n, err := strconv.Atoi(strings.TrimSpace(e.Text)); err == nil && n > 0 {
			totalPages = n
		}
	})

	contentSelector := ".tab-content"
	if !strings.Contains(params, "az-list") {
		contentSelector = "#main-content"
	}
	c.OnHTML(contentSelector+" .film_list-wrap .flw-item", func(e *colly.HTMLElement) {
		// showType
		showType := "Unknown"
		e.ForEach(".film-detail .fd-infor .fdi-item", func(_ int, el *colly.HTMLElement) {
			t := strings.ToLower(strings.TrimSpace(el.Text))
			for _, typ := range []string{"tv", "ona", "movie", "ova", "special"} {
				if strings.Contains(t, typ) {
					showType = strings.TrimSpace(el.Text)
				}
			}
		})
		poster := e.ChildAttr(".film-poster>img", "data-src")
		title := e.ChildText(".film-detail .film-name")
		jtitle := e.ChildAttr(".film-detail>.film-name>a", "data-jname")
		desc := strings.TrimSpace(e.ChildText(".film-detail .description"))
		dataID := e.ChildAttr(".film-poster>a", "data-id")
		idHref := e.ChildAttr(".film-poster>a", "href")
		id := idHref
		if parts := strings.Split(idHref, "/"); len(parts) > 0 {
			id = parts[len(parts)-1]
		}
		tv := TVInfo{ShowType: strings.TrimSpace(showType), Duration: strings.TrimSpace(e.ChildText(".film-detail .fd-infor .fdi-duration"))}
		for _, prop := range []string{"sub", "dub", "eps"} {
			val := strings.TrimSpace(e.ChildText(".tick .tick-" + prop))
			if val != "" {
				switch prop {
				case "sub":
					tv.Sub = val
				case "dub":
					tv.Dub = val
				case "eps":
					tv.Eps = val
				}
			}
		}
		adult := strings.Contains(strings.TrimSpace(e.ChildText(".film-poster>.tick-rate")), "18+")
		items = append(items, ExtractedItem{
			ID:            id,
			DataID:        dataID,
			Poster:        poster,
			Title:         title,
			JapaneseTitle: jtitle,
			Description:   desc,
			TVInfo:        tv,
			AdultContent:  adult,
		})
	})

	var visitErr error
	c.OnError(func(_ *colly.Response, err error) { visitErr = err })
	_ = c.Visit(url)
	c.Wait()
	if visitErr != nil {
		return nil, 0, visitErr
	}
	return items, totalPages, nil
}
