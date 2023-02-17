package logging

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func LogFile() *logrus.Entry {
	_, filename, line, ok := runtime.Caller(1)
	if !ok {
		panic("Unable to get caller information")
	}
	logFilename := log.WithField("filename", filepath.Base(filename)).WithField("line", line)
	return logFilename
}

func SetupLogging() *logrus.Logger {
	log := logrus.New()

	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	log.SetFormatter(customFormatter)

	// Log as JSON instead of the default ASCII formatter.
	//log.SetFormatter(&logrus.JSONFormatter{})

	// Output to stderr instead of stdout, could also be a file.
	log.SetOutput(os.Stderr)

	// Will log anything that is info or above (warn, error, fatal, panic). Default.
	logLevel := os.Getenv("LOG_LEVEL")
	switch logLevel {
	case "debug":
		log.SetLevel(logrus.DebugLevel)
		log.SetReportCaller(true)
	case "info":
		log.SetLevel(logrus.InfoLevel)
	case "warn":
		log.SetLevel(logrus.WarnLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	case "fatal":
		log.SetLevel(logrus.FatalLevel)
	case "panic":
		log.SetLevel(logrus.PanicLevel)
	default:
		log.SetLevel(logrus.InfoLevel)
	}

	return log
}
