package data

import "strings"

var spellNumbers = [...]string{
	"", "satu", "dua", "tiga", "empat", "lima",
	"enam", "tujuh", "delapan", "sembilan", "sepuluh", "sebelas",
}

// Terbilang mengubah angka menjadi kata bahasa Indonesia tanpa kata "rupiah".
// Rentang yang didukung: sampai triliun.
func Terbilang(n int64) string {
	if n == 0 {
		return "nol"
	}
	if n < 0 {
		// shortcut: math.MinInt64 overflow saat negasi; nominal dokumen tidak sebesar itu.
		return "minus " + Terbilang(-n)
	}
	return strings.Join(strings.Fields(spell(n)), " ")
}

func spell(n int64) string {
	switch {
	case n < 12:
		return spellNumbers[n]
	case n < 20:
		return spell(n-10) + " belas"
	case n < 100:
		return spell(n/10) + " puluh " + spell(n%10)
	case n < 200:
		return "seratus " + spell(n-100)
	case n < 1000:
		return spell(n/100) + " ratus " + spell(n%100)
	case n < 2000:
		return "seribu " + spell(n-1000)
	case n < 1_000_000:
		return spell(n/1000) + " ribu " + spell(n%1000)
	case n < 1_000_000_000:
		return spell(n/1_000_000) + " juta " + spell(n%1_000_000)
	case n < 1_000_000_000_000:
		return spell(n/1_000_000_000) + " miliar " + spell(n%1_000_000_000)
	case n < 1_000_000_000_000_000:
		return spell(n/1_000_000_000_000) + " triliun " + spell(n%1_000_000_000_000)
	}
	return ""
}
