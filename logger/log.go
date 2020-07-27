package logger

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.SugaredLogger

type GenericLogger interface {
	Warn(args ...interface{})
	Warnf(template string, args ...interface{})
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(template string, args ...interface{})
}

// SetupLogger sets up the global logger to be used in all the cli, configuring
// the logging level to debug whenever verbose is true, and info otherwise.
// This function will act as a no-op after the first time it is called. Returns
// an error in case could not be built from the configuration
func SetupLogger(verbose bool) error {
	if logger != nil {
		return nil
	}
	config := zap.NewProductionConfig()
	config.Encoding = "console"
	config.EncoderConfig.EncodeTime = utcRFC3339TimeEncoder
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.MessageKey = "message"
	config.EncoderConfig.LevelKey = "severity"
	config.EncoderConfig.TimeKey = "timestamp"
	config.DisableStacktrace = true
	config.OutputPaths = []string{
		"stderr",
		//fmt.Sprintf("./%s-%s.log", env, time.Now().Format("20060102_030405")),
	}
	if verbose {
		config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	}
	zapLogger, err := config.Build()
	if err != nil {
		return err
	}
	logger = zapLogger.Sugar()
	return nil
}

func Logger() *zap.SugaredLogger {
	return logger
}

// UTCRFC3339TimeEncoder encodes the time as UTC and formats it according to RFC3339, i.e: 2020-03-12T17:51:03Z
func utcRFC3339TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	t = t.UTC()
	enc.AppendString(t.Format(time.RFC3339))
}
