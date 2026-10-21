package locale

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	for _, value := range []string{"", "en"} {
		lang, err := Parse(value)
		require.NoError(t, err)
		assert.Equal(t, English, lang)
	}
	lang, err := Parse("id")
	require.NoError(t, err)
	assert.Equal(t, Indonesian, lang)

	_, err = Parse("fr")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported language")
}

func TestOrDefault(t *testing.T) {
	assert.Equal(t, English, OrDefault(""))
	assert.Equal(t, English, OrDefault("fr"))
	assert.Equal(t, Indonesian, OrDefault("id"))
}

func TestFormatDate(t *testing.T) {
	value := time.Date(2026, time.January, 2, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, "02 January 2026", FormatDate(English, value))
	assert.Equal(t, "02 Januari 2026", FormatDate(Indonesian, value))
}

func TestFormatMoney(t *testing.T) {
	tests := []struct {
		name string
		lang Language
		in   string
		want string
	}{
		{name: "english whole", lang: English, in: "1234.56", want: "$1,234.56"},
		{name: "english small", lang: English, in: "500", want: "$500.00"},
		{name: "english negative", lang: English, in: "-1500", want: "-$1,500.00"},
		{name: "indonesian whole", lang: Indonesian, in: "1234567", want: "Rp1.234.567"},
		{name: "indonesian fraction", lang: Indonesian, in: "135802.37", want: "Rp135.802,37"},
		{name: "indonesian negative", lang: Indonesian, in: "-1500", want: "-Rp1.500"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, FormatMoney(tt.lang, decimal.RequireFromString(tt.in)))
		})
	}
}

func TestLabels(t *testing.T) {
	en := Get(English)
	id := Get(Indonesian)
	assert.Equal(t, "Invoice", en.TypeTitles["invoice"])
	assert.Equal(t, "Kuitansi", id.TypeTitles["receipt"])
	assert.Equal(t, "Tax (11%)", en.Tax)
	assert.Equal(t, "PPN 11%", id.Tax)
	assert.Equal(t, "dollars", en.CurrencyName)
	assert.Equal(t, "rupiah", id.CurrencyName)
}
