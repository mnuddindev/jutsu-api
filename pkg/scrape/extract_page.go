package scrape

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type TVInfo struct {
	ShowType  string `json:"showType"`
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
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, 0, err
	}
	// totalPages
	lastSel := doc.Find(`.pre-pagination nav .pagination > .page-item a[title="Last"]`)
	nextSel := doc.Find(`.pre-pagination nav .pagination > .page-item a[title="Next"]`)
	activeText := strings.TrimSpace(doc.Find(".pre-pagination nav .pagination > .page-item.active a").Text())
	var totalPages int
	if href, ok := lastSel.Attr("href"); ok {
		parts := strings.Split(href, "=")
		last := parts[len(parts)-1]
		totalPages, _ = strconv.Atoi(last)
	}
	if totalPages == 0 {
		if href, ok := nextSel.Attr("href"); ok {
			parts := strings.Split(href, "=")
			last := parts[len(parts)-1]
			totalPages, _ = strconv.Atoi(last)
		}
	}
	if totalPages == 0 {
		totalPages, _ = strconv.Atoi(activeText)
		if totalPages == 0 {
			totalPages = 1
		}
	}
	contentSelector := ".tab-content"
	if !strings.Contains(params, "az-list") {
		contentSelector = "#main-content"
	}
	var items []ExtractedItem
	doc.Find(contentSelector + " .film_list-wrap .flw-item").Each(func(i int, s *goquery.Selection) {
		fdi := s.Find(".film-detail .fd-infor .fdi-item")
		var showType string
		fdi.EachWithBreak(func(_ int, fs *goquery.Selection) bool {
			t := strings.ToLower(strings.TrimSpace(fs.Text()))
			for _, typ := range []string{"tv", "ona", "movie", "ova", "special"} {
				if strings.Contains(t, typ) {
					showType = strings.TrimSpace(fs.Text())
					return false
				}
			}
			return true
		})
		poster, _ := s.Find(".film-poster>img").Attr("data-src")
		title := s.Find(".film-detail .film-name").Text()
		jtitle, _ := s.Find(".film-detail>.film-name>a").Attr("data-jname")
		desc := strings.TrimSpace(s.Find(".film-detail .description").Text())
		dataID, _ := s.Find(".film-poster>a").Attr("data-id")
		idHref, _ := s.Find(".film-poster>a").Attr("href")
		id := idHref
		if parts := strings.Split(idHref, "/"); len(parts) > 0 {
			id = parts[len(parts)-1]
		}
		tv := TVInfo{ShowType: strings.TrimSpace(showType), Duration: strings.TrimSpace(s.Find(".film-detail .fd-infor .fdi-duration").Text())}
		for _, prop := range []string{"sub", "dub", "eps"} {
			val := strings.TrimSpace(s.Find(".tick .tick-" + prop).Text())
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
		adult := false
		if strings.Contains(strings.TrimSpace(s.Find(".film-poster>.tick-rate").Text()), "18+") {
			adult = true
		}
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
	return items, totalPages, nil
}
