package config

// ProductionConfig holds production-specific optimizations
type ProductionConfig struct {
	EnablePrefork      bool
	EnableCompression  bool
	EnableETag         bool
	MaxRequestBodySize int

	EnableTrustedProxy bool
	TrustedProxies     []string

	EnableMetrics bool
	MetricsPath   string

	EnableRateLimit bool
	RateLimitRPS    int
}

// GetProductionConfig returns production-optimized configuration
func GetProductionConfig() *ProductionConfig {
	return &ProductionConfig{
		EnablePrefork:      true,
		EnableCompression:  true,
		EnableETag:         true,
		MaxRequestBodySize: 4 * 1024 * 1024, // 4MB
		EnableTrustedProxy: true,
		TrustedProxies:     []string{"0.0.0.0/0"},
		EnableMetrics:      true,
		MetricsPath:        "/metrics",
		EnableRateLimit:    true,
		RateLimitRPS:       100,
	}
}
