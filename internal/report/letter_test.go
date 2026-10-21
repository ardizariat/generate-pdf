package report

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ardizariat/generate-pdf/internal/data"
)

func TestLetterLayout(t *testing.T) {
	r, err := data.Generate("surat-resmi", 0, 42)
	require.NoError(t, err)

	m, err := newDocument(r)
	require.NoError(t, err)
	require.NoError(t, AddLetter(m, r))

	doc, err := m.Generate()
	require.NoError(t, err)

	path := filepath.Join(t.TempDir(), "surat-resmi.pdf")
	require.NoError(t, doc.Save(path))
	require.NoError(t, api.ValidateFile(path, model.NewDefaultConfiguration()))

	pages, err := api.PageCountFile(path)
	require.NoError(t, err)
	assert.Equal(t, 1, pages)
}

func TestLetterParagraphs(t *testing.T) {
	r, err := data.Generate("surat-resmi", 0, 42)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(r.Paragraphs), 3)

	foundFictional := false
	for _, paragraph := range r.Paragraphs {
		if strings.Contains(paragraph, "data fiktif") {
			foundFictional = true
			break
		}
	}
	assert.True(t, foundFictional, "salah satu paragraf harus memuat frasa data fiktif")
}
