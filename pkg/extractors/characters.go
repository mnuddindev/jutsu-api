package extractors

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

type CharacterData struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Profile      string `json:"profile"`
	JapaneseName string `json:"japaneseName"`
	About        struct {
		Description string `json:"description"`
		Style       string `json:"style"`
	} `json:"about"`
	VoiceActors []struct {
		Name     string `json:"name"`
		Profile  string `json:"profile"`
		Language string `json:"language"`
		ID       string `json:"id"`
	} `json:"voiceActors"`
	Animeography []struct {
		Title         string `json:"title"`
		JapaneseTitle string `json:"japanese_title"`
		ID            string `json:"id"`
		Role          string `json:"role"`
		Type          string `json:"type"`
		Poster        string `json:"poster"`
	} `json:"animeography"`
}

type CharacterResult struct {
	Success bool `json:"success"`
	Results struct {
		Data []CharacterData `json:"data"`
	} `json:"results"`
}

func ExtractCharacter(id string, baseURL string) (CharacterResult, error) {
	c := utils.NewCollector()
	var res CharacterResult
	res.Success = true
	var data CharacterData
	url := fmt.Sprintf("https://%s//character/%s", baseURL, id)

	c.OnHTML(".apw-detail .name", func(e *colly.HTMLElement) { data.Name = strings.TrimSpace(e.Text) })
	c.OnHTML(".apw-detail .sub-name", func(e *colly.HTMLElement) { data.JapaneseName = strings.TrimSpace(e.Text) })
	c.OnHTML(".avatar-circle img", func(e *colly.HTMLElement) { data.Profile = e.Attr("src") })
	c.OnHTML("#bio .bio", func(e *colly.HTMLElement) {
		data.About.Description = strings.TrimSpace(e.Text)
		if html, err := e.DOM.Html(); err == nil {
			data.About.Style = html
		}
	})
	c.OnHTML("#voiactor .per-info", func(e *colly.HTMLElement) {
		va := struct {
			Name     string `json:"name"`
			Profile  string `json:"profile"`
			Language string `json:"language"`
			ID       string `json:"id"`
		}{}
		va.Name = strings.TrimSpace(e.ChildText(".pi-name a"))
		va.Profile = e.ChildAttr(".pi-avatar img", "src")
		va.Language = strings.TrimSpace(e.ChildText(".pi-cast"))
		va.ID = lastSegment(e.ChildAttr(".pi-name a", "href"))
		if va.Name != "" && va.ID != "" {
			data.VoiceActors = append(data.VoiceActors, va)
		}
	})
	c.OnHTML(".anif-block-ul li", func(e *colly.HTMLElement) {
		anchor := e.DOM.Find(".film-name a.dynamic-name")
		title := strings.TrimSpace(anchor.Text())
		jtitle := strings.TrimSpace(anchor.AttrOr("data-jname", ""))
		cid := lastSegment(anchor.AttrOr("href", ""))
		role := strings.TrimSpace(e.ChildText(".fdi-item"))
		typeVal := strings.TrimSpace(e.DOM.Find(".fdi-item").Last().Text())
		poster := e.ChildAttr(".film-poster img", "src")
		if title != "" && cid != "" {
			data.Animeography = append(data.Animeography, struct {
				Title         string `json:"title"`
				JapaneseTitle string `json:"japanese_title"`
				ID            string `json:"id"`
				Role          string `json:"role"`
				Type          string `json:"type"`
				Poster        string `json:"poster"`
			}{
				Title:         title,
				JapaneseTitle: jtitle,
				ID:            cid,
				Role:          strings.TrimSuffix(role, " (Role)"),
				Type:          typeVal,
				Poster:        poster,
			})
		}
	})

	var errVisit error
	c.OnError(func(_ *colly.Response, err error) { errVisit = err })
	_ = c.Visit(url)
	c.Wait()
	if errVisit != nil {
		return CharacterResult{}, errVisit
	}
	data.ID = id
	res.Results.Data = []CharacterData{data}
	return res, nil
}
