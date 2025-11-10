package logger

import (
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/mnuddindev/jutsu-api/internal/config"
)

var Logger *zap.Logger
var SugaredLogger *zap.SugaredLogger

// InitLogger initializes the zap logger with configuration
func InitLogger(cfg *config.LoggerConfig) error {
	var encoderConfig zapcore.EncoderConfig

	if config.Cfg.IsProduction() {
		encoderConfig = zap.NewProductionEncoderConfig()
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
	} else {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
		encoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	}

	var encoder zapcore.Encoder
	if cfg.Encoding == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// Set log level
	var level zapcore.Level
	switch strings.ToLower(cfg.Level) {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	case "fatal":
		level = zapcore.FatalLevel
	case "panic":
		level = zapcore.PanicLevel
	default:
		level = zapcore.InfoLevel
	}

	// Setup writers
	var cores []zapcore.Core

	// Console output (stdout)
	consoleWriter := zapcore.AddSync(os.Stdout)
	consoleCore := zapcore.NewCore(encoder, consoleWriter, level)
	cores = append(cores, consoleCore)

	// File output for production (if specified)
	if cfg.OutputPath != "stdout" && cfg.OutputPath != "" {
		fileWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   cfg.OutputPath,
			MaxSize:    100, // megabytes
			MaxBackups: 10,
			MaxAge:     30, // days
			Compress:   true,
			LocalTime:  true,
		})
		fileCore := zapcore.NewCore(encoder, fileWriter, level)
		cores = append(cores, fileCore)
	}

	// Error file output (if specified and different from output path)
	if cfg.ErrorPath != "stderr" && cfg.ErrorPath != "" && cfg.ErrorPath != cfg.OutputPath {
		errorWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   cfg.ErrorPath,
			MaxSize:    100, // megabytes
			MaxBackups: 10,
			MaxAge:     30, // days
			Compress:   true,
			LocalTime:  true,
		})
		errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.ErrorLevel && level <= lvl
		})
		errorCore := zapcore.NewCore(encoder, errorWriter, errorLevel)
		cores = append(cores, errorCore)
	}

	// Combine cores
	core := zapcore.NewTee(cores...)

	// Build logger
	Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	if config.Cfg.IsDevelopment() {
		Logger = Logger.WithOptions(zap.Development())
	}

	SugaredLogger = Logger.Sugar()

	return nil
}

// Sync flushes any buffered log entries
func Sync() {
	if Logger != nil {
		_ = Logger.Sync()
	}
}

// Info logs an info message
func Info(msg string, fields ...zap.Field) {
	Logger.Info(msg, fields...)
}

// Error logs an error message
func Error(msg string, fields ...zap.Field) {
	Logger.Error(msg, fields...)
}

// Warn logs a warning message
func Warn(msg string, fields ...zap.Field) {
	Logger.Warn(msg, fields...)
}

// Debug logs a debug message
func Debug(msg string, fields ...zap.Field) {
	Logger.Debug(msg, fields...)
}

// Fatal logs a fatal message and exits
func Fatal(msg string, fields ...zap.Field) {
	Logger.Fatal(msg, fields...)
}

// Panic logs a panic message and panics
func Panic(msg string, fields ...zap.Field) {
	Logger.Panic(msg, fields...)
}

// WithFields creates a child logger with the given fields
func WithFields(fields ...zap.Field) *zap.Logger {
	return Logger.With(fields...)
}

// LogRequest logs HTTP request details
func LogRequest(method, path string, statusCode int, latency time.Duration, err error) {
	fields := []zap.Field{
		zap.String("method", method),
		zap.String("path", path),
		zap.Int("status", statusCode),
		zap.Duration("latency", latency),
	}

	if err != nil {
		fields = append(fields, zap.Error(err))
		Error("HTTP Request", fields...)
	} else {
		Info("HTTP Request", fields...)
	}
}

// LogError logs an error with context
func LogError(err error, msg string, fields ...zap.Field) {
	if err != nil {
		allFields := append(fields, zap.Error(err))
		Error(msg, allFields...)
	}
}

// LogCacheOperation logs cache operation details
func LogCacheOperation(operation, key string, err error) {
	fields := []zap.Field{
		zap.String("operation", operation),
		zap.String("key", key),
	}

	if err != nil {
		fields = append(fields, zap.Error(err))
		Warn("Cache Operation", fields...)
	} else {
		Debug("Cache Operation", fields...)
	}
}

