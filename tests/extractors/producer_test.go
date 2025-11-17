package extractors_test

import (
	"errors"
	"testing"

	"github.com/mnuddindev/jutsu-api/pkg/extractors"
	"github.com/mnuddindev/jutsu-api/pkg/scrape"
)

func TestExtractProducerUsesRouteAndPagination(t *testing.T) {
	restore := extractors.SetExtractPageFuncForTest(func(page int, params, baseURL string) ([]scrape.ExtractedItem, int, error) {
		if params != "producer/a-1-pictures" {
			t.Fatalf("unexpected params: %s", params)
		}
		if page != 2 {
			t.Fatalf("expected page 2, got %d", page)
		}
		if baseURL != "example.org" {
			t.Fatalf("unexpected baseURL: %s", baseURL)
		}
		return []scrape.ExtractedItem{
			{ID: "anime-1", Title: "Anime 1"},
		}, 5, nil
	})
	defer restore()

	result, err := extractors.ExtractProducer("a-1-pictures", 2, "example.org")
	if err != nil {
		t.Fatalf("ExtractProducer returned error: %v", err)
	}
	if result.TotalPages != 5 {
		t.Fatalf("expected totalPages 5, got %d", result.TotalPages)
	}
	if len(result.Data) != 1 || result.Data[0].ID != "anime-1" {
		t.Fatalf("unexpected data payload: %+v", result.Data)
	}
}

func TestExtractProducerDetectsInvalidPage(t *testing.T) {
	restore := extractors.SetExtractPageFuncForTest(func(page int, params, baseURL string) ([]scrape.ExtractedItem, int, error) {
		return nil, 1, nil
	})
	defer restore()

	_, err := extractors.ExtractProducer("some-producer", 3, "example.org")
	if !errors.Is(err, extractors.ErrCreatorPageOutOfRange) {
		t.Fatalf("expected ErrCreatorPageOutOfRange, got %v", err)
	}
}

func TestExtractProducerRequiresID(t *testing.T) {
	_, err := extractors.ExtractProducer("   ", 1, "example.org")
	if !errors.Is(err, extractors.ErrCreatorIDRequired) {
		t.Fatalf("expected ErrCreatorIDRequired, got %v", err)
	}
}

func TestExtractStudioSwitchesPrefix(t *testing.T) {
	restore := extractors.SetExtractPageFuncForTest(func(page int, params, baseURL string) ([]scrape.ExtractedItem, int, error) {
		if params != "studio/bones" {
			t.Fatalf("expected studio prefix, got %s", params)
		}
		return nil, 1, nil
	})
	defer restore()

	if _, err := extractors.ExtractStudio("bones", 1, "example.org"); err != nil {
		t.Fatalf("ExtractStudio returned error: %v", err)
	}
}
