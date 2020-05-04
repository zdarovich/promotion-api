package log

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
	"github.com/zdarovich/promotion-api/internal/config"
)

var (
	ip   string
	host string
	pid  string
)

//New init ...
func New(configuration *config.Configuration) {
	if configuration == nil {
		return
	}

	var fieldMap = logrus.FieldMap{
		logrus.FieldKeyMsg:   "message",
		logrus.FieldKeyLevel: "loglevel",
	}

	var logLevel logrus.Level
	if configuration.LogDebugMode {
		logLevel = logrus.DebugLevel
	} else {
		logLevel = logrus.InfoLevel
	}
	logrus.SetLevel(logLevel)
	if configuration.LogsEnabled {
		rotateFileHook, err := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
			Filename:   configuration.LogFilePath + logLevel.String() + ".log." + time.Now().Format("01-02-2006"),
			MaxSize:    configuration.LogRotateMegabytes, // megabytes
			MaxBackups: configuration.LogRotateFiles,
			MaxAge:     configuration.LogRotateDuration, //days
			Level:      logLevel,
			Formatter: &logrus.JSONFormatter{
				TimestampFormat: time.RFC822,
				FieldMap:        fieldMap,
			},
		})
		logrus.AddHook(rotateFileHook)

		if err != nil {
			logrus.Fatalf("Failed to initialize file rotate hook: %v", err)
			logrus.SetOutput(colorable.NewColorableStdout())
		}

	} else {
		logrus.SetOutput(colorable.NewColorableStdout())
		logrus.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC822,
		})
	}

	var ipErr error

	ip, ipErr = externalIP()
	if ipErr != nil {
		ip = ipErr.Error()
	}
	hosts, err := net.LookupAddr(ip)
	if err != nil {
		host = err.Error()
	}
	host = strings.Join(hosts, " ")

	pid = fmt.Sprint(os.Getpid())

}

//Trace wrapper
func Trace(args ...interface{}) {
	withContext().Trace(args...)
}

//Debug wrapper
func Debug(args ...interface{}) {
	withContext().Debug(args...)
}

//Info wrapper
func Info(args ...interface{}) {
	withContext().Info(args...)
}

//Infof wrapper
func Infof(format string, args ...interface{}) {
	withContext().Infof(format, args...)
}

//Warn wrapper
func Warn(args ...interface{}) {
	withContext().Warn(args...)
}

//Error wrapper
func Error(args ...interface{}) {
	withContext().Error(args...)
}

//Errorf wrapper
func Errorf(format string, args ...interface{}) {
	withContext().Errorf(format, args...)
}

//Fatal wrapper
func Fatal(args ...interface{}) {
	withContext().Fatal(args...)
}

//Panic wrapper
func Panic(args ...interface{}) {
	withContext().Panic(args...)
}

//HTTP ...
func HTTP(req *http.Request, res *http.Response, err error, duration time.Duration, method string) {

	withContext().Infof(
		"Request %s %s %s",
		req.Method,
		method,
		req.URL.String(),
	)
	if err != nil {
		logrus.Error(err)
		return
	}

	duration /= time.Millisecond
	withContext().Infof(
		"Response status=%d durationMs=%d",
		res.StatusCode,
		duration,
	)
}
