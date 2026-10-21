// Package report renders document layouts as PDF.
package report

import (
	"bytes"
	"fmt"
	"image/png"
	"unicode"

	scale "github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
)

const (
	barcodeWidth  = 500
	barcodeHeight = 100
)

// BarcodePNG renders value as a Code 128 barcode scaled to 500x100 and
// encoded as PNG. Code 128 only supports ASCII, so empty or non-ASCII
// input is rejected.
func BarcodePNG(value string) ([]byte, error) {
	if value == "" {
		return nil, fmt.Errorf("barcode: value is empty")
	}
	for _, r := range value {
		if r > unicode.MaxASCII {
			return nil, fmt.Errorf("barcode: value %q contains non-ASCII rune %q", value, r)
		}
	}

	code, err := code128.Encode(value)
	if err != nil {
		return nil, fmt.Errorf("encode code128: %w", err)
	}

	scaled, err := scale.Scale(code, barcodeWidth, barcodeHeight)
	if err != nil {
		return nil, fmt.Errorf("scale barcode: %w", err)
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, scaled); err != nil {
		return nil, fmt.Errorf("encode png: %w", err)
	}
	return buf.Bytes(), nil
}
