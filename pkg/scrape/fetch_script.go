package scrape

import (
	"github.com/mnuddindev/jutsu-api/pkg/httpclient"
)

func FetchScript(url string) (string, error) {
	return httpclient.Get(url)
}
