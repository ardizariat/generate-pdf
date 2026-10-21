package report

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
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

	dir := t.TempDir()
	out := filepath.Join(dir, "invoice.pdf")
	pages, _, err := Generate(r, out, true)
	require.Error(t, err)
	assert.Zero(t, pages)
	assert.Contains(t, err.Error(), "verifikasi")

	_, statErr := os.Stat(out)
	assert.True(t, errors.Is(statErr, os.ErrNotExist), "file gagal tidak boleh tercipta: %v", statErr)
	leftovers, err := filepath.Glob(filepath.Join(dir, ".pdfgen-*"))
	require.NoError(t, err)
	assert.Empty(t, leftovers, "file sementara harus dibersihkan")
}

func TestCoreVerifyFailureKeepsExistingFile(t *testing.T) {
	original := verifyPDF
	verifyPDF = func(string) (int, error) {
		return 0, errors.New("rusak")
	}
	t.Cleanup(func() { verifyPDF = original })

	r, err := data.Generate("invoice", 8, 42)
	require.NoError(t, err)

	dir := t.TempDir()
	out := filepath.Join(dir, "invoice.pdf")
	old := []byte("%PDF-lama")
	require.NoError(t, os.WriteFile(out, old, 0o600))

	_, _, err = Generate(r, out, true)
	require.Error(t, err)

	content, err := os.ReadFile(out)
	require.NoError(t, err)
	assert.Equal(t, old, content, "file lama tidak boleh tersentuh saat verifikasi gagal")
	info, err := os.Stat(out)
	require.NoError(t, err)
	if runtime.GOOS != "windows" {
		assert.Equal(t, os.FileMode(0o600), info.Mode().Perm())
	}
}

func TestCoreGeneratePagesFailureLeavesNoFile(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "two.pdf")
	build := func(int) (data.Report, error) {
		return data.Generate("kuitansi", 0, 42)
	}

	_, _, _, err := GeneratePages(build, out, 2, true)
	require.Error(t, err)

	_, statErr := os.Stat(out)
	assert.True(t, errors.Is(statErr, os.ErrNotExist), "file gagal tidak boleh tercipta: %v", statErr)
	leftovers, err := filepath.Glob(filepath.Join(dir, ".pdfgen-*"))
	require.NoError(t, err)
	assert.Empty(t, leftovers, "file sementara harus dibersihkan")
}
