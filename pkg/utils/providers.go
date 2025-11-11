package utils

// Providers holds all provider information
type Providers struct {
	BaseProviders   []string
	FallbackStreams []string
	StreamProvider  string
}

var (
	// BaseProviders are the main anime providers
	BaseProviders = []string{
		"hianime.do",
		"kaido.to",
		"aniplay.lol",
		"9animetv.to",
	}

	// FallbackStreams are fallback stream providers
	FallbackStreams = []string{
		"megaplay.buzz",
		"vidwish.live",
	}

	// StreamProvider is the main stream provider
	StreamProvider = "https://megacloud.club"
)

// GetProviders returns all provider information
func GetProviders() Providers {
	return Providers{
		BaseProviders:   BaseProviders,
		FallbackStreams: FallbackStreams,
		StreamProvider:  StreamProvider,
	}
}

// GetBaseProviders returns base providers
func GetBaseProviders() []string {
	return BaseProviders
}

// GetFallbackStreams returns fallback stream providers
func GetFallbackStreams() []string {
	return FallbackStreams
}

// GetStreamProvider returns the main stream provider
func GetStreamProvider() string {
	return StreamProvider
}

// IsValidProvider checks if a provider is valid
func IsValidProvider(provider string) bool {
	for _, p := range BaseProviders {
		if p == provider {
			return true
		}
	}
	return false
}

// IsValidFallbackStream checks if a fallback stream is valid
func IsValidFallbackStream(stream string) bool {
	for _, s := range FallbackStreams {
		if s == stream {
			return true
		}
	}
	return false
}
