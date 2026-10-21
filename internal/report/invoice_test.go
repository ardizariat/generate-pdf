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

func TestInvoiceLayout(t *testing.T) {
	r, err := data.Generate("invoice", 8, 42)
	require.NoError(t, err)

	m, err := newDocument(r)
	require.NoError(t, err)
	require.NoError(t, AddInvoice(m, r))

	doc, err := m.Generate()
	require.NoError(t, err)

	path := filepath.Join(t.TempDir(), "invoice.pdf")
	require.NoError(t, doc.Save(path))
	require.NoError(t, api.ValidateFile(path, model.NewDefaultConfiguration()))

	pages, err := api.PageCountFile(path)
	require.NoError(t, err)
	assert.Equal(t, 1, pages)
}

func TestInvoiceTotals(t *testing.T) {
	r, err := data.Generate("invoice", 8, 42)
	require.NoError(t, err)

	assert.True(t, r.InvoiceTotals.Total.Equal(r.InvoiceTotals.Subtotal.Add(r.InvoiceTotals.Tax)),
		"total %s harus sama dengan subtotal %s + PPN %s",
		r.InvoiceTotals.Total, r.InvoiceTotals.Subtotal, r.InvoiceTotals.Tax)
	assert.Len(t, r.InvoiceItems, 8)
}

func TestInvoiceFormatRupiah(t *testing.T) {
	tests := []struct {
		name  string
		value decimal.Decimal
		want  string
	}{
		{name: "whole thousands", value: decimal.NewFromInt(1234567), want: "Rp1.234.567"},
		{name: "fractional tax", value: decimal.RequireFromString("135802.37"), want: "Rp135.802,37"},
		{name: "small", value: decimal.NewFromInt(500), want: "Rp500"},
		{name: "zero", value: decimal.Zero, want: "Rp0"},
		{name: "negative", value: decimal.NewFromInt(-1500), want: "-Rp1.500"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, formatRupiah(tt.value))
		})
	}
}
