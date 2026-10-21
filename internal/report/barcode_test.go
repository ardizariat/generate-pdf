package report

import (
	"bytes"
	"image"
	_ "image/png"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBarcodePNGProducesValidPNG(t *testing.T) {
	data, err := BarcodePNG("INV-123")
	require.NoError(t, err)
	require.NotEmpty(t, data)
	assert.True(t, bytes.HasPrefix(data, []byte("\x89PNG\r\n\x1a\n")), "missing PNG signature")

	img, _, err := image.Decode(bytes.NewReader(data))
	require.NoError(t, err)
	assert.Equal(t, barcodeWidth, img.Bounds().Dx())
	assert.Equal(t, barcodeHeight, img.Bounds().Dy())
}

func TestBarcodePNGFitsDocumentID(t *testing.T) {
	// A 36-character UUID needs more than 300px, so the width must stay above it.
	_, err := BarcodePNG("0f8b4c1e-9a2d-5b7e-8c3f-1d2e3f4a5b6c")
	require.NoError(t, err)
}

func TestBarcodePNGRejectsInvalidInput(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{name: "empty", value: ""},
		{name: "non-ASCII", value: "é"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := BarcodePNG(tt.value)
			require.Error(t, err)
			assert.Empty(t, data)
		})
	}
}

func TestBarcodePNGIsDeterministic(t *testing.T) {
	first, err := BarcodePNG("INV-123")
	require.NoError(t, err)
	second, err := BarcodePNG("INV-123")
	require.NoError(t, err)
	assert.Equal(t, first, second)

	other, err := BarcodePNG("INV-124")
	require.NoError(t, err)
	assert.NotEqual(t, first, other)
}
