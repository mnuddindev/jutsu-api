package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Server   ServerConfig
	Redis    RedisConfig
	Logger   LoggerConfig
	Cors     CorsConfig
	Cache    CacheConfig
	JWT      JWTConfig
	User     UserConfig
}

// AppConfig holds application-level configuration
type AppConfig struct {
	Name             string
	Version          string
	Environment      string
	Debug            bool
	RateLimitEnabled bool
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	Database        string
	SSLMode         string
	MaxConnections  int32
	MaxConnLifetime int
	MaxConnIdleTime int
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Host         string
	Port         string
	ReadTimeout  int
	WriteTimeout int
	IdleTimeout  int
	Prefork      bool
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Enabled       bool
	Host          string
	Port          string
	Password      string
	DB            int
	MaxRetries    int
	PoolSize      int
	MinIdleConns  int
	DialTimeout   int
	ReadTimeout   int
	WriteTimeout  int
	PoolTimeout   int
	IdleTimeout   int
	IdleCheckFreq int
}

// LoggerConfig holds logger configuration
type LoggerConfig struct {
	Level      string
	Encoding   string
	OutputPath string
	ErrorPath  string
}

// CorsConfig holds CORS configuration
type CorsConfig struct {
	AllowedOrigins   []string
	AllowMethods     string
	AllowHeaders     string
	AllowCredentials bool
	ExposeHeaders    string
	MaxAge           int
}

// CacheConfig holds cache TTl configuration
type CacheConfig struct {
	Enabled        bool
	CharacterTTL   int
	ActorTTL       int
	StudioTTL      int
	ProducerTTL    int
	AnimeInfoTTL   int
	EpisodesTTL    int
	QtipTTL        int
	HomeTTL        int
	TopTenTTL      int
	TopSearchTTL   int
	TopAiringTTL   int
	GenreTTL       int
	ScheduleTTL    int
	NextEpisodeTTL int
	ServersTTL     int
	SearchTTL      int
	SuggestTTL     int
	RandomTTL      int
	StreamTTL      int
}

type JWTConfig struct {
	AccessSecretKey  string
	RefreshSecretKey string
	AccessTokenTTL   int
	RefreshTokenTTL  int
	Cost             int
}

type UserConfig struct {
	UserSessionTTL  int
	MaxUserSession  int
	UserProfileTTL  int
	JWTBlacklistTTL int
	RefreshTokenTTL int
}

var Cfg *Config

// LoadConfig loads configuration from environment variables and config files
func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, using environment variables")
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Set defaults
	setDefaults()

	config := &Config{
		App: AppConfig{
			Name:             getString("APP_NAME", "Jutsu API"),
			Version:          getString("APP_VERSION", "1.0.0"),
			Environment:      getString("APP_ENV", "development"),
			Debug:            getBool("APP_DEBUG", false),
			RateLimitEnabled: getBool("RATE_LIMIT_ENABLED", true),
		},
		Database: DatabaseConfig{
			Host:            getString("DB_HOST", "localhost"),
			Port:            getString("DB_PORT", "5432"),
			User:            getString("DB_USER", "mnuddin"),
			Password:        getString("DB_PASS", "mnuddin"),
			Database:        getString("DB_NAME", "jutsu"),
			SSLMode:         getString("DB_SSL_MODE", "disable"),
			MaxConnections:  int32(getInt("DB_MAX_CONN", 10)),
			MaxConnLifetime: getInt("DB_CONN_MAX_LIFETIME", 1800),
			MaxConnIdleTime: getInt("DB_MAX_IDLE_CONN", 5),
		},
		Server: ServerConfig{
			Host:         getString("SERVER_HOST", "0.0.0.0"),
			Port:         getString("SERVER_PORT", "8080"),
			ReadTimeout:  getInt("SERVER_READ_TIMEOUT", 15),
			WriteTimeout: getInt("SERVER_WRITE_TIMEOUT", 15),
			IdleTimeout:  getInt("SERVER_IDLE_TIMEOUT", 120),
			Prefork:      getBool("SERVER_PREFORK", false),
		},
		Redis: RedisConfig{
			Enabled:       getBool("REDIS_ENABLED", true),
			Host:          getString("REDIS_HOST", "localhost"),
			Port:          getString("REDIS_PORT", "10120"),
			Password:      getString("REDIS_PASSWORD", ""),
			DB:            getInt("REDIS_DB", 0),
			MaxRetries:    getInt("REDIS_MAX_RETRIES", 3),
			PoolSize:      getInt("REDIS_POOL_SIZE", 10),
			MinIdleConns:  getInt("REDIS_MIN_IDLE_CONNS", 5),
			DialTimeout:   getInt("REDIS_DIAL_TIMEOUT", 5),
			ReadTimeout:   getInt("REDIS_READ_TIMEOUT", 3),
			WriteTimeout:  getInt("REDIS_WRITE_TIMEOUT", 3),
			PoolTimeout:   getInt("REDIS_POOL_TIMEOUT", 4),
			IdleTimeout:   getInt("REDIS_IDLE_TIMEOUT", 5),
			IdleCheckFreq: getInt("REDIS_IDLE_CHECK_FREQ", 1),
		},
		Logger: LoggerConfig{
			Level:      getString("LOG_LEVEL", "info"),
			Encoding:   getString("LOG_ENCODING", "json"),
			OutputPath: getString("LOG_OUTPUT_PATH", "stdout"),
			ErrorPath:  getString("LOG_ERROR_PATH", "stderr"),
		},
		Cors: CorsConfig{
			AllowedOrigins:   getStringSlice("ALLOWED_ORIGINS", []string{"*"}),
			AllowMethods:     getString("CORS_ALLOW_METHODS", "GET,POST,PUT,DELETE,PATCH,OPTIONS"),
			AllowHeaders:     getString("CORS_ALLOW_HEADERS", "Origin,Content-Type,Accept,Authorization"),
			AllowCredentials: getBool("CORS_ALLOW_CREDENTIALS", false),
			ExposeHeaders:    getString("CORS_EXPOSE_HEADERS", "Content-Length"),
			MaxAge:           getInt("CORS_MAX_AGE", 3600),
		},
		Cache: CacheConfig{
			Enabled:        getBool("CACHE_ENABLED", true),
			CharacterTTL:   getInt("CACHE_TTL_CHARACTER", 86400),
			ActorTTL:       getInt("CACHE_TTL_ACTOR", 86400),
			StudioTTL:      getInt("CACHE_TTL_STUDIO", 86400),
			ProducerTTL:    getInt("CACHE_TTL_PRODUCER", 86400),
			AnimeInfoTTL:   getInt("CACHE_TTL_ANIME_INFO", 21600),
			EpisodesTTL:    getInt("CACHE_TTL_EPISODES", 3600),
			QtipTTL:        getInt("CACHE_TTL_QTIP", 21600),
			HomeTTL:        getInt("CACHE_TTL_HOME", 900),
			TopTenTTL:      getInt("CACHE_TTL_TOP_TEN", 3600),
			TopSearchTTL:   getInt("CACHE_TTL_TOP_SEARCH", 3600),
			TopAiringTTL:   getInt("CACHE_TTL_TOP_AIRING", 7200),
			GenreTTL:       getInt("CACHE_TTL_GENRE", 14400),
			ScheduleTTL:    getInt("CACHE_TTL_SCHEDULE", 600),
			NextEpisodeTTL: getInt("CACHE_TTL_NEXT_EPISODE", 1800),
			ServersTTL:     getInt("CACHE_TTL_SERVERS", 300),
			SearchTTL:      getInt("CACHE_TTL_SEARCH", 300),
			SuggestTTL:     getInt("CACHE_TTL_SUGGEST", 180),
			RandomTTL:      getInt("CACHE_TTL_RANDOM", 60),
			StreamTTL:      getInt("CACHE_TTL_STREAM", 0),
		},
		JWT: JWTConfig{
			AccessSecretKey:  getString("JWT_ACCESS_SECRET", "e85j2CFLDVHBc1geV6b8ur3309MteCYe"),
			RefreshSecretKey: getString("JWT_REFRESH_SECRET", "Di_iThHkq9eChxSTpSthrJnsDsdCiwfo"),
			AccessTokenTTL:   getInt("JWT_ACCESS_TTL", 900),
			RefreshTokenTTL:  getInt("JWT_REFRESH_TTL", 604800),
			Cost:             getInt("BCRYPT_COST", 12),
		},
		User: UserConfig{
			UserSessionTTL:  getInt("USER_SESSION_TTL", 86400),
			MaxUserSession:  getInt("USER_MAX_SESSION", 1),
			UserProfileTTL:  getInt("USER_PROFILE_TTL", 3600),
			JWTBlacklistTTL: getInt("JWT_BLACKLIST_TTL", 604800),
			RefreshTokenTTL: getInt("REFRESH_TOKEN_TTL", 604800),
		},
	}

	Cfg = config
	return config, nil
}

// GetRedisAddr returns the Redis address
func (c *Config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%s", c.Redis.Host, c.Redis.Port)
}

// GetServerAddr returns the server address
func (c *Config) GetServerAddr() string {
	return fmt.Sprintf("%s:%s", c.Server.Host, c.Server.Port)
}

// IsProduction checks if the application is running in production
func (c *Config) IsProduction() bool {
	return c.App.Environment == "production"
}

// IsDevelopment checks if the application is running in development
func (c *Config) IsDevelopment() bool {
	return c.App.Environment == "development"
}

// Helper functions
func getString(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	var intValue int
	if _, err := fmt.Sscanf(value, "%d", &intValue); err != nil {
		return defaultValue
	}
	if intValue == 0 {
		return defaultValue
	}
	return intValue
}

func getBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return strings.ToLower(value) == "true" || value == "1"
}

func getStringSlice(key string, defaultValue []string) []string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	// Split by comma and trim spaces
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	if len(result) == 0 {
		return defaultValue
	}
	return result
}

func setDefaults() {
	viper.SetDefault("APP_NAME", "Jutsu API")
	viper.SetDefault("APP_VERSION", "1.0.0")
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("APP_DEBUG", false)
}
