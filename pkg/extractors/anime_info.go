package extractors

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/mnuddindev/jutsu-api/pkg/httpclient"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

type Trailer struct {
	Title     string `json:"title"`
	URL       string `json:"url"`
	Thumbnail string `json:"thumbnail"`
}

type AnimeInfo struct {
	AdultContent          bool                    `json:"adultContent"`
	DataID                string                  `json:"data_id"`
	ID                    string                  `json:"id"`
	AnilistID             string                  `json:"anilistId"`
	MalID                 string                  `json:"malId"`
	Title                 string                  `json:"title"`
	JapaneseTitle         string                  `json:"japanese_title"`
	Synonyms              string                  `json:"synonyms"`
	Poster                string                  `json:"poster"`
	ShowType              string                  `json:"showType"`
	AnimeInfo             map[string]interface{}  `json:"animeInfo"`
	TVInfo                map[string]string       `json:"tvInfo"`
	Trailers              []Trailer               `json:"trailers"`
	CharactersVoiceActors []CharactersVoiceActors `json:"charactersVoiceActors"`
	RecommendedData       []RecommendItem         `json:"recommended_data"`
	RelatedData           []SidebarItem           `json:"related_data"`
	PopularData           []SidebarItem           `json:"popular_data"`
}

func ExtractAnimeInfo(id, baseURL string) (AnimeInfo, error) {
	info := AnimeInfo{
		ID:        id,
		AnimeInfo: make(map[string]interface{}),
		TVInfo:    make(map[string]string),
	}

	mainURL := fmt.Sprintf("https://%s/%s", baseURL, id)
	c := utils.NewCollector()

	c.OnHTML("#ani_detail .film-name", func(e *colly.HTMLElement) {
		info.Title = strings.TrimSpace(e.Text)
		info.JapaneseTitle = e.Attr("data-jname")
	})
	c.OnHTML("#ani_detail .prebreadcrumb ol li:nth-child(2) a", func(e *colly.HTMLElement) {
		info.ShowType = strings.TrimSpace(e.Text)
	})
	c.OnHTML("#ani_detail .film-poster img", func(e *colly.HTMLElement) {
		info.Poster = e.Attr("src")
	})
	c.OnHTML("#ani_detail .film-stats .tick-item", func(e *colly.HTMLElement) {
		text := strings.TrimSpace(e.Text)
		classes := e.Attr("class")
		switch {
		case strings.Contains(classes, "tick-quality"):
			info.TVInfo["quality"] = text
		case strings.Contains(classes, "tick-sub"):
			info.TVInfo["sub"] = text
		case strings.Contains(classes, "tick-dub"):
			info.TVInfo["dub"] = text
		case strings.Contains(classes, "tick-eps"):
			info.TVInfo["eps"] = text
		case strings.Contains(classes, "tick-pg"):
			info.TVInfo["rating"] = text
		}
	})
	c.OnHTML("#ani_detail .film-stats span.item", func(e *colly.HTMLElement) {
		text := strings.TrimSpace(e.Text)
		if info.TVInfo["showType"] == "" {
			info.TVInfo["showType"] = text
		} else if info.TVInfo["duration"] == "" {
			info.TVInfo["duration"] = text
		}
	})
	c.OnHTML("#ani_detail .film-description .text", func(e *colly.HTMLElement) {
		info.AnimeInfo["Overview"] = strings.TrimSpace(e.Text)
	})
	c.OnHTML("#ani_detail .anisc-info .item", func(e *colly.HTMLElement) {
		head := strings.TrimSpace(e.ChildText(".item-head"))
		if head == "" {
			return
		}
		head = strings.TrimSuffix(head, ":")
		if head == "Genres" || head == "Producers" {
			values := []string{}
			e.DOM.Find("a").Each(func(_ int, sel *goquery.Selection) {
				values = append(values, strings.ReplaceAll(strings.TrimSpace(sel.Text()), " ", "-"))
			})
			info.AnimeInfo[head] = values
		} else {
			val := strings.ReplaceAll(strings.TrimSpace(e.ChildText(".name")), " ", "-")
			info.AnimeInfo[head] = val
		}
	})
	c.OnHTML(".item.item-title:has(.item-head:contains('Synonyms')) .name", func(e *colly.HTMLElement) {
		info.Synonyms = strings.TrimSpace(e.Text)
	})
	c.OnHTML(".tick-rate", func(e *colly.HTMLElement) {
		if strings.Contains(strings.TrimSpace(e.Text), "18+") {
			info.AdultContent = true
		}
	})
	c.OnHTML("#syncData", func(e *colly.HTMLElement) {
		var sync struct {
			AnilistID interface{} `json:"anilist_id"`
			MalID     interface{} `json:"mal_id"`
		}
		if err := json.Unmarshal([]byte(e.Text), &sync); err == nil {
			info.AnilistID = toString(sync.AnilistID)
			info.MalID = toString(sync.MalID)
		}
	})
	c.OnHTML(".block_area-promotions-list .screen-items .item", func(e *colly.HTMLElement) {
		title := e.Attr("data-title")
		url := e.Attr("data-src")
		if url == "" {
			return
		}
		full := url
		if strings.HasPrefix(url, "//") {
			full = "https:" + url
		}
		tr := Trailer{Title: title, URL: full}
		if match := embedID(full); match != "" {
			tr.Thumbnail = fmt.Sprintf("https://img.youtube.com/vi/%s/hqdefault.jpg", match)
		}
		info.Trailers = append(info.Trailers, tr)
	})

	var errVisit error
	c.OnError(func(_ *colly.Response, err error) { errVisit = err })
	_ = c.Visit(mainURL)
	c.Wait()
	if errVisit != nil {
		return AnimeInfo{}, errVisit
	}

	info.DataID = lastSegment(id)
	if info.TVInfo["showType"] == "" {
		info.TVInfo["showType"] = info.ShowType
	}
	info.ID = utils.FormatTitle(info.Title, info.DataID)

	// Fetch characters AJAX
	charactersHTML, err := fetchCharacterHTML(info.DataID, baseURL)
	if err != nil {
		return AnimeInfo{}, err
	}
	if charactersHTML != "" {
		list, err := parseCharacterList(charactersHTML)
		if err == nil {
			info.CharactersVoiceActors = list
		}
	}

	recommended, _ := ExtractRecommendedData(id, baseURL)
	related, _ := ExtractRelatedData(id, baseURL)
	popular, _ := ExtractPopularData(id, baseURL)
	info.RecommendedData = recommended
	info.RelatedData = related
	info.PopularData = popular

	return info, nil
}

func fetchCharacterHTML(dataID, baseURL string) (string, error) {
	endpoint := fmt.Sprintf("https://%s/ajax/character/list/%s", baseURL, utils.ExtractDataID(dataID))
	raw, err := httpclient.Get(endpoint)
	if err != nil {
		return "", err
	}
	var payload struct {
		HTML string `json:"html"`
	}
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		return "", err
	}
	return payload.HTML, nil
}

func parseCharacterList(html string) ([]CharactersVoiceActors, error) {
	if strings.TrimSpace(html) == "" {
		return nil, nil
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}
	var list []CharactersVoiceActors
	doc.Find(".bac-list-wrap .bac-item").Each(func(_ int, sel *goquery.Selection) {
		var entry CharactersVoiceActors
		entry.Character.ID = seg(sel.Find(".per-info.ltr .pi-avatar").AttrOr("href", ""), 2)
		entry.Character.Poster = firstNonEmpty(
			sel.Find(".per-info.ltr .pi-avatar img").AttrOr("data-src", ""),
			sel.Find(".per-info.ltr .pi-avatar img").AttrOr("src", ""),
		)
		entry.Character.Name = strings.TrimSpace(sel.Find(".per-info.ltr .pi-detail a").Text())
		entry.Character.Cast = strings.TrimSpace(sel.Find(".per-info.ltr .pi-detail .pi-cast").Text())

		rtl := sel.Find(".per-info.rtl .pi-avatar")
		rtl.Each(func(_ int, h *goquery.Selection) {
			entry.VoiceActors = append(entry.VoiceActors, struct {
				ID     string `json:"id"`
				Poster string `json:"poster"`
				Name   string `json:"name"`
			}{
				ID:     lastSegment(h.AttrOr("href", "")),
				Poster: h.Find("img").AttrOr("data-src", ""),
				Name:   strings.TrimSpace(h.Parent().Find(".pi-detail .pi-name a").Text()),
			})
		})

		if len(entry.VoiceActors) == 0 {
			sel.Find(".per-info.per-info-xx .pix-list .pi-avatar").Each(func(_ int, h *goquery.Selection) {
				entry.VoiceActors = append(entry.VoiceActors, struct {
					ID     string `json:"id"`
					Poster string `json:"poster"`
					Name   string `json:"name"`
				}{
					ID:     lastSegment(h.AttrOr("href", "")),
					Poster: h.Find("img").AttrOr("data-src", ""),
					Name:   strings.TrimSpace(h.AttrOr("title", "")),
				})
			})
		}

		if len(entry.VoiceActors) == 0 {
			sel.Find(".per-info.per-info-xx .pix-list .pi-avatar").Each(func(_ int, h *goquery.Selection) {
				entry.VoiceActors = append(entry.VoiceActors, struct {
					ID     string `json:"id"`
					Poster string `json:"poster"`
					Name   string `json:"name"`
				}{
					ID:     seg(h.AttrOr("href", ""), 2),
					Poster: h.Find("img").AttrOr("data-src", ""),
					Name:   strings.TrimSpace(h.AttrOr("title", "")),
				})
			})
		}

		list = append(list, entry)
	})
	return list, nil
}
