package parsers

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/mnuddindev/jutsu-api/pkg/httpclient"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

var createServerReferenceRegex = regexp.MustCompile(`\(0,\w+\.createServerReference\)\("([a-f0-9]+)",\w+\.callServer,void 0,\w+\.findSourceMapURL,"(getSources|getEpisodes)"\)`)

type AniplayExtractor struct {
	baseURL string

	mu     sync.Mutex
	keys   map[string]string
	keysTS time.Time
}

func NewAniplayExtractor(baseURL string) *AniplayExtractor {
	if strings.TrimSpace(baseURL) == "" {
		baseURL = utils.GetV3BaseURL()
	}
	return &AniplayExtractor{
		baseURL: baseURL,
	}
}

func (a *AniplayExtractor) isCacheValid() bool {
	return a.keys != nil && time.Since(a.keysTS) < time.Hour
}

func (a *AniplayExtractor) fetchHTML(url string) (string, error) {
	return httpclient.GetWithHeaders(url, nil)
}

func (a *AniplayExtractor) fetchStaticJSURL() (string, error) {
	html, err := a.fetchHTML(fmt.Sprintf("%s/anime/watch/1", a.baseURL))
	if err != nil {
		return "", err
	}
	const prefix = "/_next/static/chunks/app/(user)/(media)/"
	idx := strings.Index(html, prefix)
	if idx == -1 {
		return "", fmt.Errorf("static chunk path not found in HTML")
	}
	slugStart := idx + len(prefix)
	slugEnd := strings.Index(html[slugStart:], "\"")
	if slugEnd == -1 {
		return "", fmt.Errorf("static chunk path terminator not found")
	}
	jsSlug := html[slugStart : slugStart+slugEnd]
	return a.baseURL + prefix + jsSlug, nil
}

func (a *AniplayExtractor) extractKeys() (map[string]string, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.isCacheValid() {
		return a.keys, nil
	}

	scriptURL, err := a.fetchStaticJSURL()
	if err != nil {
		return nil, err
	}
	script, err := a.fetchHTML(scriptURL)
	if err != nil {
		return nil, err
	}

	keys := map[string]string{
		"baseUrl": a.baseURL,
	}
	matches := createServerReferenceRegex.FindAllStringSubmatch(script, -1)
	for _, match := range matches {
		if len(match) != 3 {
			continue
		}
		hash := match[1]
		fn := match[2]
		keys[fn] = hash
	}
	if keys["getSources"] == "" || keys["getEpisodes"] == "" {
		return nil, fmt.Errorf("could not extract all required keys")
	}

	a.keys = keys
	a.keysTS = time.Now()
	return keys, nil
}

func (a *AniplayExtractor) GetNextAction() (map[string]string, error) {
	keys, err := a.extractKeys()
	if err != nil {
		return nil, err
	}
	return map[string]string{
		"watch": keys["getSources"],
		"info":  keys["getEpisodes"],
	}, nil
}

func (a *AniplayExtractor) FetchEpisode(animeID string, ep string, host string, typ string) (map[string]interface{}, error) {
	if strings.TrimSpace(host) == "" {
		host = "hika"
	}
	if strings.TrimSpace(typ) == "" {
		typ = "sub"
	}
	nextAction, err := a.GetNextAction()
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/anime/watch/%s?host=%s&ep=%s&type=%s", a.baseURL, animeID, host, ep, typ)
	payload := []string{
		animeID,
		host,
		fmt.Sprintf("%s/%s", animeID, ep),
		ep,
		typ,
	}
	headers := map[string]string{
		"Next-Action": nextAction["watch"],
	}
	resp, err := httpclient.PostJSON(url, payload, headers)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	idx := strings.Index(resp, "1:")
	if idx == -1 {
		return nil, fmt.Errorf("unexpected response format from aniplay")
	}
	jsonStr := strings.TrimSpace(resp[idx+2:])
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return nil, fmt.Errorf("failed to parse episode response: %w", err)
	}
	return data, nil
}
