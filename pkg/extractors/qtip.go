package extractors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

type QtipData struct {
	Title         string   `json:"title"`
	Rating        string   `json:"rating"`
	Quality       string   `json:"quality"`
	SubCount      string   `json:"subCount"`
	DubCount      string   `json:"dubCount"`
	EpisodeCount  string   `json:"episodeCount"`
	Type          string   `json:"type"`
	Description   string   `json:"description"`
	JapaneseTitle string   `json:"japaneseTitle"`
	Synonyms      string   `json:"Synonyms"`
	AiredDate     string   `json:"airedDate"`
	Status        string   `json:"status"`
	Genres        []string `json:"genres"`
	WatchLink     string   `json:"watchLink"`
}

func ExtractQtip(id string, baseURL string) (QtipData, error) {
	c := utils.NewCollector()
	var out QtipData
	url := fmt.Sprintf("https://%s/ajax/movie/qtip/%s", baseURL, id)
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("x-requested-with", "XMLHttpRequest")
	})
	c.OnResponse(func(r *colly.Response) {
		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(r.Body))
		if err != nil {
			return
		}
		out.Title = strings.TrimSpace(doc.Find(".pre-qtip-title").Text())
		out.Rating = strings.TrimSpace(doc.Find(".pqd-li i.fas.fa-star").Parent().Text())
		out.Quality = strings.TrimSpace(doc.Find(".tick-item.tick-quality").Text())
		out.SubCount = strings.TrimSpace(doc.Find(".tick-item.tick-sub").Text())
		out.DubCount = strings.TrimSpace(doc.Find(".tick-item.tick-dub").Text())
		out.EpisodeCount = strings.TrimSpace(doc.Find(".tick-item.tick-eps").Text())
		out.Type = strings.TrimSpace(doc.Find(".badge.badge-quality").Text())
		out.Description = strings.TrimSpace(doc.Find(".pre-qtip-description").Text())
		out.JapaneseTitle = strings.TrimSpace(doc.Find(".pre-qtip-line:contains('Japanese:') .stick-text").Text())
		out.AiredDate = strings.TrimSpace(doc.Find(".pre-qtip-line:contains('Aired:') .stick-text").Text())
		out.Status = strings.TrimSpace(doc.Find(".pre-qtip-line:contains('Status:') .stick-text").Text())
		out.Synonyms = strings.TrimSpace(doc.Find(".pre-qtip-line:contains('Synonyms:') .stick-text").Text())
		out.Genres = nil
		doc.Find(".pre-qtip-line:contains('Genres:') a").Each(func(_ int, sel *goquery.Selection) {
			genre := strings.TrimSpace(sel.Text())
			if genre != "" {
				out.Genres = append(out.Genres, genre)
			}
		})
		out.WatchLink = doc.Find(".pre-qtip-button a.btn.btn-play").AttrOr("href", "")
	})
	var errVisit error
	c.OnError(func(_ *colly.Response, err error) { errVisit = err })
	_ = c.Visit(url)
	c.Wait()
	if errVisit != nil {
		return QtipData{}, errVisit
	}
	// Ensure genres non-nil for JSON
	if out.Genres == nil {
		out.Genres = []string{}
	}
	_, _ = json.Marshal(out)
	return out, nil
}
