package pdfgen_test

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ardizariat/generate-pdf"
)

func TestRenderDefaultsAndDeterminism(t *testing.T) {
	var first, second bytes.Buffer
	opts := pdfgen.Options{Seed: 42, Rows: 2}

	firstResult, err := pdfgen.Render(&first, opts)
	require.NoError(t, err)
	secondResult, err := pdfgen.Render(&second, opts)
	require.NoError(t, err)

	assert.Equal(t, 1, firstResult.Pages)
	assert.Positive(t, firstResult.Bytes)
	assert.Equal(t, firstResult, secondResult)
	assert.True(t, bytes.HasPrefix(first.Bytes(), []byte("%PDF-")))
	assert.Equal(t, first.Bytes(), second.Bytes())
}

func TestRenderIndonesianDiffers(t *testing.T) {
	var english, indonesian bytes.Buffer
	opts := pdfgen.Options{Type: "receipt", Seed: 42}

	_, err := pdfgen.Render(&english, opts)
	require.NoError(t, err)
	opts.Language = "id"
	_, err = pdfgen.Render(&indonesian, opts)
	require.NoError(t, err)

	assert.True(t, bytes.HasPrefix(indonesian.Bytes(), []byte("%PDF-")))
	assert.NotEqual(t, english.Bytes(), indonesian.Bytes())
}

func TestGenerateWritesFileAndHonorsForce(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "invoice.pdf")

	result, err := pdfgen.Generate(path, pdfgen.Options{Seed: 42, Rows: 1})
	require.NoError(t, err)
	assert.Positive(t, result.Pages)
	assert.Positive(t, result.Bytes)

	content, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.True(t, bytes.HasPrefix(content, []byte("%PDF-")))

	_, err = pdfgen.Generate(path, pdfgen.Options{Seed: 42, Rows: 1})
	require.Error(t, err)
	assert.True(t, errors.Is(err, pdfgen.ErrOutputExists))

	_, err = pdfgen.Generate(path, pdfgen.Options{Seed: 43, Rows: 1, Force: true})
	require.NoError(t, err)
}

func TestRenderRejectsInvalidOptions(t *testing.T) {
	tests := []struct {
		name string
		opts pdfgen.Options
		want error
	}{
		{name: "unknown type", opts: pdfgen.Options{Type: "faktur"}, want: pdfgen.ErrUnsupportedType},
		{name: "unknown language", opts: pdfgen.Options{Language: "fr"}, want: pdfgen.ErrInvalidOptions},
		{name: "negative rows", opts: pdfgen.Options{Rows: -2}, want: pdfgen.ErrInvalidOptions},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := pdfgen.Render(&bytes.Buffer{}, tt.opts)
			require.Error(t, err)
			assert.True(t, errors.Is(err, tt.want))
		})
	}

	_, err := pdfgen.Render(nil, pdfgen.Options{})
	require.Error(t, err)
	assert.True(t, errors.Is(err, pdfgen.ErrInvalidOptions))

	_, err = pdfgen.Generate("", pdfgen.Options{})
	require.Error(t, err)
	assert.True(t, errors.Is(err, pdfgen.ErrInvalidOptions))
}
