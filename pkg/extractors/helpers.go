package extractors

import (
	"encoding/json"
	"fmt"
	"strings"
)

// lastSegment extracts the last segment from a URL path.
// Example: "/anime/123" -> "123", "/character/abc-123" -> "abc-123"
func lastSegment(href string) string {
	parts := strings.Split(strings.Trim(href, "/"), "/")
	if len(parts) == 0 {
		return ""
	}
	return parts[len(parts)-1]
}

// seg extracts a specific segment from a URL path by index.
// Example: seg("/anime/123/episode/456", 2) -> "123"
func seg(href string, idx int) string {
	parts := strings.Split(strings.Trim(href, "/"), "/")
	if idx >= 0 && idx < len(parts) {
		return parts[idx]
	}
	return ""
}

// firstNonEmpty returns the first non-empty string from the provided arguments.
func firstNonEmpty(strs ...string) string {
	for _, s := range strs {
		if strings.TrimSpace(s) != "" {
			return s
		}
	}
	return ""
}

// embedID extracts the video ID from an embed URL.
// Example: "https://youtube.com/embed/abc123" -> "abc123"
func embedID(u string) string {
	parts := strings.Split(u, "/embed/")
	if len(parts) < 2 {
		return ""
	}
	tail := parts[1]
	tail = strings.Split(tail, "?")[0]
	return strings.Split(tail, "&")[0]
}

// toString converts various types to string representation.
func toString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case float64:
		return fmt.Sprintf("%.0f", val)
	case int:
		return fmt.Sprintf("%d", val)
	case int64:
		return fmt.Sprintf("%d", val)
	case json.Number:
		return val.String()
	default:
		return fmt.Sprintf("%v", val)
	}
}
