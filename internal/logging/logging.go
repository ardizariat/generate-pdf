// Package logging builds the application's zap logger.
package logging

import (
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New returns a console logger writing to stderr: warn level by default, info when verbose.
// The logger is returned and passed as a parameter; there is no global state.
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
