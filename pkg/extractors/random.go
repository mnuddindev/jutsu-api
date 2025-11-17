package extractors

import (
	"fmt"
	"net/http"
)

func ExtractRandomID(baseURL string) (string, error) {
	client := &http.Client{CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse }}
	resp, err := client.Get(fmt.Sprintf("https://%s/random", baseURL))
	if err != nil {
		return "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	loc := resp.Header.Get("Location")
	if loc == "" {
		return "", nil
	}
	return lastSegment(loc), nil
}

// ExtractRandom fetches a random anime ID and returns its full info.
func ExtractRandom(baseURL string) (AnimeInfo, error) {
	id, err := ExtractRandomID(baseURL)
	if err != nil {
		return AnimeInfo{}, err
	}
	if id == "" {
		return AnimeInfo{}, fmt.Errorf("random id is empty")
	}
	return ExtractAnimeInfo(id, baseURL)
}
