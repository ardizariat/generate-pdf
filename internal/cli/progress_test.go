package cli

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ardizariat/generate-pdf/internal/config"
)

// captureStderr redirects os.Stderr for the duration of the test.
func captureStderr(t *testing.T) func() string {
	t.Helper()

	read, write, err := os.Pipe()
	require.NoError(t, err)

	previous := os.Stderr
	os.Stderr = write
	return func() string {
		os.Stderr = previous
		require.NoError(t, write.Close())

		content, err := io.ReadAll(read)
		require.NoError(t, err)
		require.NoError(t, read.Close())
		return string(content)
	}
}

func TestProgressPrintsPlainLineWithoutTerminal(t *testing.T) {
	captured := captureStderr(t)

	ran := 0
	progress(config.Config{}, "generate 1/2 · file.pdf", func() { ran++ })
	output := captured()

	assert.Equal(t, 1, ran)
	assert.Contains(t, output, "generate 1/2 · file.pdf")
}

func TestProgressSilentWhenVerbose(t *testing.T) {
	captured := captureStderr(t)

	ran := 0
	progress(config.Config{Verbose: true}, "generate 1/2 · file.pdf", func() { ran++ })
	output := captured()

	assert.Equal(t, 1, ran)
	assert.Empty(t, output)
}
