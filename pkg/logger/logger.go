// pkg/logger/logger.go
package logger

import (
	"fmt"
	"log"
	"time"
)

// ANSI color codes
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	ColorGray   = "\033[90m"

	// Bold colors
	ColorBoldRed    = "\033[1;31m"
	ColorBoldGreen  = "\033[1;32m"
	ColorBoldYellow = "\033[1;33m"
	ColorBoldBlue   = "\033[1;34m"
	ColorBoldPurple = "\033[1;35m"
	ColorBoldCyan   = "\033[1;36m"

	// Background colors
	BgRed    = "\033[41m"
	BgGreen  = "\033[42m"
	BgYellow = "\033[43m"
	BgBlue   = "\033[44m"
)

type Logger struct {
	prefix string
}

func New(prefix string) *Logger {
	return &Logger{prefix: prefix}
}

// Info - синий цвет для информационных сообщений
func (l *Logger) Info(format string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, args...)
	log.Printf("%s[%s]%s %s[INFO]%s %s%s%s %s",
		ColorGray, timestamp, ColorReset,
		ColorBoldBlue, ColorReset,
		ColorCyan, l.prefix, ColorReset,
		message)
}

// Success - зеленый цвет для успешных операций
func (l *Logger) Success(format string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, args...)
	log.Printf("%s[%s]%s %s[SUCCESS]%s %s%s%s %s",
		ColorGray, timestamp, ColorReset,
		ColorBoldGreen, ColorReset,
		ColorCyan, l.prefix, ColorReset,
		message)
}

// Error - красный цвет для ошибок
func (l *Logger) Error(format string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, args...)
	log.Printf("%s[%s]%s %s[ERROR]%s %s%s%s %s%s%s",
		ColorGray, timestamp, ColorReset,
		ColorBoldRed, ColorReset,
		ColorCyan, l.prefix, ColorReset,
		ColorRed, message, ColorReset)
}

// Warning - желтый цвет для предупреждений
func (l *Logger) Warning(format string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, args...)
	log.Printf("%s[%s]%s %s[WARNING]%s %s%s%s %s%s%s",
		ColorGray, timestamp, ColorReset,
		ColorBoldYellow, ColorReset,
		ColorCyan, l.prefix, ColorReset,
		ColorYellow, message, ColorReset)
}

// Debug - фиолетовый цвет для отладочных сообщений
func (l *Logger) Debug(format string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, args...)
	log.Printf("%s[%s]%s %s[DEBUG]%s %s%s%s %s",
		ColorGray, timestamp, ColorReset,
		ColorBoldPurple, ColorReset,
		ColorCyan, l.prefix, ColorReset,
		message)
}

// Order - зеленый цвет для операций с заказами
func (l *Logger) Order(format string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, args...)
	log.Printf("%s[%s]%s %s[ORDER]%s %s%s%s %s%s%s",
		ColorGray, timestamp, ColorReset,
		BgGreen+ColorWhite, ColorReset,
		ColorCyan, l.prefix, ColorReset,
		ColorGreen, message, ColorReset)
}

// Database - голубой цвет для операций с БД
func (l *Logger) Database(format string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, args...)
	log.Printf("%s[%s]%s %s[DB]%s %s%s%s %s",
		ColorGray, timestamp, ColorReset,
		ColorBoldCyan, ColorReset,
		ColorCyan, l.prefix, ColorReset,
		message)
}

// HTTP - белый цвет для HTTP запросов
func (l *Logger) HTTP(method, path string, statusCode int, duration time.Duration) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	var statusColor string
	switch {
	case statusCode >= 500:
		statusColor = ColorBoldRed
	case statusCode >= 400:
		statusColor = ColorBoldYellow
	case statusCode >= 300:
		statusColor = ColorBoldCyan
	case statusCode >= 200:
		statusColor = ColorBoldGreen
	default:
		statusColor = ColorWhite
	}

	var methodColor string
	switch method {
	case "GET":
		methodColor = ColorBoldBlue
	case "POST":
		methodColor = ColorBoldGreen
	case "PUT":
		methodColor = ColorBoldYellow
	case "DELETE":
		methodColor = ColorBoldRed
	default:
		methodColor = ColorWhite
	}

	log.Printf("%s[%s]%s %s[HTTP]%s %s%s%s %s %s%d%s %s",
		ColorGray, timestamp, ColorReset,
		ColorBoldPurple, ColorReset,
		methodColor, method, ColorReset,
		path,
		statusColor, statusCode, ColorReset,
		ColorGray+duration.String()+ColorReset)
}

// Fatal - критическая ошибка с красным фоном
func (l *Logger) Fatal(format string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, args...)
	log.Fatalf("%s[%s]%s %s[FATAL]%s %s%s%s %s%s%s",
		ColorGray, timestamp, ColorReset,
		BgRed+ColorWhite, ColorReset,
		ColorCyan, l.prefix, ColorReset,
		ColorRed, message, ColorReset)
}

// Startup - информация о запуске (синий с жирным шрифтом)
func (l *Logger) Startup(format string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, args...)
	log.Printf("%s[%s]%s %s[STARTUP]%s %s%s%s %s%s%s",
		ColorGray, timestamp, ColorReset,
		BgBlue+ColorWhite, ColorReset,
		ColorCyan, l.prefix, ColorReset,
		ColorBoldBlue, message, ColorReset)
}

// Global logger functions
var defaultLogger = New("")

func Info(format string, args ...interface{}) {
	defaultLogger.Info(format, args...)
}

func Success(format string, args ...interface{}) {
	defaultLogger.Success(format, args...)
}

func Error(format string, args ...interface{}) {
	defaultLogger.Error(format, args...)
}

func Warning(format string, args ...interface{}) {
	defaultLogger.Warning(format, args...)
}

func Debug(format string, args ...interface{}) {
	defaultLogger.Debug(format, args...)
}

func Order(format string, args ...interface{}) {
	defaultLogger.Order(format, args...)
}

func Database(format string, args ...interface{}) {
	defaultLogger.Database(format, args...)
}

func HTTP(method, path string, statusCode int, duration time.Duration) {
	defaultLogger.HTTP(method, path, statusCode, duration)
}

func Fatal(format string, args ...interface{}) {
	defaultLogger.Fatal(format, args...)
}

func Startup(format string, args ...interface{}) {
	defaultLogger.Startup(format, args...)
}
