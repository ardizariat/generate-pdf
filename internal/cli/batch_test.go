package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ardizariat/generate-pdf/internal/report"
)

func writeBatch(t *testing.T, content string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "jobs.yaml")
	require.NoError(t, os.WriteFile(path, []byte(content), 0o644))
	return path
}

func TestBatchGeneratesEveryJob(t *testing.T) {
	dir := t.TempDir()
	invoice := filepath.Join(dir, "invoice.pdf")
	laporan := filepath.Join(dir, "report.pdf")
	delivery := filepath.Join(dir, "delivery-note.pdf")

	batch := writeBatch(t, fmt.Sprintf(`jobs:
  - type: invoice
    pages: 1
    seed: 1
    out: %s
  - type: report
    pages: 2
    seed: 2
    out: %s
  - type: delivery-note
    pages: 3
    seed: 3
    out: %s
`, invoice, laporan, delivery))

	out, err := execute(t, "--batch", batch, "--no-input", "--verify")
	require.NoError(t, err)
	assert.Contains(t, out, "summary: 3 succeeded, 0 failed")

	for path, wantPages := range map[string]int{invoice: 1, laporan: 2, delivery: 3} {
		info, err := os.Stat(path)
		require.NoError(t, err, path)
		assert.Positive(t, info.Size(), path)

		pages, err := report.Verify(path)
		require.NoError(t, err, path)
		assert.Equal(t, wantPages, pages, path)
	}
}

func TestBatchContinuesAfterFailure(t *testing.T) {
	dir := t.TempDir()
	existing := filepath.Join(dir, "existing.pdf")
	laporan := filepath.Join(dir, "report.pdf")
	require.NoError(t, os.WriteFile(existing, []byte("%PDF-lama"), 0o644))

	batch := writeBatch(t, fmt.Sprintf(`jobs:
  - type: invoice
    out: %s
  - type: report
    pages: 1
    out: %s
`, existing, laporan))

	out, err := execute(t, "--batch", batch, "--no-input")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "1 of 2 files failed")
	assert.Contains(t, out, "failed  "+existing)
	assert.Contains(t, out, "summary: 1 succeeded, 1 failed")

	unchanged, err := os.ReadFile(existing)
	require.NoError(t, err)
	assert.Equal(t, "%PDF-lama", string(unchanged))

	_, err = os.Stat(laporan)
	require.NoError(t, err, "job lain harus tetap jalan")
}

func TestBatchForceOverwrites(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "invoice.pdf")
	require.NoError(t, os.WriteFile(path, []byte("%PDF-lama"), 0o644))

	batch := writeBatch(t, fmt.Sprintf("jobs:\n  - type: invoice\n    pages: 1\n    out: %s\n", path))

	_, err := execute(t, "--batch", batch, "--no-input")
	require.Error(t, err)

	_, err = execute(t, "--batch", batch, "--no-input", "--force")
	require.NoError(t, err)

	pages, err := report.Verify(path)
	require.NoError(t, err)
	assert.Equal(t, 1, pages)
}

func TestBatchInvalidManifestWritesNothing(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "surat.pdf")

	batch := writeBatch(t, fmt.Sprintf("jobs:\n  - type: letter\n    pages: 2\n    out: %s\n", path))

	_, err := execute(t, "--batch", batch, "--no-input")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "letter")

	_, statErr := os.Stat(path)
	assert.ErrorIs(t, statErr, os.ErrNotExist)
}

func TestBatchTitleChangesTheDocument(t *testing.T) {
	dir := t.TempDir()
	first := filepath.Join(dir, "first.pdf")
	second := filepath.Join(dir, "second.pdf")

	batch := writeBatch(t, fmt.Sprintf(`jobs:
  - type: invoice
    seed: 42
    out: %s
    title: Invoice Pertama
  - type: invoice
    seed: 42
    out: %s
    title: Invoice Kedua
`, first, second))

	_, err := execute(t, "--batch", batch, "--no-input")
	require.NoError(t, err)

	a, err := os.ReadFile(first)
	require.NoError(t, err)
	b, err := os.ReadFile(second)
	require.NoError(t, err)
	assert.NotEqual(t, a, b, "judul harus memengaruhi isi dokumen")
}
