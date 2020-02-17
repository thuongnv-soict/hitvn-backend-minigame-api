package logger

import (
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

var mLog *logrus.Logger

/**
 * Creates a new Logger
 * See more:
 * https://esc.sh/blog/golang-logging-using-logrus/
 * https://github.com/rifflock/lfshook
 */
func NewLogger(logPath string, logPrefix string) *logrus.Logger {
	if mLog != nil {
		return mLog
	}

	logPathMap := lfshook.PathMap {
		logrus.InfoLevel:	logPath + "/" + logPrefix + "_success.log",
		logrus.TraceLevel:	logPath + "/" + logPrefix + "_success.log",
		logrus.WarnLevel:	logPath + "/" + logPrefix + "_success.log",

		logrus.DebugLevel: logPath + "/" + logPrefix + "_debug.log",

		logrus.ErrorLevel:	logPath + "/" + logPrefix + "_error.log",
		logrus.FatalLevel:	logPath + "/" + logPrefix + "_error.log",
		logrus.PanicLevel:	logPath + "/" + logPrefix + "_error.log",
	}

	logFormatter := new(logrus.TextFormatter)
	logFormatter.TimestampFormat = "02-01-2006 15:04:05"
	logFormatter.FullTimestamp = true

	mLog = logrus.New()
	mLog.Hooks.Add(lfshook.NewHook(
		logPathMap,
		logFormatter,
	))

	return mLog
}

/**
 * Creates a new rotation log
 */
func NewLoggerWithRotation(logPath string, logPrefix string) *logrus.Logger {
	if mLog != nil {
		return mLog
	}

	successPath	:= logPath + "/" + logPrefix + "_success.log"
	debugPath	:= logPath + "/" + logPrefix + "_debug.log"
	errorPath	:= logPath + "/" + logPrefix + "_error.log"

	rotationMaxAge := viper.GetInt64("Log.RotationMaxAge")
	rotationTime   := viper.GetInt64("Log.RotationTime")

	successWriter, _ := rotatelogs.New(
		successPath + ".%Y%m%d%H%M",
		rotatelogs.WithLinkName(successPath),
		rotatelogs.WithMaxAge(time.Duration(rotationMaxAge)*time.Second),
		rotatelogs.WithRotationTime(time.Duration(rotationTime)*time.Second),
	)

	debugWriter, _ := rotatelogs.New(
		debugPath + ".%Y%m%d%H%M",
		rotatelogs.WithLinkName(debugPath),
		rotatelogs.WithMaxAge(time.Duration(rotationMaxAge) * time.Second),
		rotatelogs.WithRotationTime(time.Duration(rotationTime) * time.Second),
	)

	errorWriter, _ := rotatelogs.New(
		errorPath + ".%Y%m%d%H%M",
		rotatelogs.WithLinkName(errorPath),
		rotatelogs.WithMaxAge(time.Duration(rotationMaxAge)*time.Second),
		rotatelogs.WithRotationTime(time.Duration(rotationTime)*time.Second),
	)

	logFormatter := new(logrus.TextFormatter)
	logFormatter.TimestampFormat = "02-01-2006 15:04:05"
	logFormatter.FullTimestamp = true

	mLog = logrus.New()
	mLog.Hooks.Add(lfshook.NewHook(
		lfshook.WriterMap{
			logrus.InfoLevel:	successWriter,
			logrus.TraceLevel:	successWriter,
			logrus.WarnLevel:	successWriter,

			logrus.DebugLevel:	debugWriter,

			logrus.ErrorLevel:	errorWriter,
			logrus.FatalLevel:	errorWriter,
			logrus.PanicLevel:	errorWriter,
		},
		logFormatter,
	))

	return mLog
}

/**
 * Logs trace
 */
func Trace(format string, v ...interface{}) {
	mLog.Tracef(format, v)
}

/**
 * Logs info
 */
func Info(format string, v ...interface{}) {
	mLog.Infof(format, v)
}

/**
 * Logs warning
 */
func Warn(format string, v ...interface{}) {
	mLog.Warnf(format, v)
}

/**
 * Logs debug
 */
func Debug(format string, v ...interface{})  {
	mLog.Debugf(format, v)
}

/**
 * Logs error
 */
func Error(format string, v ...interface{}) {
	mLog.Errorf(format, v)
}

/**
 * Logs fatal
 */
func Fatal(format string, v ...interface{})  {
	mLog.Fatalf(format, v)
}

/**
 * Logs panic
 */
func Panic(format string, v ...interface{})  {
	mLog.Panicf(format, v)
}