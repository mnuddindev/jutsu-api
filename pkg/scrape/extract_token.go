package scrape

import (
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/mnuddindev/jutsu-api/pkg/httpclient"
)

func ExtractToken(url, v1BaseURL string) (string, error) {
	html, err := httpclient.Get(url)
	if err != nil {
		return "", err
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", err
	}
	results := map[string]string{}
	if meta, ok := doc.Find(`meta[name="_gg_fb"]`).Attr("content"); ok {
		results["meta"] = meta
	}
	if dpi, ok := doc.Find(`[data-dpi]`).Attr("data-dpi"); ok {
		results["dataDpi"] = dpi
	}
	doc.Find("script[nonce]").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if strings.Contains(s.Text(), "empty nonce script") {
			if nonce, ok := s.Attr("nonce"); ok {
				results["nonce"] = nonce
				return false
			}
		}
		return true
	})
	strAssign := regexp.MustCompile(`window\.(\w+)\s*=\s*["']([\w-]+)["']`)
	for _, m := range strAssign.FindAllStringSubmatch(html, -1) {
		if len(m) >= 3 {
			results["window."+m[1]] = m[2]
		}
	}
	commentRe := regexp.MustCompile(`<!--\s*_is_th:([\w-]+)\s*-->`)
	if m := commentRe.FindStringSubmatch(html); len(m) >= 2 {
		results["commentToken"] = m[1]
	}
	for _, v := range results {
		if v != "" {
			return v, nil
		}
	}
	return "", nil
}
