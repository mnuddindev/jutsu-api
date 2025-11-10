package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

// GenerateID generates a random ID
func GenerateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// GenerateSlug generates a URL-friendly slug from a string
func GenerateSlug(s string) string {
	// Convert to lowercase
	s = strings.ToLower(s)
	// Replace spaces with hyphens
	s = strings.ReplaceAll(s, " ", "-")
	// Remove special characters
	s = strings.ReplaceAll(s, "_", "-")
	// Remove multiple consecutive hyphens
	for strings.Contains(s, "--") {
		s = strings.ReplaceAll(s, "--", "-")
	}
	// Remove leading and trailing hyphens
	s = strings.Trim(s, "-")
	return s
}

// FormatDuration formats a duration to a human-readable string
func FormatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.0fs", d.Seconds())
	}
	if d < time.Hour {
		return fmt.Sprintf("%.0fm", d.Minutes())
	}
	return fmt.Sprintf("%.0fh", d.Hours())
}

// TruncateString truncates a string to a maximum length
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// Contains checks if a string slice contains a string
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// Unique removes duplicates from a string slice
func Unique(slice []string) []string {
	keys := make(map[string]bool)
	var result []string
	for _, item := range slice {
		if !keys[item] {
			keys[item] = true
			result = append(result, item)
		}
	}
	return result
}

// PointerToString converts a string pointer to a string
func PointerToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// StringToPointer converts a string to a string pointer
func StringToPointer(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// Now returns the current UTC time
func Now() time.Time {
	return time.Now().UTC()
}

// ParseTime parses a time string in RFC3339 format
func ParseTime(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}

// FormatTime formats a time to RFC3339 format
func FormatTime(t time.Time) string {
	return t.Format(time.RFC3339)
}

// FormatTimeString formats a time string to RFC3339 format
func FormatTimeString(t time.Time) string {
	return t.Format(time.RFC3339)
}

