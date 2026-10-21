package report

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ardizariat/generate-pdf/internal/data"
)

func TestGeneratePagesExactTarget(t *testing.T) {
	tests := []struct {
		docType string
		pages   int
	}{
		{docType: "invoice", pages: 1},
		{docType: "invoice", pages: 2},
		{docType: "invoice", pages: 3},
		{docType: "laporan", pages: 1},
		{docType: "laporan", pages: 2},
		{docType: "laporan", pages: 3},
		{docType: "surat-jalan", pages: 1},
		{docType: "surat-jalan", pages: 2},
		{docType: "surat-jalan", pages: 3},
		{docType: "surat-resmi", pages: 1},
		{docType: "kuitansi", pages: 1},
		{docType: "sertifikat", pages: 1},
	}

	for _, tc := range tests {
		t.Run(tc.docType, func(t *testing.T) {
			out := filepath.Join(t.TempDir(), "out.pdf")
			build := func(rows int) (data.Report, error) {
				return data.Generate(tc.docType, rows, 42)
			}

			rows, pages, size, err := GeneratePages(build, out, tc.pages, true)
			require.NoError(t, err)
			assert.Positive(t, rows)
			assert.Equal(t, tc.pages, pages)
			assert.Positive(t, size)

			verified, err := Verify(out)
			require.NoError(t, err)
			assert.Equal(t, tc.pages, verified)
		})
	}
}

func TestGeneratePagesIsDeterministic(t *testing.T) {
	dir := t.TempDir()
	build := func(rows int) (data.Report, error) {
		return data.Generate("invoice", rows, 42)
	}

	first := filepath.Join(dir, "first.pdf")
	second := filepath.Join(dir, "second.pdf")
	_, _, _, err := GeneratePages(build, first, 2, true)
	require.NoError(t, err)
	_, _, _, err = GeneratePages(build, second, 2, true)
	require.NoError(t, err)

	a, err := os.ReadFile(first)
	require.NoError(t, err)
	b, err := os.ReadFile(second)
	require.NoError(t, err)
	assert.Equal(t, a, b)
}

func TestGeneratePagesRejectsInvalidTarget(t *testing.T) {
	out := filepath.Join(t.TempDir(), "out.pdf")
	build := func(rows int) (data.Report, error) {
		return data.Generate("invoice", rows, 42)
	}

	_, _, _, err := GeneratePages(build, out, 0, true)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "target halaman")

	_, statErr := os.Stat(out)
	assert.ErrorIs(t, statErr, os.ErrNotExist)
}

func TestGeneratePagesUnreachableTarget(t *testing.T) {
	out := filepath.Join(t.TempDir(), "out.pdf")
	build := func(rows int) (data.Report, error) {
		return data.Generate("laporan", rows, 42)
	}

	_, _, _, err := GeneratePages(build, out, 1000, true)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "1000")

	_, statErr := os.Stat(out)
	assert.ErrorIs(t, statErr, os.ErrNotExist)
}
