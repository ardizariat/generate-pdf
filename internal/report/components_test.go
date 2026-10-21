package report

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ardizariat/generate-pdf/internal/data"
)

func testReport() data.Report {
	return data.Report{
		Type:  "invoice",
		Title: "Invoice",
		DocID: "INV-123",
		Date:  time.Unix(42, 0).UTC(),
		Seed:  42,
	}
}

func TestComponentsTableValidation(t *testing.T) {
	tests := []struct {
		name   string
		cells  []string
		widths []int
	}{
		{name: "widths below grid", cells: []string{"a", "b"}, widths: []int{5, 6}},
		{name: "widths above grid", cells: []string{"a", "b"}, widths: []int{7, 6}},
		{name: "cells and widths differ", cells: []string{"a", "b", "c"}, widths: []int{6, 6}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Error(t, addTableHeader(nil, tt.cells, tt.widths))
			assert.Error(t, addTableRow(nil, tt.cells, tt.widths))
		})
	}
}

func TestComponentsEmptyDocument(t *testing.T) {
	m, err := newDocument(testReport())
	require.NoError(t, err)

	doc, err := m.Generate()
	require.NoError(t, err)

	path := filepath.Join(t.TempDir(), "empty.pdf")
	require.NoError(t, doc.Save(path))
	require.NoError(t, api.ValidateFile(path, model.NewDefaultConfiguration()))

	pages, err := api.PageCountFile(path)
	require.NoError(t, err)
	assert.Equal(t, 1, pages)
}

func TestComponentsHeaderBarcodeError(t *testing.T) {
	r := testReport()
	r.DocID = "é"

	_, err := newDocument(r)
	assert.Error(t, err)
}

func TestComponentsTableRenders(t *testing.T) {
	m, err := newDocument(testReport())
	require.NoError(t, err)

	require.NoError(t, addTableHeader(m, []string{"Barang", "Qty", "Harga", "Subtotal"}, []int{5, 2, 2, 3}))
	require.NoError(t, addTableRow(m, []string{"Kertas", "2", "Rp10.000", "Rp20.000"}, []int{5, 2, 2, 3}))
	addParagraph(m, "Dokumen ini dihasilkan otomatis dan datanya fiktif.")

	doc, err := m.Generate()
	require.NoError(t, err)

	path := filepath.Join(t.TempDir(), "table.pdf")
	require.NoError(t, doc.Save(path))
	assert.NoError(t, api.ValidateFile(path, model.NewDefaultConfiguration()))
}
