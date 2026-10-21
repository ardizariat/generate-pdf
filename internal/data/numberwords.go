package data

import (
	"strconv"
	"strings"

	"github.com/ardizariat/generate-pdf/internal/locale"
)

var englishOnes = [...]string{
	"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine",
	"ten", "eleven", "twelve", "thirteen", "fourteen", "fifteen", "sixteen",
	"seventeen", "eighteen", "nineteen",
}

var englishTens = [...]string{"", "", "twenty", "thirty", "forty", "fifty", "sixty", "seventy", "eighty", "ninety"}

var englishScales = [...]struct {
	value int64
	name  string
}{
	{1_000_000_000_000, "trillion"},
	{1_000_000_000, "billion"},
	{1_000_000, "million"},
	{1_000, "thousand"},
}

var indonesianNumbers = [...]string{
	"", "satu", "dua", "tiga", "empat", "lima",
	"enam", "tujuh", "delapan", "sembilan", "sepuluh", "sebelas",
}

// NumberWords spells n in the requested language without the currency name.
// Supported range: up to trillions.
func NumberWords(lang locale.Language, n int64) string {
	if lang == locale.Indonesian {
		return indonesianWords(n)
	}
	return englishWords(n)
}

func englishWords(n int64) string {
	if n == 0 {
		return "zero"
	}
	if n < 0 {
		// shortcut: negating math.MinInt64 overflows; document amounts never get that large.
		return "minus " + englishWords(-n)
	}
	if n >= 1_000_000_000_000_000 {
		return strconv.FormatInt(n, 10)
	}

	parts := make([]string, 0, 4)
	for _, scale := range englishScales {
		if n < scale.value {
			continue
		}
		parts = append(parts, englishWords(n/scale.value)+" "+scale.name)
		n %= scale.value
	}
	if n > 0 {
		parts = append(parts, englishUnderThousand(n))
	}
	return strings.Join(parts, " ")
}

func englishUnderThousand(n int64) string {
	switch {
	case n < 20:
		return englishOnes[n]
	case n < 100:
		if n%10 == 0 {
			return englishTens[n/10]
		}
		return englishTens[n/10] + " " + englishOnes[n%10]
	default:
		result := englishOnes[n/100] + " hundred"
		if n%100 != 0 {
			result += " " + englishUnderThousand(n%100)
		}
		return result
	}
}

func indonesianWords(n int64) string {
	if n == 0 {
		return "nol"
	}
	if n < 0 {
		// shortcut: negating math.MinInt64 overflows; document amounts never get that large.
		return "minus " + indonesianWords(-n)
	}
	return strings.Join(strings.Fields(indonesianSpell(n)), " ")
}

func indonesianSpell(n int64) string {
	switch {
	case n < 12:
		return indonesianNumbers[n]
	case n < 20:
		return indonesianSpell(n-10) + " belas"
	case n < 100:
		return indonesianSpell(n/10) + " puluh " + indonesianSpell(n%10)
	case n < 200:
		return "seratus " + indonesianSpell(n-100)
	case n < 1000:
		return indonesianSpell(n/100) + " ratus " + indonesianSpell(n%100)
	case n < 2000:
		return "seribu " + indonesianSpell(n-1000)
	case n < 1_000_000:
		return indonesianSpell(n/1000) + " ribu " + indonesianSpell(n%1000)
	case n < 1_000_000_000:
		return indonesianSpell(n/1_000_000) + " juta " + indonesianSpell(n%1_000_000)
	case n < 1_000_000_000_000:
		return indonesianSpell(n/1_000_000_000) + " miliar " + indonesianSpell(n%1_000_000_000)
	case n < 1_000_000_000_000_000:
		return indonesianSpell(n/1_000_000_000_000) + " triliun " + indonesianSpell(n%1_000_000_000_000)
	}
	return strconv.FormatInt(n, 10)
}
