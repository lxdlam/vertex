package common

import (
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
	logger.Debug(log)
}

// Debugf works like Printf and works same as Debug
func Debugf(template string, a ...interface{}) {
	logger.Debugf(template, a)
}

// Info will generate a new log into to the logger
func Info(log string) {
	logger.Info(log)
}

// Infof works like Printf and works same as Info
func Infof(template string, a ...interface{}) {
	logger.Infof(template, a)
}

// Warn will generate a new log into to the logger
func Warn(log string) {
	logger.Warn(log)
}

// Warnf works like Printf and works same as Warn
func Warnf(template string, a ...interface{}) {
	logger.Warnf(template, a)
}

// Error will generate a new log into to the logger
func Error(log string) {
	logger.Error(log)
}

// Errorf works like Printf and works same as Error
func Errorf(template string, a ...interface{}) {
	logger.Errorf(template, a)
}

// Fatal will generate a new log into to the logger
func Fatal(log string) {
	logger.Fatal(log)
}

// Fatalf works like Printf and works same as Fatal
func Fatalf(template string, a ...interface{}) {
	logger.Fatalf(template, a)
}
