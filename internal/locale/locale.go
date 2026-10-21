// Package locale holds the runtime language selection and the localized
// strings and formatting used by the CLI and the PDF layouts.
package locale

import "fmt"

// Language is a supported output language.
type Language string

// Supported languages. English is the default.
const (
	English    Language = "en"
	Indonesian Language = "id"
)

// Parse returns the language for s. An empty string means English.
func Parse(s string) (Language, error) {
	switch Language(s) {
	case "", English:
		return English, nil
	case Indonesian:
		return Indonesian, nil
	default:
		return "", fmt.Errorf("unsupported language %q (want en or id)", s)
	}
}

// OrDefault returns the parsed language or English for an empty/invalid value.
// It is used by display-only paths; validation paths must call Parse.
func OrDefault(s string) Language {
	lang, err := Parse(s)
	if err != nil {
		return English
	}
	return lang
}

// Labels is the localized text used by the CLI and PDF layouts. Fields that
// contain percent signs are printf formats and are used with fmt.Fprintf.
type Labels struct {
	TypeTitles       map[string]string
	TypeDescriptions map[string]string

	InvoiceNo string
	Date      string
	Seller    string
	Buyer     string
	Item      string
	Qty       string
	Price     string
	Subtotal  string
	Tax       string
	Total     string

	Ref       string
	Subject   string
	Attn      string
	Dear      string
	Sincerely string

	PagePattern string

	No    string
	Name  string
	City  string
	Email string
	Value string

	Receipt      string
	ReceiptNo    string
	ReceivedFrom string
	Amount       string
	InWords      string
	ForPayment   string
	ReceivedBy   string
	CurrencyName string

	DeliveryNote string
	DeliveryNo   string
	To           string
	Unit         string
	Note         string
	Sender       string
	Receiver     string

	Certificate         string
	CertificateFallback string

	PromptFileCount  string
	PromptTypes      string
	PromptPages      string
	PromptInteger    string
	PromptMinOneFile string
	PromptMinOnePage string
	PromptMinOneType string
	ProgressRows     string
	ProgressPages    string
	BatchFailed      string
	BatchOK          string
	BatchSummary     string
}

// Get returns the labels for lang. Unknown languages fall back to English.
func Get(lang Language) Labels {
	if lang == Indonesian {
		return indonesian
	}
	return english
}
