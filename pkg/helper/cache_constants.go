package helper

import "time"

// Cache TTL constants for different data types
// These are optimized based on data volatility and update frequency

const (
	// HomeInfoCacheTTL - Home page data (spotlights, trending, etc.)
	// Updated frequently, cache for shorter duration
	HomeInfoCacheTTL = 15 * time.Minute

	// CategoryCacheTTL - Category/genre listings
	// Relatively stable, can cache longer
	CategoryCacheTTL = 30 * time.Minute

	// AnimeInfoCacheTTL - Individual anime details
	// Stable data, cache for longer
	AnimeInfoCacheTTL = 1 * time.Hour

	// EpisodeListCacheTTL - Episode lists
	// Updates when new episodes are added
	EpisodeListCacheTTL = 30 * time.Minute

	// TopTenCacheTTL - Top 10 anime
	// Updates daily, cache for shorter duration
	TopTenCacheTTL = 20 * time.Minute

	// CreatorCacheTTL - Producer/Studio listings
	// Relatively stable
	CreatorCacheTTL = 30 * time.Minute

	// ScheduleCacheTTL - Anime schedule
	// Updates daily, cache for shorter duration
	ScheduleCacheTTL = 10 * time.Minute

	// SearchCacheTTL - Search results
	// Short cache as results may vary
	SearchCacheTTL = 5 * time.Minute

	// FilterCacheTTL - Filter results
	// Short cache as results may vary
	FilterCacheTTL = 5 * time.Minute

	// SuggestionCacheTTL - Search suggestions
	// Very short cache for real-time feel
	SuggestionCacheTTL = 2 * time.Minute

	// StreamInfoCacheTTL - Streaming information
	// Very short cache as links may expire
	StreamInfoCacheTTL = 5 * time.Minute

	// CharacterCacheTTL - Character details
	// Stable data, cache longer
	CharacterCacheTTL = 1 * time.Hour

	// VoiceActorCacheTTL - Voice actor details
	// Stable data, cache longer
	VoiceActorCacheTTL = 1 * time.Hour

	// QtipCacheTTL - Qtip data
	// Stable data
	QtipCacheTTL = 30 * time.Minute

	// WatchlistCacheTTL - User watchlist
	// User-specific, cache for shorter duration
	WatchlistCacheTTL = 10 * time.Minute

	// TopSearchCacheTTL - Top search keywords
	// Updates frequently
	TopSearchCacheTTL = 15 * time.Minute

	// RandomCacheTTL - Random anime
	// Should not be cached, but if needed, very short
	RandomCacheTTL = 0 // No cache for random

	// NextEpisodeScheduleCacheTTL - Next episode schedule
	// Updates when episodes air
	NextEpisodeScheduleCacheTTL = 10 * time.Minute

	// ServersCacheTTL - Available servers for episodes
	// Relatively stable, cache for medium duration
	ServersCacheTTL = 10 * time.Minute
)
