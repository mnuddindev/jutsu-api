package extractors

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

type ScheduleItem struct {
	ID            string `json:"id"`
	DataID        string `json:"data_id"`
	Title         string `json:"title"`
	JapaneseTitle string `json:"japanese_title"`
	ReleaseDate   string `json:"releaseDate"`
	Time          string `json:"time"`
	EpisodeNo     string `json:"episode_no"`
}

func ExtractSchedule(date string, tzOffset int, baseURL string) ([]ScheduleItem, error) {
	c := utils.NewCollector()
	var items []ScheduleItem
	url := fmt.Sprintf("https://%s/ajax/schedule/list?tzOffset=%d&date=%s", baseURL, tzOffset, date)
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
		doc.Find("li").Each(func(_ int, sel *goquery.Selection) {
			var it ScheduleItem
			href := sel.Find("a").AttrOr("href", "")
			href = strings.Split(href, "?")[0]
			it.ID = strings.TrimPrefix(href, "/")
			it.DataID = lastSegment(it.ID)
			it.Title = strings.TrimSpace(sel.Find(".film-name").Text())
			it.JapaneseTitle = strings.TrimSpace(sel.Find(".film-name").AttrOr("data-jname", ""))
			it.ReleaseDate = date
			it.Time = strings.TrimSpace(sel.Find(".time").Text())
			it.EpisodeNo = lastSegment(strings.TrimSpace(sel.Find(".btn-play").Text()))
			items = append(items, it)
		})
	})
	var errVisit error
	c.OnError(func(_ *colly.Response, err error) { errVisit = err })
	_ = c.Visit(url)
	c.Wait()
	if errVisit != nil {
		return nil, errVisit
	}
	return items, nil
}

// ExtractNextEpisodeSchedule extracts next episode schedule information.
// It loads the watch page and reads `.schedule-alert > .alert.small > span:last`
// data-value attribute, which contains the next episode timestamp.
func ExtractNextEpisodeSchedule(id string, baseURL string) (string, error) {
	c := utils.NewCollector()
	var value string
	url := fmt.Sprintf("https://%s/watch/%s", baseURL, id)
	c.OnHTML(".schedule-alert > .alert.small > span:last-child", func(e *colly.HTMLElement) {
		if v := e.Attr("data-value"); strings.TrimSpace(v) != "" {
			value = strings.TrimSpace(v)
		}
	})
	var errVisit error
	c.OnError(func(_ *colly.Response, err error) { errVisit = err })
	_ = c.Visit(url)
	c.Wait()
	if errVisit != nil {
		return "", errVisit
	}
	return value, nil
}
