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
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Red     = "\033[31m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"
)

type CustomTextFormatter struct{}

func (f *CustomTextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var color string

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

	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	level := strings.ToUpper(entry.Level.String())

	logLine := fmt.Sprintf("[%s] %s%s%s: %s", timestamp, color, level, Reset, entry.Message)

	for key, value := range entry.Data {
		logLine += fmt.Sprintf(" %s=%v", key, value)
	}

	return []byte(logLine + "\n"), nil
}

var (
	Instance *logrus.Logger
	once     sync.Once
)

func InitLogger() *logrus.Logger {
	once.Do(func() {
		Instance = logrus.New()
		Instance.SetOutput(os.Stdout)
		Instance.SetFormatter(&CustomTextFormatter{})
		Instance.SetLevel(logrus.InfoLevel)
	})

	return Instance
}

func Middleware(c *gin.Context) {

	methodColor := getMethodColor(c.Request.Method)

	Instance.WithFields(logrus.Fields{
		"method": fmt.Sprintf("%s%s%s", methodColor, c.Request.Method, Reset),
		"path":   c.Request.URL.Path,
	}).Info("Incoming request")

	c.Next()

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
