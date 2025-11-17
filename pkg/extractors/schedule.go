package extractors

import (
	"bytes"
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
		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(r.Body))
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
