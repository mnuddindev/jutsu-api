package parsers

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/mnuddindev/jutsu-api/pkg/httpclient"
	"github.com/mnuddindev/jutsu-api/pkg/scrape"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

var (
	hostRegex   = regexp.MustCompile(`^(https?:\/\/(?:www\.)?[^\/\?]+)`)
	sourceIDReg = regexp.MustCompile(`/([^\/\?]+)\?`)
)

type embedV2Response struct {
	Encrypted bool            `json:"encrypted"`
	Sources   json.RawMessage `json:"sources"`
	Tracks    json.RawMessage `json:"tracks"`
	Intro     json.RawMessage `json:"intro"`
	Outro     json.RawMessage `json:"outro"`
}

// DecryptSourcesV2 mirrors decrypt_v2.decryptor.js and is kept as a fallback
// in case Rapid/Megacloud reverts to the legacy architecture.
func DecryptSourcesV2(id, name, typ string) (DecryptedSources, error) {
	baseURL := utils.GetV2BaseURL()
	srcURL := fmt.Sprintf("%s/ajax/episode/sources?id=%s", baseURL, id)
	raw, err := httpclient.Get(srcURL)
	if err != nil {
		return DecryptedSources{}, err
	}
	var payload episodeSourceData
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		return DecryptedSources{}, err
	}
	if strings.TrimSpace(payload.Link) == "" {
		return DecryptedSources{}, errors.New("decrypt_v2: empty iframe link")
	}

	host, sourceID, err := parseV2Embed(payload.Link)
	if err != nil {
		return DecryptedSources{}, err
	}

	script, err := scrape.FetchScript(playerScriptURL)
	if err != nil {
		return DecryptedSources{}, err
	}
	keyPairs := utils.GetKeys(script)
	if len(keyPairs) == 0 {
		return DecryptedSources{}, errors.New("decrypt_v2: unable to extract key pairs")
	}

	embedURL := fmt.Sprintf("%s/ajax/embed-6-v2/getSources?id=%s", host, sourceID)
	embedRaw, err := httpclient.Get(embedURL)
	if err != nil {
		return DecryptedSources{}, err
	}
	var embedResp embedV2Response
	if err := json.Unmarshal([]byte(embedRaw), &embedResp); err != nil {
		return DecryptedSources{}, err
	}

	link, err := resolveV2Sources(embedResp, keyPairs)
	if err != nil {
		return DecryptedSources{}, err
	}

	return DecryptedSources{
		ID:     id,
		Type:   typ,
		Link:   link,
		Tracks: embedResp.Tracks,
		Intro:  embedResp.Intro,
		Outro:  embedResp.Outro,
		Server: name,
	}, nil
}

func parseV2Embed(embedLink string) (string, string, error) {
	mHost := hostRegex.FindStringSubmatch(embedLink)
	if len(mHost) < 2 {
		return "", "", errors.New("decrypt_v2: unable to determine host")
	}
	mID := sourceIDReg.FindStringSubmatch(embedLink)
	if len(mID) < 2 {
		return "", "", errors.New("decrypt_v2: unable to find source id")
	}
	return mHost[1], mID[1], nil
}

func resolveV2Sources(resp embedV2Response, pairs [][]int) (StreamLink, error) {
	if resp.Encrypted {
		var cipherText string
		if err := json.Unmarshal(resp.Sources, &cipherText); err != nil {
			return StreamLink{}, err
		}
		key, cipher, err := extractCipherAndKey(cipherText, pairs)
		if err != nil {
			return StreamLink{}, err
		}
		plain, err := aesDecrypt(cipher, key)
		if err != nil {
			return StreamLink{}, err
		}
		var assets []StreamLink
		if err := json.Unmarshal(plain, &assets); err != nil {
			return StreamLink{}, err
		}
		if len(assets) == 0 {
			return StreamLink{}, errors.New("decrypt_v2: decrypted payload empty")
		}
		return StreamLink{File: assets[0].File, Type: "hls"}, nil
	}

	var link StreamLink
	if err := json.Unmarshal(resp.Sources, &link); err == nil && link.File != "" {
		if strings.TrimSpace(link.Type) == "" {
			link.Type = "hls"
		}
		return link, nil
	}

	var arr []StreamLink
	if err := json.Unmarshal(resp.Sources, &arr); err == nil && len(arr) > 0 {
		link = arr[0]
		if strings.TrimSpace(link.Type) == "" {
			link.Type = "hls"
		}
		return link, nil
	}
	return StreamLink{}, errors.New("decrypt_v2: unable to interpret sources")
}

func extractCipherAndKey(raw string, pairs [][]int) (string, string, error) {
	runes := []rune(raw)
	var keyBuilder strings.Builder
	current := 0
	for _, pair := range pairs {
		if len(pair) != 2 {
			continue
		}
		start := pair[0] + current
		end := start + pair[1]
		if start < 0 || end > len(runes) {
			return "", "", errors.New("decrypt_v2: key pair exceeds bounds")
		}
		for i := start; i < end; i++ {
			keyBuilder.WriteRune(runes[i])
			runes[i] = 0
		}
		current += pair[1]
	}
	var cipherBuilder strings.Builder
	cipherBuilder.Grow(len(runes))
	for _, r := range runes {
		if r != 0 {
			cipherBuilder.WriteRune(r)
		}
	}
	return keyBuilder.String(), cipherBuilder.String(), nil
}
