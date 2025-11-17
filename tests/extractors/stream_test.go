package extractors_test

import (
	"errors"
	"testing"

	"github.com/mnuddindev/jutsu-api/pkg/extractors"
	"github.com/mnuddindev/jutsu-api/pkg/parsers"
)

func TestFindServerMatchMatchesByName(t *testing.T) {
	servers := []extractors.ServerItem{
		{ServerName: "Vidcloud", DataID: "123"},
		{ServerName: "HD-1", DataID: "456"},
	}
	match := extractors.FindServerMatchForTest(servers, "vidcloud")
	if match == nil || match.DataID != "123" {
		t.Fatalf("expected match for Vidcloud, got %#v", match)
	}
}

func TestFindServerMatchMatchesByDataID(t *testing.T) {
	servers := []extractors.ServerItem{
		{ServerName: "Vidcloud", DataID: "123"},
	}
	match := extractors.FindServerMatchForTest(servers, "123")
	if match == nil || match.ServerName != "Vidcloud" {
		t.Fatalf("expected to match by data-id, got %#v", match)
	}
}

func TestResolveStreamingLinkPrimary(t *testing.T) {
	restore := extractors.SetStreamDecryptorsForTest(
		func(id, name, typ string) (parsers.DecryptedSources, error) {
			if id != "server-id" {
				t.Fatalf("expected megacloud to receive server-id, got %s", id)
			}
			return parsers.DecryptedSources{Server: name, Type: typ}, nil
		},
		func(epID, id, name, typ string, fallback bool) (parsers.DecryptedSources, error) {
			t.Fatalf("legacy decryptor should not be invoked")
			return parsers.DecryptedSources{}, nil
		},
	)
	defer restore()

	res, err := extractors.ResolveStreamingLinkForTest("episode", "server-id", "Vidcloud", "sub", false)
	if err != nil {
		t.Fatalf("ResolveStreamingLinkForTest returned error: %v", err)
	}
	if res.Server != "Vidcloud" || res.Type != "sub" {
		t.Fatalf("unexpected response: %+v", res)
	}
}

func TestResolveStreamingLinkFallbackFlow(t *testing.T) {
	restore := extractors.SetStreamDecryptorsForTest(
		func(id, name, typ string) (parsers.DecryptedSources, error) {
			return parsers.DecryptedSources{}, errors.New("boom")
		},
		func(epID, id, name, typ string, fallback bool) (parsers.DecryptedSources, error) {
			if fallback {
				if id != "episode-slug" {
					t.Fatalf("expected fallback to receive episode slug, got %s", id)
				}
			} else if id != "server-id" {
				t.Fatalf("expected legacy decryptor to receive server-id, got %s", id)
			}
			return parsers.DecryptedSources{Server: name, Type: typ}, nil
		},
	)
	defer restore()

	// Case 1: automatic legacy fallback after megacloud error.
	res, err := extractors.ResolveStreamingLinkForTest("episode", "server-id", "Vidcloud", "sub", false)
	if err != nil {
		t.Fatalf("legacy fallback returned error: %v", err)
	}
	if res.Server != "Vidcloud" {
		t.Fatalf("unexpected server: %+v", res)
	}

	// Case 2: explicit fallback flag should bypass server id requirement.
	res, err = extractors.ResolveStreamingLinkForTest("episode-slug", "", "HD-1", "dub", true)
	if err != nil {
		t.Fatalf("explicit fallback returned error: %v", err)
	}
	if res.Type != "dub" {
		t.Fatalf("unexpected type from fallback: %+v", res)
	}
}
