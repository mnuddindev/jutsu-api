package extractors

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/mnuddindev/jutsu-api/pkg/parsers"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

var (
	decryptMegacloudFn = parsers.DecryptMegacloud
	decryptSourcesV1Fn = parsers.DecryptSourcesV1
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
	StreamingLink parsers.DecryptedSources `json:"streamingLink"`
	Servers       []ServerItem             `json:"servers"`
}

func ExtractStreamingInfo(episodeID, name, typ string, fallback bool, baseURL string) (StreamingInfo, error) {
	servers, err := ExtractServers(episodeID, baseURL)
	if err != nil {
		return StreamingInfo{}, err
	}

	requestedServer := strings.TrimSpace(name)
	if requestedServer == "" {
		return StreamingInfo{Servers: servers}, nil
	}

	streamType := strings.TrimSpace(typ)
	if streamType == "" {
		streamType = "sub"
	}

	selected := findServerMatch(servers, requestedServer)
	if selected == nil && !fallback {
		return StreamingInfo{Servers: servers}, nil
	}

	serverID := ""
	serverName := requestedServer
	if selected != nil {
		serverID = strings.TrimSpace(selected.DataID)
		if trimmed := strings.TrimSpace(selected.ServerName); trimmed != "" {
			serverName = trimmed
		}
	}

	stream, err := resolveStreamingLink(episodeID, serverID, serverName, streamType, fallback)
	if err != nil {
		return StreamingInfo{Servers: servers}, nil
	}
	return StreamingInfo{StreamingLink: stream, Servers: servers}, nil
}

func findServerMatch(servers []ServerItem, target string) *ServerItem {
	normalized := strings.TrimSpace(strings.ToLower(target))
	for i := range servers {
		switch {
		case normalized != "" && strings.ToLower(strings.TrimSpace(servers[i].ServerName)) == normalized:
			return &servers[i]
		case normalized != "" && strings.ToLower(strings.TrimSpace(servers[i].DataID)) == normalized:
			return &servers[i]
		case normalized != "" && strings.ToLower(strings.TrimSpace(servers[i].ServerID)) == normalized:
			return &servers[i]
		}
	}
	return nil
}

func resolveStreamingLink(episodeID, serverID, serverName, typ string, useFallback bool) (parsers.DecryptedSources, error) {
	if useFallback {
		fallbackID := serverID
		if strings.TrimSpace(fallbackID) == "" {
			fallbackID = episodeID
		}
		return decryptSourcesV1Fn(episodeID, fallbackID, serverName, typ, true)
	}

	if strings.TrimSpace(serverID) == "" {
		return parsers.DecryptedSources{}, fmt.Errorf("missing server id for %s", serverName)
	}

	stream, err := decryptMegacloudFn(serverID, serverName, typ)
	if err == nil {
		return stream, nil
	}

	legacy, legacyErr := decryptSourcesV1Fn(episodeID, serverID, serverName, typ, false)
	if legacyErr == nil {
		return legacy, nil
	}
	return parsers.DecryptedSources{}, err
}

// SetStreamDecryptorsForTest overrides the decryptor functions and returns a restore
// callback. Intended for use in unit tests located outside this package.
func SetStreamDecryptorsForTest(
	megacloud func(string, string, string) (parsers.DecryptedSources, error),
	legacy func(string, string, string, string, bool) (parsers.DecryptedSources, error),
) func() {
	prevMega := decryptMegacloudFn
	prevLegacy := decryptSourcesV1Fn
	if megacloud != nil {
		decryptMegacloudFn = megacloud
	}
	if legacy != nil {
		decryptSourcesV1Fn = legacy
	}
	return func() {
		decryptMegacloudFn = prevMega
		decryptSourcesV1Fn = prevLegacy
	}
}

// FindServerMatchForTest exposes findServerMatch for black-box tests.
func FindServerMatchForTest(servers []ServerItem, target string) *ServerItem {
	return findServerMatch(servers, target)
}

// ResolveStreamingLinkForTest exposes resolveStreamingLink for unit tests.
func ResolveStreamingLinkForTest(episodeID, serverID, serverName, typ string, fallback bool) (parsers.DecryptedSources, error) {
	return resolveStreamingLink(episodeID, serverID, serverName, typ, fallback)
}
