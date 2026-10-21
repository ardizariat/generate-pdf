package report

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ardizariat/generate-pdf/internal/data"
	"github.com/ardizariat/generate-pdf/internal/locale"
)

func TestGeneratePagesExactTarget(t *testing.T) {
	tests := []struct {
		docType string
		pages   int
	}{
		{docType: "invoice", pages: 1},
		{docType: "invoice", pages: 2},
		{docType: "invoice", pages: 3},
		{docType: "report", pages: 1},
		{docType: "report", pages: 2},
		{docType: "report", pages: 3},
		{docType: "delivery-note", pages: 1},
		{docType: "delivery-note", pages: 2},
		{docType: "delivery-note", pages: 3},
		{docType: "letter", pages: 1},
		{docType: "receipt", pages: 1},
		{docType: "certificate", pages: 1},
	}

	for _, tc := range tests {
		t.Run(tc.docType, func(t *testing.T) {
			out := filepath.Join(t.TempDir(), "out.pdf")
			build := func(rows int) (data.Report, error) {
				return data.Generate(tc.docType, rows, 42, locale.English)
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
		return data.Generate("invoice", rows, 42, locale.English)
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
		return data.Generate("invoice", rows, 42, locale.English)
	}

	_, _, _, err := GeneratePages(build, out, 0, true)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "target pages")

	_, statErr := os.Stat(out)
	assert.ErrorIs(t, statErr, os.ErrNotExist)
}

func TestGeneratePagesUnreachableTarget(t *testing.T) {
	out := filepath.Join(t.TempDir(), "out.pdf")
	build := func(rows int) (data.Report, error) {
		return data.Generate("report", rows, 42, locale.English)
	}

	_, _, _, err := GeneratePages(build, out, 1000, true)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "1000")

	_, statErr := os.Stat(out)
	assert.ErrorIs(t, statErr, os.ErrNotExist)
}
