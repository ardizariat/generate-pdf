package report

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ardizariat/generate-pdf/internal/data"
)

func TestCoreUnknownType(t *testing.T) {
	out := filepath.Join(t.TempDir(), "unknown.pdf")

	r, err := data.Generate("invoice", 8, 42)
	require.NoError(t, err)
	r.Type = "tidak-ada"

	pages, size, err := Generate(r, out, true)
	require.Error(t, err)
	assert.Contains(t, err.Error(), `jenis dokumen tidak didukung: "tidak-ada"`)
	assert.Zero(t, pages)
	assert.Zero(t, size)

	_, statErr := os.Stat(out)
	assert.True(t, errors.Is(statErr, os.ErrNotExist), "file tidak boleh tercipta: %v", statErr)
}

func TestCoreDeterministic(t *testing.T) {
	dir := t.TempDir()

	first := filepath.Join(dir, "first.pdf")
	second := filepath.Join(dir, "second.pdf")

	r1, err := data.Generate("invoice", 8, 42)
	require.NoError(t, err)
	pages, size, err := Generate(r1, first, true)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, pages, 1)
	assert.Positive(t, size)

	r2, err := data.Generate("invoice", 8, 42)
	require.NoError(t, err)
	_, _, err = Generate(r2, second, true)
	require.NoError(t, err)

	a, err := os.ReadFile(first)
	require.NoError(t, err)
	b, err := os.ReadFile(second)
	require.NoError(t, err)
	assert.Equal(t, a, b, "dua generate dengan seed sama harus byte-identik")
}

func TestCoreVerifyRejectsTruncated(t *testing.T) {
	path := filepath.Join(t.TempDir(), "truncated.pdf")

	r, err := data.Generate("kuitansi", 0, 42)
	require.NoError(t, err)
	_, _, err = Generate(r, path, true)
	require.NoError(t, err)

	content, err := os.ReadFile(path)
	require.NoError(t, err)
	require.Greater(t, len(content), 50)
	require.NoError(t, os.WriteFile(path, content[:len(content)-50], 0o600))

	pages, err := Verify(path)
	require.Error(t, err)
	assert.Zero(t, pages)
	assert.Contains(t, err.Error(), path)
}

func TestCoreVerifyFailureIsNotSilent(t *testing.T) {
	original := verifyPDF
	verifyPDF = func(string) (int, error) {
		return 0, errors.New("rusak")
	}
	t.Cleanup(func() { verifyPDF = original })

	r, err := data.Generate("invoice", 8, 42)
	require.NoError(t, err)

	pages, _, err := Generate(r, filepath.Join(t.TempDir(), "invoice.pdf"), true)
	require.Error(t, err)
	assert.Zero(t, pages)
	assert.Contains(t, err.Error(), "verifikasi")
}
