package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	"github.com/ardizariat/generate-pdf/internal/config"
	"github.com/ardizariat/generate-pdf/internal/locale"
)

// execute runs the root command with args and captures its output.
func execute(t *testing.T, args ...string) (string, error) {
	t.Helper()

	cmd := New(zap.NewNop())
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs(args)

	err := cmd.Execute()
	return out.String(), err
}

// ansiEscape strips the styling lipgloss adds to `types` output.
var ansiEscape = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func TestTypesListsAllDocumentTypes(t *testing.T) {
	out, err := execute(t, "types")
	require.NoError(t, err)

	var listed []string
	for _, line := range strings.Split(strings.TrimSpace(ansiEscape.ReplaceAllString(out, "")), "\n") {
		fields := strings.Fields(line)
		require.NotEmpty(t, fields, "baris tanpa slug: %q", line)
		listed = append(listed, fields[0])
	}

	assert.Equal(t, config.Types(), listed, "urutan dan isi harus sama dengan config.Types()")
	for _, docType := range config.Types() {
		assert.Contains(t, out, locale.Get(locale.English).TypeDescriptions[docType])
	}
}

func TestInvalidTypeFailsWithoutDoubleMessage(t *testing.T) {
	out, err := execute(t, "-t", "salah", "--no-input")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Type")
	assert.Empty(t, out, "cobra tidak boleh mencetak error atau usage sendiri")
}

func TestUnknownFlagFails(t *testing.T) {
	_, err := execute(t, "--tidak-ada")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown flag")
}

func TestOverwriteRequiresForce(t *testing.T) {
	out := filepath.Join(t.TempDir(), "invoice.pdf")
	args := []string{"-t", "invoice", "-o", out, "--no-input", "-s", "42"}

	_, err := execute(t, args...)
	require.NoError(t, err)

	original, err := os.ReadFile(out)
	require.NoError(t, err)

	_, err = execute(t, args...)
	require.Error(t, err)
	assert.Contains(t, err.Error(), out)
	assert.Contains(t, err.Error(), "--force")

	after, err := os.ReadFile(out)
	require.NoError(t, err)
	assert.True(t, bytes.Equal(original, after), "file lama tidak boleh berubah tanpa --force")

	_, err = execute(t, append(args, "--force")...)
	require.NoError(t, err)
}

func TestGenerateWithVerify(t *testing.T) {
	out := filepath.Join(t.TempDir(), "receipt.pdf")

	_, err := execute(t, "-t", "receipt", "-o", out, "--no-input", "--verify")
	require.NoError(t, err)

	info, err := os.Stat(out)
	require.NoError(t, err)
	assert.Positive(t, info.Size())

	header, err := os.ReadFile(out)
	require.NoError(t, err)
	assert.True(t, strings.HasPrefix(string(header), "%PDF-"))
}

func TestVerboseSelectsInfoLogger(t *testing.T) {
	core, logs := observer.New(zapcore.DebugLevel)
	sawVerbose := false

	original := newLogger
	newLogger = func(verbose bool) *zap.Logger {
		sawVerbose = verbose
		return zap.New(core)
	}
	t.Cleanup(func() { newLogger = original })

	out := filepath.Join(t.TempDir(), "invoice.pdf")
	_, err := execute(t, "-t", "invoice", "-o", out, "--no-input", "-v")
	require.NoError(t, err)

	assert.True(t, sawVerbose)
	entries := logs.All()
	require.Len(t, entries, 1)
	assert.Equal(t, zapcore.InfoLevel, entries[0].Level)
	assert.Equal(t, "generate complete", entries[0].Message)
}
