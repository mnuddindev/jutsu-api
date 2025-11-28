package parsers

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/mnuddindev/jutsu-api/pkg/httpclient"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

func DecryptSourcesV1(epID string, id string, name string, typ string, fallback bool) (DecryptedSources, error) {
	if fallback {
		return decryptFallback(epID, id, name, typ)
	}
	return decryptPrimary(id, name, typ)
}

func decryptFallback(epID, id, name, typ string) (DecryptedSources, error) {
	targetServer := utils.GetFallback2Host()
	switch strings.ToLower(name) {
	case "hd-1", "hd-3":
		targetServer = utils.GetFallback1Host()
	}
	baseURL := ensureHTTPS(targetServer)
	if typ == "" {
		typ = "sub"
	}

	iframeURL := fmt.Sprintf("%s/stream/s-2/%s/%s", baseURL, epID, typ)
	html, err := httpclient.GetWithHeaders(iframeURL, map[string]string{
		"Referer": baseURL + "/",
	})
	if err != nil {
		return DecryptedSources{}, err
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return DecryptedSources{}, err
	}
	dataID, exists := doc.Find("#megaplay-player").Attr("data-id")
	if !exists || strings.TrimSpace(dataID) == "" {
		return DecryptedSources{}, fmt.Errorf("failed to extract data-id from fallback iframe")
	}

	sourceURL := fmt.Sprintf("%s/stream/getSources?id=%s", baseURL, dataID)
	raw, err := httpclient.GetWithHeaders(sourceURL, map[string]string{
		"X-Requested-With": "XMLHttpRequest",
	})
	if err != nil {
		return DecryptedSources{}, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &data); err != nil {
		return DecryptedSources{}, err
	}

	linkFile := extractSourceFile(data)
	if strings.TrimSpace(linkFile) == "" {
		return DecryptedSources{}, fmt.Errorf("failed to extract source file from fallback response")
	}
	return DecryptedSources{
		ID:   id,
		Type: typ,
		Link: StreamLink{
			File: linkFile,
			Type: "hls",
		},
		Tracks: toRawMessage(data["tracks"]),
		Intro:  toRawMessage(data["intro"]),
		Outro:  toRawMessage(data["outro"]),
		Iframe: iframeURL,
		Server: name,
	}, nil
}

func decryptPrimary(id, name, typ string) (DecryptedSources, error) {
	baseURL := utils.GetV4BaseURL()
	sourceEndpoint := fmt.Sprintf("%s/ajax/episode/sources?id=%s", baseURL, id)
	// Add Referer header
	raw, err := httpclient.GetWithHeaders(sourceEndpoint, map[string]string{
		"Referer": baseURL + "/",
	})
	if err != nil {
		return DecryptedSources{}, err
	}
	var payload struct {
		Link string `json:"link"`
	}
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		return DecryptedSources{}, err
	}
	if strings.TrimSpace(payload.Link) == "" {
		return DecryptedSources{}, fmt.Errorf("missing link in sourcesData")
	}

	iframeURL, getSourcesURL, err := buildSourceURLs(payload.Link)
	if err != nil {
		return DecryptedSources{}, err
	}

	// Parse the base URL from getSourcesURL to set Referer
	parsedSourcesURL, err := url.Parse(getSourcesURL)
	if err != nil {
		return DecryptedSources{}, err
	}
	refererURL := fmt.Sprintf("%s://%s/", parsedSourcesURL.Scheme, parsedSourcesURL.Host)

	// Add Referer and X-Requested-With headers for getSources request
	rawSource, err := httpclient.GetWithHeaders(getSourcesURL, map[string]string{
		"Referer":          refererURL,
		"X-Requested-With": "XMLHttpRequest",
	})
	if err != nil {
		return DecryptedSources{}, err
	}
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(rawSource), &data); err != nil {
		return DecryptedSources{}, err
	}

	linkFile := extractSourceFile(data)
	if strings.TrimSpace(linkFile) == "" {
		return DecryptedSources{}, fmt.Errorf("failed to extract source file from primary response")
	}
	return DecryptedSources{
		ID:   id,
		Type: typ,
		Link: StreamLink{
			File: linkFile,
			Type: "hls",
		},
		Tracks: toRawMessage(data["tracks"]),
		Intro:  toRawMessage(data["intro"]),
		Outro:  toRawMessage(data["outro"]),
		Iframe: iframeURL,
		Server: name,
	}, nil
}

func buildSourceURLs(ajaxLink string) (iframe string, sources string, err error) {
	parsed, err := url.Parse(ajaxLink)
	if err != nil {
		return "", "", err
	}
	pathParts := strings.Split(strings.Trim(parsed.Path, "/"), "/")
	if len(pathParts) < 3 {
		return "", "", fmt.Errorf("could not extract base URL from link")
	}
	sourceID := pathParts[len(pathParts)-1]
	basePath := strings.Join(pathParts[:len(pathParts)-1], "/")
	baseURL := fmt.Sprintf("%s://%s/%s", parsed.Scheme, parsed.Host, basePath)
	iframe = fmt.Sprintf("%s/%s?k=1&autoPlay=0&oa=0&asi=1", baseURL, sourceID)
	sources = fmt.Sprintf("%s/getSources?id=%s", baseURL, sourceID)
	return iframe, sources, nil
}

func extractSourceFile(data map[string]interface{}) string {
	if data == nil {
		return ""
	}
	if sources, ok := data["sources"]; ok {
		switch s := sources.(type) {
		case map[string]interface{}:
			if file, ok := s["file"].(string); ok {
				return file
			}
		case []interface{}:
			if len(s) > 0 {
				if entry, ok := s[0].(map[string]interface{}); ok {
					if file, ok := entry["file"].(string); ok {
						return file
					}
				}
			}
		}
	}
	return ""
}

func toRawMessage(v interface{}) json.RawMessage {
	if v == nil {
		return nil
	}
	b, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	return json.RawMessage(b)
}

func ensureHTTPS(host string) string {
	if strings.HasPrefix(host, "http://") || strings.HasPrefix(host, "https://") {
		return host
	}
	return "https://" + host
}
