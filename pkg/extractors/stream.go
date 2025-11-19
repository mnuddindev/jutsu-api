package extractors

import (
	"encoding/json"
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

func ExtractStreamingInfo(fullID, name, typ string, fallback bool, baseURL string) (StreamingInfo, error) {
	// Extract just the episode ID from the full format (anime-id?ep=episode-id)
	episodeID := extractEpisodeID(fullID)
	if episodeID == "" {
		return StreamingInfo{}, fmt.Errorf("invalid ID format, expected 'anime-id?ep=episode-id'")
	}

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

	// Match server by both name AND type
	selected := findServerMatchByType(servers, requestedServer, streamType)
	// If no match, try with "raw" type
	if selected == nil {
		selected = findServerMatchByType(servers, requestedServer, "raw")
	}

	// If no server match found, return servers only
	if selected == nil {
		if fallback {
			return StreamingInfo{Servers: servers}, fmt.Errorf("no matching server found for name: %s, type: %s", requestedServer, streamType)
		}
		return StreamingInfo{Servers: servers}, nil
	}

	serverID := strings.TrimSpace(selected.DataID)
	serverName := selected.ServerName
	if strings.TrimSpace(serverName) == "" {
		serverName = requestedServer
	}

	// Extract episode ID for fallback case (needs just episode ID, not full format)
	episodeIDOnly := extractEpisodeID(fullID)
	if episodeIDOnly == "" {
		return StreamingInfo{Servers: servers}, fmt.Errorf("failed to extract episode ID from: %s", fullID)
	}

	// Use full ID format for decryption (matching Node.js behavior)
	// For fallback, we need to pass just the episode ID, not the full format
	stream, err := resolveStreamingLink(fullID, episodeIDOnly, serverID, serverName, streamType, fallback)
	if err != nil {
		return StreamingInfo{Servers: servers}, err
	}
	return StreamingInfo{StreamingLink: stream, Servers: servers}, nil
}

// extractEpisodeID extracts the episode ID from format "anime-id?ep=episode-id"
func extractEpisodeID(fullID string) string {
	parts := strings.Split(fullID, "?ep=")
	if len(parts) < 2 {
		return ""
	}
	return strings.TrimSpace(parts[1])
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

// findServerMatchByType matches server by both name AND type (matching Node.js behavior)
func findServerMatchByType(servers []ServerItem, name, typ string) *ServerItem {
	normalizedName := strings.TrimSpace(strings.ToLower(name))
	normalizedType := strings.TrimSpace(strings.ToLower(typ))
	for i := range servers {
		serverName := strings.ToLower(strings.TrimSpace(servers[i].ServerName))
		serverType := strings.ToLower(strings.TrimSpace(servers[i].Type))
		if serverName == normalizedName && serverType == normalizedType {
			return &servers[i]
		}
	}
	return nil
}

func resolveStreamingLink(fullEpisodeID, episodeIDOnly, serverID, serverName, typ string, useFallback bool) (parsers.DecryptedSources, error) {
	// fullEpisodeID is the full format (anime-id?ep=episode-id) - used for primary path
	// episodeIDOnly is just the episode ID (e.g., "107257") - used for fallback path
	// serverID is the data_id from the matched server
	// For fallback, serverID is not required (decryptFallback doesn't use it)
	if !useFallback && strings.TrimSpace(serverID) == "" {
		return parsers.DecryptedSources{}, fmt.Errorf("missing server id for %s", serverName)
	}

	// Use v1 decryptor (matching Node.js behavior exactly)
	// For fallback: epID should be just the episode ID (e.g., "107257"), serverID can be empty
	// For primary: epID is not used, only the serverID (data_id) is used
	// Node.js calls: decryptSources_v1(id, requestedServer[0].data_id, name, type, fallback)
	// where id is the full format (anime-id?ep=episode-id) for primary, but just episode ID for fallback
	if useFallback {
		return decryptSourcesV1Fn(episodeIDOnly, serverID, serverName, typ, useFallback)
	}
	return decryptSourcesV1Fn(fullEpisodeID, serverID, serverName, typ, useFallback)
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

// FindServerMatchByTypeForTest exposes findServerMatchByType for tests.
func FindServerMatchByTypeForTest(servers []ServerItem, name, typ string) *ServerItem {
	return findServerMatchByType(servers, name, typ)
}

// ResolveStreamingLinkForTest exposes resolveStreamingLink for unit tests.
func ResolveStreamingLinkForTest(fullEpisodeID, episodeIDOnly, serverID, serverName, typ string, fallback bool) (parsers.DecryptedSources, error) {
	return resolveStreamingLink(fullEpisodeID, episodeIDOnly, serverID, serverName, typ, fallback)
}
