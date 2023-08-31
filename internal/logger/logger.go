package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

// NewLogger initialize logger. Can set output.
func NewLogger(f *os.File) *zap.SugaredLogger {
	pe := zap.NewProductionEncoderConfig()
	pe.TimeKey = "timestamp"
	pe.EncodeTime = zapcore.ISO8601TimeEncoder

	fileEncoder := zapcore.NewJSONEncoder(pe)

	core := zapcore.NewCore(fileEncoder, zapcore.AddSync(f), zap.InfoLevel)

	l := zap.New(core)

	return l.Sugar()
}
