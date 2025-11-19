package extractors

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

type EpisodeHTMLResponse struct {
	Status bool   `json:"status"`
	HTML   string `json:"html"`
}

type EpisodeList struct {
	TotalEpisodes int            `json:"totalEpisodes"`
	Episodes      []EpisodeEntry `json:"episodes"`
}

type EpisodeEntry struct {
	EpisodeNo     int    `json:"episode_no"`
	ID            string `json:"id"`
	Title         string `json:"title"`
	JapaneseTitle string `json:"japanese_title"`
	Filler        bool   `json:"filler"`
}

func ExtractEpisodeList(id string, baseURL string) (EpisodeList, error) {
	c := utils.NewCollector()
	var out EpisodeList
	url := fmt.Sprintf("https://%s/ajax/v2/episode/list/%s", baseURL, utils.ExtractDataID(id))
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("X-Requested-With", "XMLHttpRequest")
		r.Headers.Set("Referer", fmt.Sprintf("https://%s/watch/%s", baseURL, id))
	})
	c.OnResponse(func(r *colly.Response) {
		var res EpisodeHTMLResponse

		if err := json.Unmarshal(r.Body, &res); err != nil {
			fmt.Println("JSON parse error:", err)
			return
		}

		body := utils.CleanHTML(res.HTML)
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
		if err != nil {
			return
		}

		doc.Find(".detail-infor-content .ss-list a").Each(func(_ int, sel *goquery.Selection) {
			var ent EpisodeEntry
			ent.ID = lastSegment(sel.AttrOr("href", ""))
			ent.Title = strings.TrimSpace(sel.AttrOr("title", ""))
			ent.JapaneseTitle = sel.Find(".ep-name").AttrOr("data-jname", "")
			ent.Filler = strings.Contains(sel.AttrOr("class", ""), "ssl-item-filler")
			if num, err := strconv.Atoi(sel.AttrOr("data-number", "0")); err == nil {
				ent.EpisodeNo = num
			}
			out.Episodes = append(out.Episodes, ent)
		})
		out.TotalEpisodes = len(out.Episodes)
	})
	var errVisit error
	c.OnError(func(_ *colly.Response, err error) { errVisit = err })
	_ = c.Visit(url)
	c.Wait()
	if errVisit != nil {
		return EpisodeList{}, errVisit
	}
	return out, nil
}
