package log

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type LogLevel string
type LogLevelValue byte

const (
	traceLevel LogLevelValue = iota + 1
	debugLevel
	infoLevel
	warningLevel
	errorLevel
	fatalLevel
	TRACE LogLevel = LogLevel("TRACE")
	DEBUG LogLevel = LogLevel("DEBUG")
	INFO  LogLevel = LogLevel("INFO")
	WARN  LogLevel = LogLevel("WARN")
	ERROR LogLevel = LogLevel("ERROR")
	FATAL LogLevel = LogLevel("FATAL")
)

type Logger interface {
	Tracef(format string, in ...interface{})
	Trace(in ...interface{})
	Debugf(format string, in ...interface{})
	Debug(in ...interface{})
	Infof(format string, in ...interface{})
	Info(in ...interface{})
	Warnf(format string, in ...interface{})
	Warn(in ...interface{})
	Errorf(format string, in ...interface{})
	Error(in ...interface{})
	Fatalf(format string, in ...interface{})
	Fatal(in ...interface{})
	Printf(format string, in ...interface{})
	Println(in ...interface{})
	SetVerbosity(verbosity LogLevel)
	GetVerbosity() LogLevel
}

type logger struct {
	verbosity LogLevelValue
	logger    *log.Logger
}

func (logger *logger) Tracef(format string, in ...interface{}) {
	logger.log(traceLevel, fmt.Sprintf(format, in...))
}

func (logger *logger) Trace(in ...interface{}) {
	logger.log(traceLevel, in...)
}

func (logger *logger) Debugf(format string, in ...interface{}) {
	logger.log(debugLevel, fmt.Sprintf(format, in...))
}

func (logger *logger) Debug(in ...interface{}) {
	logger.log(debugLevel, in...)
}

func (logger *logger) Infof(format string, in ...interface{}) {
	logger.log(infoLevel, fmt.Sprintf(format, in...))
}

func (logger *logger) Info(in ...interface{}) {
	logger.log(infoLevel, in...)
}

func (logger *logger) Warnf(format string, in ...interface{}) {
	logger.log(warningLevel, fmt.Sprintf(format, in...))
}

func (logger *logger) Warn(in ...interface{}) {
	logger.log(warningLevel, in...)
}

func (logger *logger) Errorf(format string, in ...interface{}) {
	logger.log(errorLevel, fmt.Sprintf(format, in...))
}

func (logger *logger) Error(in ...interface{}) {
	logger.log(errorLevel, in...)
}

func (logger *logger) Fatalf(format string, in ...interface{}) {
	logger.log(fatalLevel, fmt.Sprintf(format, in...))
}

func (logger *logger) Fatal(in ...interface{}) {
	logger.log(fatalLevel, in...)
}

func (logger *logger) Printf(format string, in ...interface{}) {
	logger.logger.Printf(format, in...)
}

func (logger *logger) Println(in ...interface{}) {
	fmt.Println(in...)
}

func (logger *logger) SetVerbosity(verbosity LogLevel) {
	logger.verbosity = toVerbosityLevelValue(verbosity)
}
func (logger *logger) GetVerbosity() LogLevel {
	return toVerbosityLevel(logger.verbosity)
}

func (logger *logger) log(level LogLevelValue, in ...interface{}) {
	if level >= logger.verbosity {
		var itfs string = " " + string(toVerbosityLevel(level)) + " " + fmt.Sprint(in...)
		logger.logger.Println(itfs)
	}
}

func NewLogger(verbosity LogLevel) Logger {
	return &logger{
		verbosity: toVerbosityLevelValue(verbosity),
		logger:    log.New(os.Stdout, "[go-deploy] ", log.LstdFlags|log.LUTC),
	}
}

func NewAppLogger(appName string, verbosity LogLevel) Logger {
	return &logger{
		verbosity: toVerbosityLevelValue(verbosity),
		logger:    log.New(os.Stdout, "["+appName+"] ", log.LstdFlags|log.LUTC),
	}
}

func VerbosityLevelFromString(verbosity string) LogLevel {
	switch strings.ToUpper(verbosity) {
	case "TRACE":
		return TRACE
	case "DEBUG":
		return DEBUG
	case "INFO":
		return INFO
	case "WARN":
		return WARN
	case "ERROR":
		return ERROR
	case "FATAL":
		return FATAL
	}
	return INFO
}

func toVerbosityLevelValue(verbosity LogLevel) LogLevelValue {
	switch strings.ToUpper(string(verbosity)) {
	case "TRACE":
		return traceLevel
	case "DEBUG":
		return debugLevel
	case "INFO":
		return infoLevel
	case "WARN":
		return warningLevel
	case "ERROR":
		return errorLevel
	case "FATAL":
		return fatalLevel
	}
	return infoLevel
}

func toVerbosityLevel(verbosity LogLevelValue) LogLevel {
	switch verbosity {
	case traceLevel:
		return TRACE
	case debugLevel:
		return DEBUG
	case infoLevel:
		return INFO
	case warningLevel:
		return WARN
	case errorLevel:
		return ERROR
	case fatalLevel:
		return FATAL
	}
	return INFO
}
