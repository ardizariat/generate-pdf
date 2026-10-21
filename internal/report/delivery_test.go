package report

import (
	"path/filepath"
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ardizariat/generate-pdf/internal/data"
	"github.com/ardizariat/generate-pdf/internal/locale"
)

func TestDeliveryLayout(t *testing.T) {
	tests := []struct {
		name     string
		rows     int
		minPages int
	}{
		{name: "one page", rows: 25, minPages: 1},
		{name: "flows to more pages", rows: 200, minPages: 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := data.Generate("delivery-note", tt.rows, 42, locale.English)
			require.NoError(t, err)
			require.Len(t, r.DeliveryItems, tt.rows)

			m, err := newDocument(r)
			require.NoError(t, err)
			require.NoError(t, AddDelivery(m, r))

			doc, err := m.Generate()
			require.NoError(t, err)

			path := filepath.Join(t.TempDir(), "delivery-note.pdf")
			require.NoError(t, doc.Save(path))
			require.NoError(t, api.ValidateFile(path, model.NewDefaultConfiguration()))

			pages, err := api.PageCountFile(path)
			require.NoError(t, err)
			assert.GreaterOrEqual(t, pages, tt.minPages)
			if tt.rows == 25 {
				assert.Equal(t, 1, pages)
			}
		})
	}
}
