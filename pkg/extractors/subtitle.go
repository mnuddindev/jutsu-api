package extractors

import (
	"fmt"
	"strings"

	"github.com/mnuddindev/jutsu-api/pkg/httpclient"
)

type SubtitleResult struct {
	Subtitles interface{} `json:"subtitles"`
	Intro     interface{} `json:"intro"`
	Outro     interface{} `json:"outro"`
}

func ExtractSubtitle(id, baseURL, provider string) (SubtitleResult, error) {
	// First: get sources link from v1
	resp, err := httpclient.Get(fmt.Sprintf("https://%s/ajax/v2/episode/sources/?id=%s", baseURL, id))
	if err != nil {
		return SubtitleResult{}, err
	}
	// naive parse to find link (JSON); we could unmarshal but keep light
	link := between(resp, "\"link\":\"", "\"")
	if link == "" {
		return SubtitleResult{}, nil
	}
	last := lastSegment(link)
	last = strings.ReplaceAll(last, "?k=", "")
	src, err := httpclient.Get(fmt.Sprintf("%s/embed-2/ajax/e-1/getSources?id=%s", provider, last))
	if err != nil {
		return SubtitleResult{}, err
	}
	// Return raw JSON fragments (for now)
	return SubtitleResult{Subtitles: src, Intro: nil, Outro: nil}, nil
}

func between(s, a, b string) string {
	i := strings.Index(s, a)
	if i < 0 {
		return ""
	}
	s = s[i+len(a):]
	j := strings.Index(s, b)
	if j < 0 {
		return ""
	}
	return s[:j]
}
