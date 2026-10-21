package report

import (
	"path/filepath"
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ardizariat/generate-pdf/internal/data"
	"github.com/ardizariat/generate-pdf/internal/locale"
)

func TestReceiptLayout(t *testing.T) {
	r, err := data.Generate("receipt", 0, 42, locale.English)
	require.NoError(t, err)

	m, err := newDocument(r)
	require.NoError(t, err)
	require.NoError(t, AddReceipt(m, r))

	doc, err := m.Generate()
	require.NoError(t, err)

	path := filepath.Join(t.TempDir(), "receipt.pdf")
	require.NoError(t, doc.Save(path))
	require.NoError(t, api.ValidateFile(path, model.NewDefaultConfiguration()))

	pages, err := api.PageCountFile(path)
	require.NoError(t, err)
	assert.Equal(t, 1, pages)
}

func TestReceiptAmountWords(t *testing.T) {
	r, err := data.Generate("receipt", 0, 42, locale.English)
	require.NoError(t, err)

	assert.True(t, r.Amount.GreaterThanOrEqual(decimal.NewFromInt(100)))
	assert.True(t, r.Amount.LessThanOrEqual(decimal.NewFromInt(50_000)))
	assert.True(t, r.Amount.Equal(r.Amount.Round(2)))

	words := data.NumberWords(locale.English, r.Amount.IntPart())
	assert.NotEmpty(t, words)
	assert.NotRegexp(t, `[0-9]`, words)
}
