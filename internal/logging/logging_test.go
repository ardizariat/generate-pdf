package logging

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestNewLevel(t *testing.T) {
	quiet := New(false)
	assert.False(t, quiet.Core().Enabled(zapcore.InfoLevel))
	assert.True(t, quiet.Core().Enabled(zapcore.WarnLevel))

	verbose := New(true)
	assert.True(t, verbose.Core().Enabled(zapcore.InfoLevel))
}

func TestNewWritesInfoOnlyWhenVerbose(t *testing.T) {
	var quiet bytes.Buffer
	quietLogger := newLogger(&quiet, false)
	quietLogger.Info("hidden")
	quietLogger.Warn("shown")
	assert.NotContains(t, quiet.String(), "hidden")
	assert.Contains(t, quiet.String(), "shown")

	var verbose bytes.Buffer
	verboseLogger := newLogger(&verbose, true)
	verboseLogger.Info("shown-info")
	assert.Contains(t, verbose.String(), "shown-info")
}

func TestNewUsesConsoleEncoder(t *testing.T) {
	var buf bytes.Buffer
	newLogger(&buf, true).Info("hello", zap.String("type", "invoice"))
	assert.Contains(t, buf.String(), `"type": "invoice"`)
	assert.NotContains(t, buf.String(), `{"level"`)
}
