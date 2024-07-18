package ioc

import (
	"github.com/ac-zht/super-job/admin/pkg/logger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
)

func InitDevLogger() logger.Logger {
	cfg := zap.NewDevelopmentConfig()
	err := viper.UnmarshalKey("log", &cfg)
	if err != nil {
		panic(err)
	}
	l, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	return logger.NewZapLogger(l)
}

func InitProdLogger() logger.Logger {
	cfg := zap.NewProductionConfig()
	err := viper.UnmarshalKey("log", &cfg)
	if err != nil {
		panic(err)
	}
	l, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	return logger.NewZapLogger(l)
}

func WrapCore(file string) zap.Option {
	return zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return ZapCore(file)
	})
}

func ZapCore(file string) zapcore.Core {
	writeSyncer := zapcore.AddSync(logWriter(file))
	encoder := Encoder()
	return zapcore.NewCore(encoder, writeSyncer, zap.DebugLevel)
}

func logWriter(file string) io.Writer {
	return &lumberjack.Logger{
		Filename:   file,
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
	}
}

func Encoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}
