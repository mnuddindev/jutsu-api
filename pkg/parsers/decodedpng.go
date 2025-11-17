package parsers

import (
	_ "embed"
	"errors"
	"strconv"
	"strings"
	"sync"
)

//go:embed assets/decodedpng.js
var decodedPNGSource string

var (
	decodedPNGOnce sync.Once
	decodedPNGData []byte
	decodedPNGErr  error
)

// GetDecodedPNGData returns the decoded Uint8ClampedArray as raw bytes.
// The data is parsed once from the embedded JS asset to avoid shipping
// large literals inside the Go source.
func GetDecodedPNGData() ([]byte, error) {
	decodedPNGOnce.Do(func() {
		start := strings.Index(decodedPNGSource, "[")
		end := strings.LastIndex(decodedPNGSource, "]")
		if start == -1 || end == -1 || end <= start {
			decodedPNGErr = ErrDecodedPNGFormat
			return
		}
		content := decodedPNGSource[start+1 : end]
		fields := strings.Split(content, ",")
		data := make([]byte, 0, len(fields))
		for _, field := range fields {
			trimmed := strings.TrimSpace(field)
			if trimmed == "" {
				continue
			}
			val, err := strconv.Atoi(trimmed)
			if err != nil {
				decodedPNGErr = err
				return
			}
			data = append(data, byte(val))
		}
		decodedPNGData = data
	})
	return decodedPNGData, decodedPNGErr
}

var ErrDecodedPNGFormat = errors.New("decoded png asset has unexpected format")
