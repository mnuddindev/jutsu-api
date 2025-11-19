package utils

import "strings"

// Providers holds all provider information
type Providers struct {
	BaseProviders   []string
	FallbackStreams []string
	StreamProvider  string
	BaseV1          string
	BaseV2          string
	BaseV3          string
	BaseV4          string
	Fallback1       string
	Fallback2       string
}

const (
	baseV1Host = "hianime.do"
	baseV2Host = "kaido.to"
	baseV3Host = "aniplay.lol"
	baseV4Host = "9animetv.to"

	fallback1Host = "megaplay.buzz"
	fallback2Host = "vidwish.live"

	streamProviderHost = "megacloud.club"
)

var (
	// BaseProviders are the main anime providers
	BaseProviders = []string{
		baseV1Host,
		baseV2Host,
		baseV3Host,
		baseV4Host,
	}

	// FallbackStreams are fallback stream providers
	FallbackStreams = []string{
		fallback1Host,
		fallback2Host,
	}
)

// GetProviders returns all provider information
func GetProviders() Providers {
	return Providers{
		BaseProviders:   BaseProviders,
		FallbackStreams: FallbackStreams,
		StreamProvider:  ensureHTTPS(streamProviderHost),
		BaseV1:          baseV1Host,
		BaseV2:          baseV2Host,
		BaseV3:          baseV3Host,
		BaseV4:          baseV4Host,
		Fallback1:       fallback1Host,
		Fallback2:       fallback2Host,
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

// GetStreamProvider returns the main stream provider with https scheme
func GetStreamProvider() string {
	return ensureHTTPS(streamProviderHost)
}

// GetStreamProviderHost returns the raw stream provider host
func GetStreamProviderHost() string {
	return streamProviderHost
}

// GetV1BaseHost returns the v1 provider host
func GetV1BaseHost() string {
	return baseV1Host
}

// GetV2BaseHost returns the v2 provider host
func GetV2BaseHost() string {
	return baseV2Host
}

// GetV3BaseHost returns the v3 provider host
func GetV3BaseHost() string {
	return baseV3Host
}

// GetV4BaseHost returns the v4 provider host
func GetV4BaseHost() string {
	return baseV4Host
}

// GetV1BaseURL returns the v1 provider URL with https scheme
func GetV1BaseURL() string {
	return ensureHTTPS(baseV1Host)
}

// GetV2BaseURL returns the v2 provider URL with https scheme
func GetV2BaseURL() string {
	return ensureHTTPS(baseV2Host)
}

// GetV3BaseURL returns the v3 provider URL with https scheme
func GetV3BaseURL() string {
	return ensureHTTPS(baseV3Host)
}

// GetV4BaseURL returns the v4 provider URL with https scheme
func GetV4BaseURL() string {
	return ensureHTTPS(baseV4Host)
}

// GetFallback1Host returns the primary fallback host
func GetFallback1Host() string {
	return fallback1Host
}

// GetFallback2Host returns the secondary fallback host
func GetFallback2Host() string {
	return fallback2Host
}

// GetFallback1URL returns the primary fallback URL with https scheme
func GetFallback1URL() string {
	return ensureHTTPS(fallback1Host)
}

// GetFallback2URL returns the secondary fallback URL with https scheme
func GetFallback2URL() string {
	return ensureHTTPS(fallback2Host)
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

// ensureHTTPS prefixes https:// if missing
func ensureHTTPS(host string) string {
	if strings.HasPrefix(host, "http://") || strings.HasPrefix(host, "https://") {
		return host
	}
	return "https://" + host
}
