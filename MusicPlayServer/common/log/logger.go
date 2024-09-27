package log

import (
	"MusicPlayServer/common"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
)

var (
	logger *zap.Logger
	stdlog *log.Logger
)

func InitLog(cfg *common.Config, options ...LogOptions) error {
	var err error

	config := DefaultLogOptions
	switch cfg.Mode {
	case common.ModeTest:
		config.Env = common.ModeTest
		config.LogPath = cfg.Log.Output + "/" + config.ProcessName
	case common.ModeProd:
		config.Env = common.ModeProd
		config.Level = zapcore.InfoLevel
		config.LogPath = cfg.Log.Output + "/" + config.ProcessName
	}

	if err = checkLog(&config); err != nil {
		return err
	}

	for _, option := range options {
		option(&config)
	}

	if logger, err = LogInit(&config); err != nil {
		fmt.Printf("ZapLogInit err:%v", err)
		return err
	}

	logger = logger.WithOptions(zap.AddCallerSkip(1))

	// redirect standard go log to zap
	zap.RedirectStdLogAt(logger, config.Level)
	stdlog, _ = zap.NewStdLogAt(logger, config.Level)
	return nil
}

func checkLog(conf *logConf) error {
	if len(conf.LogPath) > 0 {
		if _, err := os.Stat(conf.LogPath); err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("log path does not exist: %s \n", conf.LogPath)

				if err := os.MkdirAll(conf.LogPath, os.ModePerm); err != nil {
					fmt.Printf("create log path error: %v\n", err)
					return err
				}
			} else {
				fmt.Printf("check log path error: %v\n", err)

				return err
			}
		}
	}

	return nil
}

func Logger() *zap.Logger {
	return logger
}

func StdLogger() *log.Logger {
	return stdlog
}

// Debug logs a message at DebugLevel.
func Debug(args ...interface{}) {
	logger.Debug(formatLog("", args))
}

func Debugf(format string, args ...interface{}) {
	logger.Debug(formatLog(format, args))
}

// Info logs a message at InfoLevel.
func Info(args ...interface{}) {
	logger.Info(formatLog("", args))
}

func Infof(format string, args ...interface{}) {
	logger.Info(formatLog(format, args))
}

// Warn logs a message at WarnLevel.
func Warn(args ...interface{}) {
	logger.Warn(formatLog("", args))
}

func Warnf(format string, args ...interface{}) {
	logger.Warn(formatLog(format, args))
}

// Error logs a message at ErrorLevel.
func Error(args ...interface{}) {
	logger.Error(formatLog("", args))
}

func Errorf(format string, args ...interface{}) {
	logger.Error(formatLog(format, args))
}

// DPanic logs a message at DPanicLevel.
//
// If the logger is in development mode, it then panics (DPanic means
// "development panic"). This is useful for catching errors that are
// recoverable, but shouldn't ever happen.
func DPanic(args ...interface{}) {
	logger.DPanic(formatLog("", args))
}

func DPanicf(format string, args ...interface{}) {
	logger.DPanic(formatLog(format, args))
}

// Panic logs a message at PanicLevel.
//
// The logger then panics, even if logging at PanicLevel is disabled.
func Panic(args ...interface{}) {
	logger.Panic(formatLog("", args))
}

func Panicf(format string, args ...interface{}) {
	logger.Panic(formatLog(format, args))
}

// Fatal logs a message at FatalLevel.
//
// The logger then calls os.Exit(1), even if logging at FatalLevel is disabled.
func Fatal(args ...interface{}) {
	logger.Fatal(formatLog("", args))
}

func Fatalf(format string, args ...interface{}) {
	logger.Fatal(formatLog(format, args))
}

func formatLog(template string, fmtArgs []interface{}) string {
	// Format with Sprint, Sprintf, or neither.
	msg := template
	if msg == "" && len(fmtArgs) > 0 {
		msg = fmt.Sprint(fmtArgs...)
	} else if msg != "" && len(fmtArgs) > 0 {
		msg = fmt.Sprintf(template, fmtArgs...)
	}

	return msg
}

func Sync() error {
	return logger.Sync()
}
