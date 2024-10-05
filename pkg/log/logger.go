package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

type LoggerContext struct {
	Debug    string
	File     *os.File
	FileName string
}

func (l *LoggerContext) Open() {

	var logLevel zapcore.Level
	switch l.Debug {
	case "debug":
		logLevel = zap.DebugLevel
	case "development":
		logLevel = zap.InfoLevel
	case "production":
		logLevel = zap.WarnLevel
	default:
		logLevel = zap.InfoLevel
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Setting up encoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	// Create a custom WriteSyncer
	customSyncer := &CustomWriteSyncer{writer: l.File}

	// Define log level encoder
	core := zapcore.NewCore(encoder, zapcore.AddSync(customSyncer), logLevel)
	// Construct the logger
	Log = zap.New(core)

	defer Log.Sync() // Flushes buffer, if any
}

func (l *LoggerContext) OpenFileLogs(fileName string) *LoggerContext {
	// Create a file to log to
	logFile, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	// Append log file
	l.File = logFile
	return l
}
