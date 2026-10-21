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
)

func TestReceiptLayout(t *testing.T) {
	r, err := data.Generate("kuitansi", 0, 42)
	require.NoError(t, err)

	m, err := newDocument(r)
	require.NoError(t, err)
	require.NoError(t, AddReceipt(m, r))

	doc, err := m.Generate()
	require.NoError(t, err)

	path := filepath.Join(t.TempDir(), "kuitansi.pdf")
	require.NoError(t, doc.Save(path))
	require.NoError(t, api.ValidateFile(path, model.NewDefaultConfiguration()))

	pages, err := api.PageCountFile(path)
	require.NoError(t, err)
	assert.Equal(t, 1, pages)
}

func TestReceiptAmountWords(t *testing.T) {
	r, err := data.Generate("kuitansi", 0, 42)
	require.NoError(t, err)

	assert.True(t, r.Amount.GreaterThanOrEqual(decimal.NewFromInt(100_000)))
	assert.True(t, r.Amount.LessThanOrEqual(decimal.NewFromInt(50_000_000)))
	assert.True(t, r.Amount.Mod(decimal.NewFromInt(1_000)).IsZero())

	words := data.Terbilang(r.Amount.IntPart())
	assert.NotEmpty(t, words)
	assert.NotRegexp(t, `[0-9]`, words)
}
