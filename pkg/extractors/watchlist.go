package extractors

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

type WatchlistItem struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Poster   string `json:"poster"`
	Duration string `json:"duration"`
	Type     string `json:"type"`
	SubCount string `json:"subCount"`
	DubCount string `json:"dubCount"`
	Link     string `json:"link"`
	ShowType string `json:"showType"`
	TVInfo   struct {
		ShowType string `json:"showType"`
		Duration string `json:"duration"`
		Sub      string `json:"sub"`
		Dub      string `json:"dub"`
	} `json:"tvInfo"`
}

type WatchlistResult struct {
	Watchlist  []WatchlistItem `json:"watchlist"`
	TotalPages int             `json:"totalPages"`
}

func ExtractWatchlist(userID string, page int, baseURL string) (WatchlistResult, error) {
	c := utils.NewCollector()
	var out WatchlistResult
	url := fmt.Sprintf("https://%s/community/user/%s/watch-list?page=%d", baseURL, userID, page)

	// total pages via Last/Next/Active
	c.OnHTML(`.pre-pagination nav .pagination > .page-item a[title="Last"]`, func(e *colly.HTMLElement) {
		if href := e.Attr("href"); href != "" {
			parts := strings.Split(href, "=")
			last := parts[len(parts)-1]
			if n, err := strconv.Atoi(last); err == nil {
				out.TotalPages = n
			}
		}
	})
	c.OnHTML(`.pre-pagination nav .pagination > .page-item a[title="Next"]`, func(e *colly.HTMLElement) {
		if out.TotalPages == 0 {
			if href := e.Attr("href"); href != "" {
				parts := strings.Split(href, "=")
				last := parts[len(parts)-1]
				if n, err := strconv.Atoi(last); err == nil {
					out.TotalPages = n
				}
			}
		}
	})
	c.OnHTML(`.pre-pagination nav .pagination > .page-item.active a`, func(e *colly.HTMLElement) {
		if out.TotalPages == 0 {
			if n, err := strconv.Atoi(strings.TrimSpace(e.Text)); err == nil {
				out.TotalPages = n
			}
		}
	})
	if out.TotalPages == 0 {
		out.TotalPages = 1
	}

	c.OnHTML(".flw-item", func(e *colly.HTMLElement) {
		var item WatchlistItem
		item.Title = strings.TrimSpace(e.ChildText(".film-name a"))
		item.Poster = e.ChildAttr(".film-poster img", "data-src")
		item.Duration = strings.TrimSpace(e.ChildText(".fdi-duration"))
		item.Type = strings.TrimSpace(e.ChildText(".fdi-item"))
		item.SubCount = strings.TrimSpace(e.DOM.Find(".tick-item.tick-sub").Text())
		item.DubCount = strings.TrimSpace(e.DOM.Find(".tick-item.tick-dub").Text())
		link := e.ChildAttr(".film-name a", "href")
		item.ID = lastSegment(link)
		item.Link = fmt.Sprintf("https://%s%s", baseURL, link)
		item.ShowType = item.Type
		item.TVInfo.ShowType = item.Type
		item.TVInfo.Duration = item.Duration
		item.TVInfo.Sub = item.SubCount
		item.TVInfo.Dub = item.DubCount
		out.Watchlist = append(out.Watchlist, item)
	})

	var visitErr error
	c.OnError(func(_ *colly.Response, err error) { visitErr = err })
	_ = c.Visit(url)
	c.Wait()
	if visitErr != nil {
		return WatchlistResult{}, visitErr
	}
	if out.TotalPages == 0 {
		out.TotalPages = 1
	}
	return out, nil
}
