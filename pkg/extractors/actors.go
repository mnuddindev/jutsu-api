package extractors

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

func ExtractVoiceActor(id string, baseURL string) (VoiceActorResult, error) {
	c := utils.NewCollector()
	var result VoiceActorResult
	result.Success = true
	var data VoiceActorData

	url := fmt.Sprintf("https://%s/people/%s", baseURL, id)

	c.OnHTML(".apw-detail .name", func(e *colly.HTMLElement) {
		data.Name = strings.TrimSpace(e.Text)
	})
	c.OnHTML(".apw-detail .sub-name", func(e *colly.HTMLElement) {
		data.JapaneseName = strings.TrimSpace(e.Text)
	})
	c.OnHTML(".avatar-circle img", func(e *colly.HTMLElement) {
		data.Profile = e.Attr("src")
	})
	c.OnHTML("#bio .bio", func(e *colly.HTMLElement) {
		data.About.Description = strings.TrimSpace(e.Text)
	})
	c.OnHTML("#bio .bio", func(e *colly.HTMLElement) {
		if html, err := e.DOM.Html(); err == nil {
			data.About.Style = html
		}
	})

	c.OnHTML(".bac-list-wrap .bac-item", func(e *colly.HTMLElement) {
		var role Role
		animeSel := e.DOM.Find(".per-info.anime-info.ltr")
		charSel := e.DOM.Find(".per-info.rtl")
		// Anime
		role.Anime.ID = strings.TrimPrefix(strings.TrimSpace(animeSel.Find(".pi-name a").AttrOr("href", "")), "/")
		if parts := strings.Split(role.Anime.ID, "/"); len(parts) > 0 {
			role.Anime.ID = parts[len(parts)-1]
		}
		role.Anime.Title = strings.TrimSpace(animeSel.Find(".pi-name a").Text())
		role.Anime.Poster = animeSel.Find(".pi-avatar img").AttrOr("data-src", animeSel.Find(".pi-avatar img").AttrOr("src", ""))
		cast := strings.TrimSpace(animeSel.Find(".pi-cast").Text())
		if cast != "" {
			sp := strings.Split(cast, ",")
			role.Anime.Type = strings.TrimSpace(sp[0])
			if len(sp) > 1 {
				role.Anime.Year = strings.TrimSpace(sp[1])
			}
		}
		// Character
		role.Character.ID = strings.TrimPrefix(strings.TrimSpace(charSel.Find(".pi-name a").AttrOr("href", "")), "/")
		if parts := strings.Split(role.Character.ID, "/"); len(parts) > 0 {
			role.Character.ID = parts[len(parts)-1]
		}
		role.Character.Name = strings.TrimSpace(charSel.Find(".pi-name a").Text())
		role.Character.Profile = charSel.Find(".pi-avatar img").AttrOr("data-src", charSel.Find(".pi-avatar img").AttrOr("src", ""))
		role.Character.Role = strings.TrimSpace(charSel.Find(".pi-cast").Text())
		data.Roles = append(data.Roles, role)
	})

	var visitErr error
	c.OnError(func(_ *colly.Response, err error) { visitErr = err })
	_ = c.Visit(url)
	c.Wait()
	if visitErr != nil {
		return VoiceActorResult{}, visitErr
	}
	data.ID = id
	result.Results.Data = []VoiceActorData{data}
	return result, nil
}
