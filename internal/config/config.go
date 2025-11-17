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
	App    AppConfig
	Server ServerConfig
	Redis  RedisConfig
	Logger LoggerConfig
	Cors   CorsConfig
}

// AppConfig holds application-level configuration
type AppConfig struct {
	Name        string
	Version     string
	Environment string
	Debug       bool
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

var Cfg *Config

// LoadConfig loads configuration from environment variables and config files
func LoadConfig() (*Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		// .env file is optional
		fmt.Println("No .env file found, using environment variables")
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Set defaults
	setDefaults()

	config := &Config{
		App: AppConfig{
			Name:        getString("APP_NAME", "Jutsu API"),
			Version:     getString("APP_VERSION", "1.0.0"),
			Environment: getString("APP_ENV", "development"),
			Debug:       getBool("APP_DEBUG", false),
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
			Host:          getString("REDIS_HOST", "localhost"),
			Port:          getString("REDIS_PORT", "6379"),
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
