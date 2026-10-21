package locale

import (
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

var indonesianMonths = [...]string{
	"Januari", "Februari", "Maret", "April", "Mei", "Juni",
	"Juli", "Agustus", "September", "Oktober", "November", "Desember",
}

// FormatDate renders t in the locale's long date form.
func FormatDate(lang Language, t time.Time) string {
	if lang == Indonesian {
		return fmt.Sprintf("%02d %s %d", t.Day(), indonesianMonths[t.Month()-1], t.Year())
	}
	return t.Format("02 January 2006")
}

// FormatMoney renders value with the locale's currency symbol and separators.
func FormatMoney(lang Language, value decimal.Decimal) string {
	negative := value.IsNegative()
	parts := strings.SplitN(value.Abs().StringFixed(2), ".", 2)
	whole := groupDigits(parts[0], ',')

	if lang == Indonesian {
		result := "Rp" + groupDigits(parts[0], '.')
		if parts[1] != "00" {
			result += "," + parts[1]
		}
		if negative {
			result = "-" + result
		}
		return result
	}

	result := "$" + whole + "." + parts[1]
	if negative {
		result = "-" + result
	}
	return result
}

func groupDigits(value string, separator byte) string {
	var grouped strings.Builder
	for i, digit := range value {
		if i > 0 && (len(value)-i)%3 == 0 {
			grouped.WriteByte(separator)
		}
		grouped.WriteRune(digit)
	}
	return grouped.String()
}
