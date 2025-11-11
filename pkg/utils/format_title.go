package utils

import (
	"regexp"
	"strings"
)

func FormatTitle(title string, dataID string) string {
	re := regexp.MustCompile(`[^\w\s]`)
	formatted := re.ReplaceAllString(title, "")
	formatted = strings.ToLower(formatted)
	formatted = strings.Join(strings.Fields(formatted), "-")
	return formatted + "-" + dataID
}
