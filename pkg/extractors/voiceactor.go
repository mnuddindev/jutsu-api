package extractors

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

type VoiceActorPageResult struct {
	TotalPages            int                     `json:"totalPages"`
	CharactersVoiceActors []CharactersVoiceActors `json:"charactersVoiceActors"`
}

func ExtractVoiceActorPage(id string, page int, baseURL string) (VoiceActorPageResult, error) {
	c := utils.NewCollector()
	var out VoiceActorPageResult
	url := fmt.Sprintf("https://%s/ajax/character/list/%s?page=%d", baseURL, lastSegment(id), page)

	// pagination
	c.OnHTML(".pre-pagination nav ul", func(e *colly.HTMLElement) {
		pageNum := 1
		// last li a data-url or text
		lastLi := e.DOM.Find("li").Last().Find("a")
		if v, ok := lastLi.Attr("data-url"); ok {
			re := regexp.MustCompile(`page=(\d+)`)
			if m := re.FindStringSubmatch(v); len(m) == 2 {
				if n, err := strconv.Atoi(m[1]); err == nil {
					pageNum = n
				}
			}
		} else {
			t := strings.TrimSpace(lastLi.Text())
			if n, err := strconv.Atoi(t); err == nil {
				pageNum = n
			}
		}
		out.TotalPages = pageNum
	})

	c.OnHTML(".bac-list-wrap .bac-item", func(e *colly.HTMLElement) {
		var item CharactersVoiceActors
		// character
		charSel := e.DOM
		item.Character.ID = seg(charSel.Find(".per-info.ltr .pi-avatar").AttrOr("href", ""), 2)
		item.Character.Poster = charSel.Find(".per-info.ltr .pi-avatar img").AttrOr("data-src", "")
		item.Character.Name = charSel.Find(".per-info.ltr .pi-detail a").Text()
		item.Character.Cast = charSel.Find(".per-info.ltr .pi-detail .pi-cast").Text()
		// voice actors (rtl)
		rtl := charSel.Find(".per-info.rtl")
		if rtl.Length() > 0 {
			rtl.Find(".pi-avatar").Each(func(_ int, s *goquery.Selection) {
				item.VoiceActors = append(item.VoiceActors, struct {
					ID     string `json:"id"`
					Poster string `json:"poster"`
					Name   string `json:"name"`
				}{
					ID:     lastSegment(s.AttrOr("href", "")),
					Poster: s.Find("img").AttrOr("data-src", ""),
					Name:   strings.TrimSpace(s.Parent().Find(".pi-detail .pi-name a").Text()),
				})
			})
		} else {
			charSel.Find(".per-info.per-info-xx .pix-list .pi-avatar").Each(func(_ int, s *goquery.Selection) {
				item.VoiceActors = append(item.VoiceActors, struct {
					ID     string `json:"id"`
					Poster string `json:"poster"`
					Name   string `json:"name"`
				}{
					ID:     lastSegment(s.AttrOr("href", "")),
					Poster: s.Find("img").AttrOr("data-src", ""),
					Name:   strings.TrimSpace(s.AttrOr("title", "")),
				})
			})
		}
		if len(item.VoiceActors) == 0 {
			charSel.Find(".per-info.per-info-xx .pix-list .pi-avatar").Each(func(_ int, s *goquery.Selection) {
				item.VoiceActors = append(item.VoiceActors, struct {
					ID     string `json:"id"`
					Poster string `json:"poster"`
					Name   string `json:"name"`
				}{
					ID:     seg(s.AttrOr("href", ""), 2),
					Poster: s.Find("img").AttrOr("data-src", ""),
					Name:   strings.TrimSpace(s.AttrOr("title", "")),
				})
			})
		}
		out.CharactersVoiceActors = append(out.CharactersVoiceActors, item)
	})

	var errVisit error
	c.OnError(func(_ *colly.Response, err error) { errVisit = err })
	_ = c.Visit(url)
	c.Wait()
	if errVisit != nil {
		return VoiceActorPageResult{}, errVisit
	}
	return out, nil
}

func lastSegment(href string) string {
	parts := strings.Split(strings.Trim(href, "/"), "/")
	if len(parts) == 0 {
		return ""
	}
	return parts[len(parts)-1]
}

func seg(href string, idx int) string {
	parts := strings.Split(strings.Trim(href, "/"), "/")
	if idx >= 0 && idx < len(parts) {
		return parts[idx]
	}
	return ""
}
