package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/ardizariat/generate-pdf/internal/config"
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

func TestTypesListsAllDocumentTypes(t *testing.T) {
	out, err := execute(t, "types")
	require.NoError(t, err)

	for _, docType := range config.Types() {
		assert.Contains(t, out, docType)
		assert.Contains(t, out, typeDescriptions[docType])
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
	out := filepath.Join(t.TempDir(), "kuitansi.pdf")

	_, err := execute(t, "-t", "kuitansi", "-o", out, "--no-input", "--verify")
	require.NoError(t, err)

	info, err := os.Stat(out)
	require.NoError(t, err)
	assert.Positive(t, info.Size())

	header, err := os.ReadFile(out)
	require.NoError(t, err)
	assert.True(t, strings.HasPrefix(string(header), "%PDF-"))
}
