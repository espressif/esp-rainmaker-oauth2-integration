package utils

import (
	"os"

	log "github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
)

var logger = &log.Logger{
	Out:   os.Stderr,
	Level: log.DebugLevel,
	Formatter: &easy.Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
		LogFormat:       "[%lvl%]: %time% - %msg%\n",
	},
}

func init() {
	var logLevel string
	logLevel = GetLogLevel()
	switch logLevel {
	case "Info":
		logger.SetLevel(log.InfoLevel)
	case "Warn":
		logger.SetLevel(log.WarnLevel)
	case "Error":
		logger.SetLevel(log.ErrorLevel)
	case "Debug":
		logger.SetLevel(log.DebugLevel)
	default:
		logger.SetLevel(log.InfoLevel)
	}
}

/* Following Log functions are supposed take two arguments :
 * RequsetID	(optional) AWSRequestID passed in context object to lambda handler
 * message	(compulsory) Log message to be printed
 */

func LogDebug(args ...string) {
	if len(args) == 2 {
		logger.Debug(args[0] + " " + args[1])
	} else {
		logger.Debug(args[0])
	}
}

func LogInfo(args ...string) {
	if len(args) == 2 {
		logger.Info(args[0] + " " + args[1])
	} else {
		logger.Info(args[0])
	}
}

func LogWarn(args ...string) {
	if len(args) == 2 {
		logger.Warn(args[0] + " " + args[1])
	} else {
		logger.Warn(args[0])
	}
}

func LogError(args ...string) {
	if len(args) == 2 {
		logger.Error(args[0] + " " + args[1])
	} else {
		logger.Error(args[0])
	}
}
