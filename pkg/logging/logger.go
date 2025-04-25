package logging

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
	"sync"
)

const (
	Reset   = "\033[0m"
	Green   = "\033[32m" // INFO & 2xx status
	Yellow  = "\033[33m" // WARN & 4xx status
	Red     = "\033[31m" // ERROR & 5xx status
	Blue    = "\033[34m" // 3xx status & POST method
	Magenta = "\033[35m" // PUT method
	Cyan    = "\033[36m" // GET method
	White   = "\033[37m" // Default
)

type CustomTextFormatter struct{}

func (f *CustomTextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var color string
	// Apply color based on the log level
	switch entry.Level {
	case logrus.InfoLevel:
		color = Green
	case logrus.WarnLevel:
		color = Yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		color = Red
	default:
		color = White
	}

	// Format the timestamp, log level, and message
	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	level := strings.ToUpper(entry.Level.String())

	// Color the log level and format the log as [timestamp] LEVEL: message
	logLine := fmt.Sprintf("[%s] %s%s%s: %s", timestamp, color, level, Reset, entry.Message)

	// Include any fields in the log entry
	for key, value := range entry.Data {
		logLine += fmt.Sprintf(" %s=%v", key, value)
	}

	// Return the formatted log line
	return []byte(logLine + "\n"), nil
}

var (
	Instance *logrus.Logger
	once     sync.Once
)

func InitLogger() {
	once.Do(func() {
		Instance = logrus.New()
		Instance.SetOutput(os.Stdout)
		Instance.SetFormatter(&CustomTextFormatter{})
		Instance.SetLevel(logrus.InfoLevel)
	})
}

func Middleware(c *gin.Context) {
	// Get color for HTTP method
	methodColor := getMethodColor(c.Request.Method)

	// Log the incoming request
	Instance.WithFields(logrus.Fields{
		"method": fmt.Sprintf("%s%s%s", methodColor, c.Request.Method, Reset),
		"path":   c.Request.URL.Path,
	}).Info("Incoming request")

	c.Next()

	// Log after handling the request
	statusCode := c.Writer.Status()
	statusColor := getStatusColor(statusCode)

	Instance.WithFields(logrus.Fields{
		"status": fmt.Sprintf("%s%d%s", statusColor, statusCode, Reset),
		"method": fmt.Sprintf("%s%s%s", methodColor, c.Request.Method, Reset),
		"path":   c.Request.URL.Path,
	}).Info("Request handled")
}

func getMethodColor(method string) string {
	switch method {
	case "GET":
		return Cyan
	case "POST":
		return Blue
	case "PUT":
		return Magenta
	case "DELETE":
		return Red
	default:
		return White
	}
}

func getStatusColor(status int) string {
	switch {
	case status >= 200 && status < 300:
		return Green
	case status >= 300 && status < 400:
		return Blue
	case status >= 400 && status < 500:
		return Yellow
	case status >= 500:
		return Red
	default:
		return White
	}
}
