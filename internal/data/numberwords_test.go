package data

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ardizariat/generate-pdf/internal/locale"
)

func TestNumberWordsEnglish(t *testing.T) {
	tests := []struct {
		n    int64
		want string
	}{
		{0, "zero"},
		{1, "one"},
		{11, "eleven"},
		{20, "twenty"},
		{21, "twenty one"},
		{100, "one hundred"},
		{101, "one hundred one"},
		{999, "nine hundred ninety nine"},
		{1000, "one thousand"},
		{1500, "one thousand five hundred"},
		{1_000_000, "one million"},
		{1_234_567, "one million two hundred thirty four thousand five hundred sixty seven"},
		{-250, "minus two hundred fifty"},
		{1_000_000_000_000, "one trillion"},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.want, NumberWords(locale.English, tt.n), tt.n)
	}
}

func TestNumberWordsIndonesian(t *testing.T) {
	tests := []struct {
		n    int64
		want string
	}{
		{0, "nol"},
		{1, "satu"},
		{11, "sebelas"},
		{20, "dua puluh"},
		{100, "seratus"},
		{101, "seratus satu"},
		{1000, "seribu"},
		{1500, "seribu lima ratus"},
		{1_000_000, "satu juta"},
		{-250, "minus dua ratus lima puluh"},
		{1_000_000_000_000, "satu triliun"},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.want, NumberWords(locale.Indonesian, tt.n), tt.n)
	}
}
