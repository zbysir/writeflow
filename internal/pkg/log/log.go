package log

import (
	"bytes"
	"github.com/zbysir/writeflow/internal/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"time"
)

func Debugf(template string, args ...interface{}) { innerLogger.Debugf(template, args...) }

func Infof(template string, args ...interface{})  { innerLogger.Infof(template, args...) }
func Warnf(template string, args ...interface{})  { innerLogger.Warnf(template, args...) }
func Errorf(template string, args ...interface{}) { innerLogger.Errorf(template, args...) }
func Fatalf(template string, args ...interface{}) { innerLogger.Fatalf(template, args...) }
func Panicf(template string, args ...interface{}) { innerLogger.Panicf(template, args...) }

var innerLogger *zap.SugaredLogger

func Logger() *zap.SugaredLogger {
	return innerLogger
}

func init() {
	SetDev(false)
}

// SetDev 会影响最低等级
func SetDev(logDebug bool) {
	// 如果开启了 env  DEBUG=true，才会打印 Caller
	disableCaller := !config.IsDebug()

	innerLogger = New(Options{
		IsDev:         logDebug,
		To:            nil,
		DisableTime:   false,
		DisableCaller: disableCaller,
		CallerSkip:    1,
		Name:          "",
	})
}

type BuffSink struct {
	buf bytes.Buffer
}

func (b *BuffSink) Write(p []byte) (n int, err error) {
	return b.buf.Write(p)
}

func (b *BuffSink) Sync() error {
	return nil
}

func (b *BuffSink) Close() error {
	return nil
}

type Options struct {
	IsDev         bool
	To            io.Writer
	DisableTime   bool
	DisableLevel  bool
	DisableCaller bool
	CallerSkip    int
	Name          string
}

func New(o Options) *zap.SugaredLogger {
	zapconfig := zap.NewDevelopmentConfig()
	if !o.IsDev {
		zapconfig.Level.SetLevel(zap.InfoLevel)
	}

	//config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	if o.DisableTime {
		zapconfig.EncoderConfig.EncodeTime = nil
	} else {
		zapconfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.StampMilli)
	}
	if o.DisableLevel {
		zapconfig.EncoderConfig.EncodeLevel = nil
	} else {
		zapconfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	var ops []zap.Option

	ops = append(ops)
	if o.CallerSkip != 0 {
		ops = append(ops, zap.AddCallerSkip(o.CallerSkip))
	}

	var sink zapcore.WriteSyncer
	if o.To == nil {
		var err error
		var closeOut func()
		sink, closeOut, err = zap.Open(zapconfig.OutputPaths...)
		if err != nil {
			panic(err)

		}
		errSink, _, err := zap.Open(zapconfig.ErrorOutputPaths...)
		if err != nil {
			closeOut()
			panic(err)
		}

		ops = append(ops, zap.ErrorOutput(errSink))
	} else {
		sink = zapcore.AddSync(o.To)
	}
	if !o.DisableCaller {
		ops = append(ops, zap.AddCaller())
	}
	logger := zap.New(zapcore.NewCore(zapcore.NewConsoleEncoder(zapconfig.EncoderConfig), sink, zapconfig.Level), ops...)

	sugar := logger.Sugar()
	if o.Name != "" {
		sugar = sugar.Named(o.Name)
	}
	return sugar
}
