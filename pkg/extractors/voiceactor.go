package extractors

import (
	"encoding/json"
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
	url := fmt.Sprintf("https://%s/ajax/character/list/%s?page=%d", baseURL, utils.ExtractDataID(id), page)

	c.OnResponse(func(r *colly.Response) {
		var res struct {
			Status        bool   `json:"status"`
			HTML          string `json:"html"`
			TotalItems    int    `json:"totalItems"`
			ContinueWatch string `json:"continueWatch"`
		}
		if err := json.Unmarshal(r.Body, &res); err != nil {
			fmt.Println("JSON parse error:", err)
			return
		}

		body := utils.CleanHTML(res.HTML)

		doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
		if err != nil {
			fmt.Println("goquery error:", err)
			return
		}

		pageNum := 1

		// pagination
		lastLi := doc.Find(".pre-pagination nav ul li").Last().Find("a")
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

		// Characters & Voice Actors
		doc.Find(".bac-list-wrap .bac-item").Each(func(_ int, charSel *goquery.Selection) {
			var item CharactersVoiceActors

			item.Character.ID = seg(charSel.Find(".per-info.ltr .pi-avatar").AttrOr("href", ""), 2)
			item.Character.Poster = charSel.Find(".per-info.ltr .pi-avatar img").AttrOr("data-src", "")
			item.Character.Name = charSel.Find(".per-info.ltr .pi-detail a").Text()
			item.Character.Cast = charSel.Find(".per-info.ltr .pi-detail .pi-cast").Text()

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

			out.CharactersVoiceActors = append(out.CharactersVoiceActors, item)
		})
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
