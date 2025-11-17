package parsers

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "embed"

	"github.com/mnuddindev/jutsu-api/pkg/httpclient"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

//go:embed assets/wasm_bridge.js
var wasmBridgeScript string

const (
	megacloudUserAgent    = "Mozilla/5.0 (X11; Linux x86_64; rv:133.0) Gecko/20100101 Firefox/133.0"
	megacloudBrowserValue = 1676800512
	wasmBridgeName        = "wasm-bridge.js"
)

var (
	metaRegex         = regexp.MustCompile(`name="j_crt"\s+content="([A-Za-z0-9]+)"`)
	wasmBridgeOnce    sync.Once
	wasmBridgePath    string
	wasmBridgeErr     error
	defaultHTTPClient = &http.Client{Timeout: 15 * time.Second}
)

// DecryptMegacloud resolves the Megacloud iframe and returns the decrypted sources.
func DecryptMegacloud(episodeID, serverName, typ string) (DecryptedSources, error) {
	link, err := fetchMegacloudEmbedLink(episodeID)
	if err != nil {
		return DecryptedSources{}, err
	}
	payload, err := decryptMegacloudSource(link)
	if err != nil {
		return DecryptedSources{}, err
	}

	stream, err := extractMegacloudStream(payload.Sources)
	if err != nil {
		return DecryptedSources{}, err
	}

	return DecryptedSources{
		ID:     episodeID,
		Type:   typ,
		Link:   stream,
		Tracks: payload.Tracks,
		Intro:  payload.Intro,
		Outro:  payload.Outro,
		Iframe: link,
		Server: serverName,
	}, nil
}

type episodeSourceData struct {
	Link string `json:"link"`
}

type megacloudSourcesResponse struct {
	Sources json.RawMessage `json:"sources"`
	Tracks  json.RawMessage `json:"tracks"`
	Intro   json.RawMessage `json:"intro"`
	Outro   json.RawMessage `json:"outro"`
	T       int             `json:"t"`
	K       []int           `json:"k"`
}

type wasmBridgeResult struct {
	PID      string `json:"pid"`
	KVersion string `json:"kversion"`
	Kid      string `json:"kid"`
	Bytes    []byte `json:"bytes"`
	Error    string `json:"error"`
}

func fetchMegacloudEmbedLink(episodeID string) (string, error) {
	url := fmt.Sprintf("%s/ajax/episode/sources?id=%s", utils.GetV4BaseURL(), episodeID)
	body, err := httpclient.Get(url)
	if err != nil {
		return "", err
	}
	var resp episodeSourceData
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		return "", err
	}
	if strings.TrimSpace(resp.Link) == "" {
		return "", errors.New("megacloud: empty iframe link")
	}
	return resp.Link, nil
}

func decryptMegacloudSource(embedURL string) (*megacloudSourcesResponse, error) {
	parsed, err := url.Parse(embedURL)
	if err != nil {
		return nil, err
	}
	baseURL := fmt.Sprintf("%s://%s", parsed.Scheme, parsed.Host)
	referrer := baseURL
	if strings.Contains(strings.ToLower(parsed.Host), "mega") {
		referrer = utils.GetV4BaseURL()
	}

	metaContent, err := fetchMetaContent(embedURL, referrer)
	if err != nil {
		return nil, err
	}

	wasmResult, err := runWasmBridge(baseURL, metaContent)
	if err != nil {
		return nil, err
	}

	getSourcesURL, err := buildGetSourcesURL(embedURL, baseURL, wasmResult)
	if err != nil {
		return nil, err
	}
	payload, err := requestSources(getSourcesURL, embedURL)
	if err != nil {
		return nil, err
	}

	keyInt, err := strconv.Atoi(wasmResult.KVersion)
	if err != nil {
		return nil, fmt.Errorf("invalid kversion: %w", err)
	}
	key := intToRGBA(keyInt)
	decryptionKey, err := deriveMegacloudKey(payload, wasmResult.Bytes, key)
	if err != nil {
		return nil, err
	}

	decryptedSources, err := decryptSourcesPayload(payload.Sources, decryptionKey)
	if err != nil {
		return nil, err
	}
	payload.Sources = decryptedSources
	return payload, nil
}

func fetchMetaContent(embedURL, referrer string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, embedURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", megacloudUserAgent)
	req.Header.Set("Referer", referrer)

	resp, err := defaultHTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	match := metaRegex.FindStringSubmatch(string(body))
	if len(match) != 2 {
		return "", errors.New("megacloud: meta tag not found")
	}
	return match[1] + "==", nil
}

func runWasmBridge(baseURL, meta string) (*wasmBridgeResult, error) {
	scriptPath, err := ensureWasmBridge()
	if err != nil {
		return nil, err
	}
	nodePath, err := exec.LookPath("node")
	if err != nil {
		return nil, errors.New("node executable not found in PATH")
	}

	wasmURL := baseURL + "/images/loading.png?v=0.0.9"
	imageURL := baseURL + "/images/image.png?v=0.0.9"

	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, nodePath, scriptPath, baseURL, wasmURL, imageURL, meta)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("wasm bridge failed: %w (%s)", err, strings.TrimSpace(stderr.String()))
	}

	var result wasmBridgeResult
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		return nil, fmt.Errorf("invalid bridge output: %w", err)
	}
	if result.Error != "" {
		return nil, fmt.Errorf("wasm bridge error: %s", result.Error)
	}
	if result.PID == "" || result.KVersion == "" || result.Kid == "" {
		return nil, errors.New("wasm bridge returned incomplete data")
	}
	return &result, nil
}

func ensureWasmBridge() (string, error) {
	wasmBridgeOnce.Do(func() {
		dir := filepath.Join(os.TempDir(), "jutsu-megacloud")
		if err := os.MkdirAll(dir, 0o755); err != nil {
			wasmBridgeErr = err
			return
		}
		path := filepath.Join(dir, wasmBridgeName)
		if err := os.WriteFile(path, []byte(wasmBridgeScript), 0o644); err != nil {
			wasmBridgeErr = err
			return
		}
		wasmBridgePath = path
	})
	return wasmBridgePath, wasmBridgeErr
}

func buildGetSourcesURL(embedURL, baseURL string, result *wasmBridgeResult) (string, error) {
	parts := strings.Split(embedURL, "/")
	if len(parts) < 5 {
		return "", errors.New("megacloud: unexpected embed URL format")
	}
	hostLower := strings.ToLower(parts[2])
	query := fmt.Sprintf("id=%s&v=%s&h=%s&b=%d", url.QueryEscape(result.PID), url.QueryEscape(result.KVersion), url.QueryEscape(result.Kid), megacloudBrowserValue)
	if strings.Contains(hostLower, "mega") {
		return fmt.Sprintf("%s/%s/ajax/%s/getSources?%s", baseURL, parts[3], parts[4], query), nil
	}
	return fmt.Sprintf("%s/ajax/%s/getSources?%s", baseURL, parts[3], query), nil
}

func requestSources(getSourcesURL, embedURL string) (*megacloudSourcesResponse, error) {
	req, err := http.NewRequest(http.MethodGet, getSourcesURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", megacloudUserAgent)
	req.Header.Set("Referer", embedURL+"&autoPlay=1&oa=0&asi=1")
	req.Header.Set("Accept-Language", "en,bn;q=0.9,en-US;q=0.8")
	req.Header.Set("sec-ch-ua", `"Google Chrome";v="131", "Chromium";v="131", "Not_A Brand";v="24"`)
	req.Header.Set("sec-ch-ua-mobile", "?1")
	req.Header.Set("sec-ch-ua-platform", `"Android"`)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	resp, err := defaultHTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var payload megacloudSourcesResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}
	return &payload, nil
}

func deriveMegacloudKey(payload *megacloudSourcesResponse, wasmBytes []byte, key []int) ([]byte, error) {
	if payload == nil {
		return nil, errors.New("nil payload")
	}
	if payload.T != 0 {
		buf := make([]byte, len(wasmBytes))
		copy(buf, wasmBytes)
		xorBytes(buf, key)
		return buf, nil
	}
	if len(payload.K) == 0 {
		return nil, errors.New("payload missing key material")
	}
	buf := make([]byte, len(payload.K))
	for i, v := range payload.K {
		buf[i] = byte(v)
	}
	xorBytes(buf, key)
	return buf, nil
}

func decryptSourcesPayload(raw json.RawMessage, key []byte) (json.RawMessage, error) {
	if len(raw) == 0 {
		return nil, errors.New("empty sources payload")
	}
	var cipherText string
	if err := json.Unmarshal(raw, &cipherText); err != nil {
		// already decrypted JSON
		return raw, nil
	}
	decodedKey := base64.StdEncoding.EncodeToString(key)
	plain, err := aesDecrypt(cipherText, decodedKey)
	if err != nil {
		return nil, err
	}
	trimmed := bytes.TrimSpace(plain)
	return json.RawMessage(trimmed), nil
}

func extractMegacloudStream(raw json.RawMessage) (StreamLink, error) {
	if len(raw) == 0 {
		return StreamLink{}, errors.New("no sources available")
	}
	var arr []struct {
		File string `json:"file"`
		Type string `json:"type"`
	}
	if err := json.Unmarshal(raw, &arr); err == nil && len(arr) > 0 && strings.TrimSpace(arr[0].File) != "" {
		return StreamLink{File: arr[0].File, Type: fallbackType(arr[0].Type)}, nil
	}
	var single struct {
		File string `json:"file"`
		Type string `json:"type"`
	}
	if err := json.Unmarshal(raw, &single); err == nil && strings.TrimSpace(single.File) != "" {
		return StreamLink{File: single.File, Type: fallbackType(single.Type)}, nil
	}
	var file string
	if err := json.Unmarshal(raw, &file); err == nil && strings.TrimSpace(file) != "" {
		return StreamLink{File: file, Type: "hls"}, nil
	}
	return StreamLink{}, errors.New("unable to determine stream source")
}

func fallbackType(t string) string {
	if strings.TrimSpace(t) == "" {
		return "hls"
	}
	return t
}

func intToRGBA(a int) []int {
	return []int{
		(a & 0xFF000000) >> 24,
		(a & 0x00FF0000) >> 16,
		(a & 0x0000FF00) >> 8,
		a & 0x000000FF,
	}
}

func xorBytes(data []byte, key []int) {
	if len(key) == 0 {
		return
	}
	for i := 0; i < len(data); i++ {
		data[i] ^= byte(key[i%len(key)])
	}
}

// ioReadAll mirrors io.ReadAll but avoids importing deprecated helper.
