package scrape

import (
	"regexp"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

func ExtractToken(url, v1BaseURL string) (string, error) {
	c := utils.NewCollector()
	results := map[string]string{}
	strAssign := regexp.MustCompile(`window\.(\w+)\s*=\s*["']([\w-]+)["']`)
	commentRe := regexp.MustCompile(`<!--\s*_is_th:([\w-]+)\s*-->`)

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Referer", "https://"+v1BaseURL+"/")
	})

	c.OnHTML(`meta[name="_gg_fb"]`, func(e *colly.HTMLElement) {
		if v := e.Attr("content"); v != "" {
			results["meta"] = v
		}
	})
	c.OnHTML(`[data-dpi]`, func(e *colly.HTMLElement) {
		if v := e.Attr("data-dpi"); v != "" {
			results["dataDpi"] = v
		}
	})
	c.OnHTML("script[nonce]", func(e *colly.HTMLElement) {
		if strings.Contains(e.Text, "empty nonce script") {
			if n := e.Attr("nonce"); n != "" {
				results["nonce"] = n
			}
		}
	})
	c.OnResponse(func(r *colly.Response) {
		html := string(r.Body)
		for _, m := range strAssign.FindAllStringSubmatch(html, -1) {
			if len(m) >= 3 {
				results["window."+m[1]] = m[2]
			}
		}
		if m := commentRe.FindStringSubmatch(html); len(m) >= 2 {
			results["commentToken"] = m[1]
		}
	})
	var visitErr error
	c.OnError(func(_ *colly.Response, err error) { visitErr = err })
	_ = c.Visit(url)
	c.Wait()
	if visitErr != nil {
		return "", visitErr
	}
	for _, v := range results {
		if v != "" {
			return v, nil
		}
	}
	return "", nil
}
