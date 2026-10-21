// Package logging membuat logger zap untuk aplikasi.
package logging

import (
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New membuat logger console ke stderr: level warn default, info bila verbose.
// Logger dikembalikan dan diteruskan sebagai parameter; tidak ada state global.
func New(verbose bool) *zap.Logger {
	return newLogger(os.Stderr, verbose)
}

func newLogger(w io.Writer, verbose bool) *zap.Logger {
	level := zap.WarnLevel
	if verbose {
		level = zap.InfoLevel
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), zapcore.Lock(zapcore.AddSync(w)), level)
	return zap.New(core)
}
