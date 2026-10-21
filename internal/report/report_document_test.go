package report

import (
	"path/filepath"
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ardizariat/generate-pdf/internal/data"
)

func TestReportDocumentLayout(t *testing.T) {
	tests := []struct {
		name      string
		rows      int
		minPages  int
		exactPage int
	}{
		{name: "many rows break pages", rows: 70, minPages: 3},
		{name: "few rows stay on one page", rows: 5, minPages: 1, exactPage: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := data.Generate("laporan", tt.rows, 42)
			require.NoError(t, err)
			require.Len(t, r.ReportRows, tt.rows)

			m, err := newDocument(r)
			require.NoError(t, err)
			require.NoError(t, AddReportDocument(m, r))

			doc, err := m.Generate()
			require.NoError(t, err)

			path := filepath.Join(t.TempDir(), "laporan.pdf")
			require.NoError(t, doc.Save(path))
			require.NoError(t, api.ValidateFile(path, model.NewDefaultConfiguration()))

			pages, err := api.PageCountFile(path)
			require.NoError(t, err)
			assert.GreaterOrEqual(t, pages, tt.minPages)
			if tt.exactPage > 0 {
				assert.Equal(t, tt.exactPage, pages)
			}
			// Manual visual check: 70 rows are expected to show the table header
			// on page two. Pages >= 3 already proves a natural page break without
			// an explicit AddPage.
		})
	}
}
