package parsers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/mnuddindev/jutsu-api/pkg/httpclient"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

type serverResponse struct {
	HTML string `json:"html"`
}

func FetchServerDataV1(episodeID string) ([]ServerData, error) {
	baseURL := utils.GetV1BaseURL()
	url := fmt.Sprintf("%s/ajax/v2/episode/servers?episodeId=%s", baseURL, episodeID)
	return fetchAndFilterServers(url, func(name string) bool {
		normalized := strings.TrimSpace(strings.ToUpper(name))
		return normalized == "HD-1" || normalized == "HD-2"
	})
}

func FetchServerDataV2(episodeID string) ([]ServerData, error) {
	baseURL := utils.GetV2BaseURL()
	url := fmt.Sprintf("%s/ajax/episode/servers?episodeId=%s", baseURL, episodeID)
	return fetchAndFilterServers(url, func(name string) bool {
		return strings.EqualFold(strings.TrimSpace(name), "Vidcloud")
	})
}

func fetchAndFilterServers(url string, predicate func(string) bool) ([]ServerData, error) {
	raw, err := httpclient.Get(url)
	if err != nil {
		return nil, err
	}
	var payload serverResponse
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		return nil, err
	}
	if strings.TrimSpace(payload.HTML) == "" {
		return []ServerData{}, nil
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(payload.HTML))
	if err != nil {
		return nil, err
	}

	var servers []ServerData
	doc.Find("div.ps_-block > div.ps__-list > div.server-item").Each(func(_ int, sel *goquery.Selection) {
		name := strings.TrimSpace(sel.Find("a.btn").Text())
		if !predicate(name) {
			return
		}
		server := ServerData{
			Name: name,
			ID:   sel.AttrOr("data-id", ""),
			Type: sel.AttrOr("data-type", ""),
		}
		servers = append(servers, server)
	})
	return servers, nil
}
