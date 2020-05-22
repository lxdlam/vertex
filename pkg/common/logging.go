package common

import (
	"errors"
	"fmt"
	"os"
	"time"

	"go.uber.org/zap/zapcore"

	"go.uber.org/zap"
)

var (
	logger     *zap.SugaredLogger = nil
	syncPeriod time.Duration      = 30 * time.Second
	exitChan   chan bool          = make(chan bool)
	debugMode  bool               = false
)

// InitLog will initialize the logger using the given config, should be called at the program start
func InitLog(c Config, debug bool) {
	debugMode = debug

	syncer := getSyncer(c)
	encoder := getEncoder(c)
	logLevel := zap.LevelEnablerFunc(getLogLevel(c))

	logger = zap.New(zapcore.NewCore(encoder, syncer, logLevel), zap.AddCaller(), zap.AddStacktrace(zapcore.FatalLevel)).Sugar()

	go checkAndSync()
}

func getSyncer(c Config) zapcore.WriteSyncer {
	file, err := os.OpenFile(c.LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(fmt.Sprintf("get syncer: creative file %s failed, err=%v", c.LogPath, err))
	}

	return zapcore.AddSync(file)
}

func getEncoder(c Config) zapcore.Encoder {
	if debugMode {
		return zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	}
	return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
}

func getLogLevel(c Config) func(zapcore.Level) bool {
	level := zapcore.InfoLevel

	if debugMode {
		level = zapcore.DebugLevel
	} else {
		// No debug level support in production
		switch c.LogLevel {
		case "INFO":
			level = zapcore.InfoLevel
		case "WARN":
			level = zapcore.WarnLevel
		case "ERROR":
			level = zapcore.ErrorLevel
		case "FATAL":
			level = zapcore.FatalLevel
		}
	}

	return func(lvl zapcore.Level) bool {
		return lvl >= level
	}
}

// ShutdownLog will do the final sync work, should be called when program exits
func ShutdownLog() {
	close(exitChan)
	_ = logger.Sync()
}

func checkAndSync() {
	for {
		select {
		case <-exitChan:
			return
		case <-time.After(syncPeriod):
			_ = logger.Sync()
		}
	}
}

// Debug will generate a new log into to the logger
func Debug(log string) {
	if logger != nil {
		logger.Debug(log)
	}
}

// Debugf works like Printf and works same as Debug
func Debugf(template string, a ...interface{}) {
	if logger != nil {
		logger.Debugf(template, a)
	}
}

// Info will generate a new info log into to the logger
func Info(log string) {
	if logger != nil {
		logger.Info(log)
	}
}

// Infof works like Printf and works same as Info
func Infof(template string, a ...interface{}) {
	if logger != nil {
		logger.Infof(template, a)
	}
}

// Warn will generate a new warn log into to the logger
func Warn(log string) {
	if logger != nil {
		logger.Warn(log)
	}
}

// Warnf works like Printf and works same as Warn
func Warnf(template string, a ...interface{}) {
	if logger != nil {
		logger.Warnf(template, a)
	}
}

// Error will generate a new error log into to the logger, a new error with the same message will also be raised
func Error(log string) error {
	if logger != nil {
		logger.Error(log)
	}

	return errors.New(log)
}

// Errorf works like Printf and works same as Error
func Errorf(template string, a ...interface{}) error {
	if logger != nil {
		logger.Errorf(template, a...)
	}

	return fmt.Errorf(template, a...)
}

// Panic will generate a new panic log into to the logger,
// and then call a panic() with the given log
func Panic(log string) {
	if logger != nil {
		logger.Panic(log)
	}
}

// Panicf works like Printf and works same as Panic
func Panicf(template string, a ...interface{}) {
	if logger != nil {
		logger.Panicf(template, a)
	}
}

// Fatal will generate a new fatal log into to the logger
func Fatal(log string) {
	if logger != nil {
		logger.Fatal(log)
	}
}

// Fatalf works like Printf and works same as Fatal
func Fatalf(template string, a ...interface{}) {
	if logger != nil {
		logger.Fatalf(template, a)
	}
}
