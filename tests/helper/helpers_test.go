package helper_test

import (
	"testing"
	"time"

	"github.com/mnuddindev/jutsu-api/pkg/helper"
	"github.com/mnuddindev/jutsu-api/pkg/utils"
)

func TestFormatTitleParity(t *testing.T) {
	got := helper.FormatTitle("My Anime!! Season 1", "123")
	want := utils.FormatTitle("My Anime!! Season 1", "123")
	if got != want {
		t.Fatalf("FormatTitle mismatch: got %q, want %q", got, want)
	}
}

func TestCacheSetNoClient(t *testing.T) {
	if err := helper.SetCachedData("test:key", "value", time.Second); err != nil {
		t.Fatalf("SetCachedData returned error with nil client: %v", err)
	}
}

func TestCacheGetNoClient(t *testing.T) {
	var out string
	if err := helper.GetCachedData("test:key", &out); err != nil {
		t.Fatalf("GetCachedData returned error with nil client: %v", err)
	}
}
