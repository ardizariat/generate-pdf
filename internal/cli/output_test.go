package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSingleCreatesOutputDir(t *testing.T) {
	out := filepath.Join(t.TempDir(), "nested", "sub", "invoice.pdf")

	_, err := execute(t, "-t", "invoice", "--rows", "1", "--out", out, "--no-input")
	require.NoError(t, err)

	_, err = os.Stat(out)
	require.NoError(t, err)
}

func TestSingleDefaultLandsInOutputDir(t *testing.T) {
	t.Chdir(t.TempDir())

	_, err := execute(t, "-t", "kuitansi", "--no-input")
	require.NoError(t, err)

	matches, err := filepath.Glob(filepath.Join("output", "kuitansi_*.pdf"))
	require.NoError(t, err)
	require.Len(t, matches, 1)
}

func TestBatchUsesOutDirForDefaultsAndRelativePaths(t *testing.T) {
	outDir := filepath.Join(t.TempDir(), "hasil")
	batch := writeBatch(t, `jobs:
  - type: invoice
  - type: laporan
    out: nested/laporan.pdf
`)

	_, err := execute(t, "--batch", batch, "--out-dir", outDir, "--no-input")
	require.NoError(t, err)

	matches, err := filepath.Glob(filepath.Join(outDir, "invoice_*.pdf"))
	require.NoError(t, err)
	require.Len(t, matches, 1)

	_, err = os.Stat(filepath.Join(outDir, "nested", "laporan.pdf"))
	require.NoError(t, err)
}
