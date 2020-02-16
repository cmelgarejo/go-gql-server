package logger

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

var logger *StandardLogger

func init() {
	logger = NewLogger()
}

// Event stores messages to log later, from our standard interface
type Event struct {
	id      int
	message string
}

// StandardLogger enforces specific log message formats
type StandardLogger struct {
	*logrus.Logger
}

// NewLogger initializes the standard logger
func NewLogger() *StandardLogger {
	var baseLogger = logrus.New()

	var standardLogger = &StandardLogger{baseLogger}

	standardLogger.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
	}
	// We could transform the errors into a JSON format, for external log SaaS tools
	// standardLogger.Formatter = &logrus.JSONFormatter{
	// 	PrettyPrint: true,
	// }

	return standardLogger
}

// Declare variables to store log messages as new Events
var (
	invalidArgMessage      = Event{1, "Invalid arg: %s"}
	invalidArgValueMessage = Event{2, "Invalid value for argument: %s: %v"}
	missingArgMessage      = Event{3, "Missing arg: %s"}
)

// Expose some log functions:

// Debug Log
func Debug(args ...interface{}) {
	logger.Debugln(args...)
}

// Debugf Log
func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

// Errorfn Log errors of a [fn] with format
func Errorfn(fn string, err error) error {
	outerr := fmt.Errorf("[%s]: %v", fn, err)
	logger.Errorln(outerr)
	return outerr
}

// InvalidArg is a standard error message
func InvalidArg(argumentName string) error {
	outerr := fmt.Errorf(invalidArgMessage.message, argumentName)
	logger.Errorln(outerr)
	return outerr
}

// InvalidArgValue is a standard error message
func InvalidArgValue(argumentName string, argumentValue string) error {
	outerr := fmt.Errorf(invalidArgValueMessage.message, argumentName, argumentValue)
	logger.Errorln(outerr)
	return outerr
}

// MissingArg is a standard error message
func MissingArg(argumentName string) error {
	outerr := fmt.Errorf(missingArgMessage.message, argumentName)
	logger.Errorln(outerr)
	return outerr
}

// Info Log
func Info(args ...interface{}) {
	logger.Infoln(args...)
}

// Infof Log
func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

// Warn Log
func Warn(args ...interface{}) {
	logger.Warnln(args...)
}

// Warnf Log
func Warnf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}

// Panic Log
func Panic(args ...interface{}) {
	logger.Panicln(args...)
}

// Panicf Log
func Panicf(format string, args ...interface{}) {
	logger.Panicf(format, args...)
}

// Error Log
func Error(args ...interface{}) {
	logger.Errorln(args...)
}

// Errorf Log
func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

// Fatal Log
func Fatal(args ...interface{}) {
	logger.Fatalln(args...)
}

// Fatalf Log
func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
}
