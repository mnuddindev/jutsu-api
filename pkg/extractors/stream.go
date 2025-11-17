package extractors

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

type ServerItem struct {
	Type       string `json:"type"`
	DataID     string `json:"data_id"`
	ServerID   string `json:"server_id"`
	ServerName string `json:"serverName"`
}

func ExtractServers(episodeID string, baseURL string) ([]ServerItem, error) {
	c := utils.NewCollector()
	var servers []ServerItem
	url := fmt.Sprintf("https://%s/ajax/v2/episode/servers?episodeId=%s", baseURL, episodeID)
	c.OnResponse(func(r *colly.Response) {
		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(r.Body))
		if err != nil {
			return
		}
		doc.Find(".server-item").Each(func(_ int, sel *goquery.Selection) {
			var s ServerItem
			s.DataID = sel.AttrOr("data-id", "")
			s.ServerID = sel.AttrOr("data-server-id", "")
			s.Type = sel.AttrOr("data-type", "")
			s.ServerName = strings.TrimSpace(sel.Find("a").Text())
			servers = append(servers, s)
		})
	})
	var errVisit error
	c.OnError(func(_ *colly.Response, err error) { errVisit = err })
	_ = c.Visit(url)
	c.Wait()
	if errVisit != nil {
		return nil, errVisit
	}
	return servers, nil
}

type StreamingInfo struct {
	StreamingLink []interface{} `json:"streamingLink"`
	Servers       []ServerItem  `json:"servers"`
}

func ExtractStreamingInfo(id, name, typ string, fallback bool, baseURL string) (StreamingInfo, error) {
	// Placeholder: decrypt not implemented, return servers only
	episodeID := id
	if i := len(id) - 1; i >= 0 {
		episodeID = id
	}
	servers, err := ExtractServers(episodeID, baseURL)
	if err != nil {
		return StreamingInfo{StreamingLink: []interface{}{}, Servers: []ServerItem{}}, nil
	}
	return StreamingInfo{StreamingLink: []interface{}{}, Servers: servers}, nil
}
