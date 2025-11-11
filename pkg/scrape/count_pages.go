package scrape

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/mnuddindev/jutsu-api/pkg/httpclient"
)

func CountPages(url string) (int, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}
	for k, v := range httpclient.DEFAULT_HEADERS {
		req.Header.Set(k, v)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return 0, errors.New("non-200 response")
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return 0, err
	}
	lastHref, exists := doc.Find(".tab-content .pagination .page-item:last-child a").Attr("href")
	if !exists || strings.TrimSpace(lastHref) == "" {
		return 1, nil
	}
	parts := strings.Split(lastHref, "=")
	last := parts[len(parts)-1]
	n, _ := strconv.Atoi(last)
	if n <= 0 {
		return 1, nil
	}
	return n, nil
}
