package report

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ardizariat/generate-pdf/internal/data"
	"github.com/ardizariat/generate-pdf/internal/locale"
)

// TestCrossTypeGenerate exercises dispatch, writing, and verification for every
// supported type in one table, including the page expectations from the spec.
func TestCrossTypeGenerate(t *testing.T) {
	tests := []struct {
		docType  string
		rows     int
		minPages int
		maxPages int // 0 = no upper bound
	}{
		{docType: "invoice", rows: 8, minPages: 1, maxPages: 1},
		{docType: "letter", rows: 0, minPages: 1, maxPages: 1},
		{docType: "report", rows: 70, minPages: 3},
		{docType: "receipt", rows: 0, minPages: 1, maxPages: 1},
		{docType: "delivery-note", rows: 25, minPages: 1, maxPages: 1},
		{docType: "certificate", rows: 0, minPages: 1, maxPages: 1},
	}

	for _, tt := range tests {
		t.Run(tt.docType, func(t *testing.T) {
			r, err := data.Generate(tt.docType, tt.rows, 42, locale.English)
			require.NoError(t, err)

			path := filepath.Join(t.TempDir(), tt.docType+".pdf")
			pages, size, err := Generate(r, path, true)
			require.NoError(t, err)
			assert.Positive(t, size)
			assert.GreaterOrEqual(t, pages, tt.minPages)
			if tt.maxPages > 0 {
				assert.LessOrEqual(t, pages, tt.maxPages)
			}

			verified, err := Verify(path)
			require.NoError(t, err)
			assert.Equal(t, pages, verified)
		})
	}
}
