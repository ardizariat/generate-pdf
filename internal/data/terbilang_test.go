package data

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTerbilang(t *testing.T) {
	tests := []struct {
		n    int64
		want string
	}{
		{n: 0, want: "nol"},
		{n: 7, want: "tujuh"},
		{n: 11, want: "sebelas"},
		{n: 25, want: "dua puluh lima"},
		{n: 100, want: "seratus"},
		{n: 1000, want: "seribu"},
		{n: 1_500_000, want: "satu juta lima ratus ribu"},
		{n: 2_000_000_000, want: "dua miliar"},
		{n: 1_000_000_000_000, want: "satu triliun"},
		{n: -250, want: "minus dua ratus lima puluh"},
	}

	for _, tc := range tests {
		t.Run(tc.want, func(t *testing.T) {
			assert.Equal(t, tc.want, Terbilang(tc.n))
		})
	}
}

func TestTerbilangTanpaDigitAtauSpasiGanda(t *testing.T) {
	got := Terbilang(12_345_678)
	assert.NotContains(t, got, "  ")
	assert.False(t, strings.ContainsAny(got, "0123456789"), got)
}
