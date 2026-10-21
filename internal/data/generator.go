package data

import (
	"fmt"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/ardizariat/generate-pdf/internal/locale"
)

// Party is a party involved in a document.
type Party struct {
	Name     string
	Company  string
	Address  string
	City     string
	Email    string
	Phone    string
	Position string
}

// InvoiceItem is one line item on an invoice.
type InvoiceItem struct {
	Name     string
	Qty      int
	Price    decimal.Decimal
	Subtotal decimal.Decimal
}

// InvoiceTotals holds the invoice amount summary; Tax is the 11% VAT.
type InvoiceTotals struct {
	Subtotal decimal.Decimal
	Tax      decimal.Decimal
	Total    decimal.Decimal
}

// DeliveryItem is one line item on a delivery note.
type DeliveryItem struct {
	Name string
	Qty  int
	Unit string
	Note string
}

// ReportRow is one row of a report table.
type ReportRow struct {
	No    int
	Name  string
	City  string
	Email string
	Value decimal.Decimal
	Date  time.Time
}

// Report is the full payload needed by the PDF layouts.
type Report struct {
	Type          string
	Title         string
	Language      locale.Language
	DocID         string
	Number        string
	Date          time.Time
	Seed          uint64
	Rows          int
	From          Party
	To            Party
	Purpose       string
	InvoiceItems  []InvoiceItem
	InvoiceTotals InvoiceTotals
	DeliveryItems []DeliveryItem
	ReportRows    []ReportRow
	Amount        decimal.Decimal
	Paragraphs    []string
}

// Generate builds the document payload deterministically from type, row count,
// seed, and language.
func Generate(docType string, rows int, seed uint64, lang locale.Language) (Report, error) {
	title, ok := locale.Get(lang).TypeTitles[docType]
	if !ok {
		return Report{}, fmt.Errorf("unknown document type: %q", docType)
	}
	if rows < 0 {
		rows = 0
	}

	ds := dataFor(lang)
	faker := gofakeit.New(seed)
	date := time.Unix(int64(seed), 0).UTC()

	report := Report{
		Type:       docType,
		Title:      title,
		Language:   lang,
		DocID:      uuid.NewSHA1(uuid.NameSpaceOID, []byte(fmt.Sprintf("%s-%d", docType, seed))).String(),
		Number:     fmt.Sprintf("%03d/GENPDF/%d", faker.IntRange(1, 999), date.Year()),
		Date:       date,
		Seed:       seed,
		Rows:       rows,
		From:       newParty(faker, ds),
		To:         newParty(faker, ds),
		Purpose:    faker.RandomString(ds.purposes),
		Paragraphs: newParagraphs(faker, title, ds),
	}

	switch docType {
	case "invoice":
		report.InvoiceItems, report.InvoiceTotals = newInvoice(faker, rows, ds)
	case "delivery-note":
		report.DeliveryItems = newDeliveryItems(faker, rows, ds)
	case "report":
		report.ReportRows = newReportRows(faker, rows, date, ds)
	case "receipt":
		report.Amount = ds.receipt.fake(faker)
	}

	return report, nil
}

func newParty(faker *gofakeit.Faker, ds dataset) Party {
	name := faker.RandomString(ds.firstNames) + " " + faker.RandomString(ds.lastNames)

	return Party{
		Name:     name,
		Company:  faker.RandomString(ds.companies),
		Address:  ds.address(faker),
		City:     faker.RandomString(ds.cities),
		Email:    strings.ToLower(strings.ReplaceAll(name, " ", ".")) + "@example.com",
		Phone:    faker.Numerify(ds.phonePattern),
		Position: faker.RandomString(ds.positions),
	}
}

func newInvoice(faker *gofakeit.Faker, rows int, ds dataset) ([]InvoiceItem, InvoiceTotals) {
	items := make([]InvoiceItem, 0, rows)
	subtotal := decimal.Zero

	for i := 0; i < rows; i++ {
		qty := faker.IntRange(1, 10)
		price := ds.price.fake(faker)
		line := price.Mul(decimal.NewFromInt(int64(qty)))

		items = append(items, InvoiceItem{
			Name:     faker.RandomString(ds.goods),
			Qty:      qty,
			Price:    price,
			Subtotal: line,
		})
		subtotal = subtotal.Add(line)
	}

	tax := subtotal.Mul(decimal.NewFromInt(11)).Div(decimal.NewFromInt(100)).Round(2)
	return items, InvoiceTotals{Subtotal: subtotal, Tax: tax, Total: subtotal.Add(tax)}
}

func newDeliveryItems(faker *gofakeit.Faker, rows int, ds dataset) []DeliveryItem {
	items := make([]DeliveryItem, 0, rows)
	for i := 0; i < rows; i++ {
		items = append(items, DeliveryItem{
			Name: faker.RandomString(ds.goods),
			Qty:  faker.IntRange(1, 20),
			Unit: faker.RandomString(ds.units),
			Note: faker.RandomString(ds.deliveryNotes),
		})
	}
	return items
}

func newReportRows(faker *gofakeit.Faker, rows int, date time.Time, ds dataset) []ReportRow {
	result := make([]ReportRow, 0, rows)
	for i := 0; i < rows; i++ {
		name := faker.RandomString(ds.firstNames) + " " + faker.RandomString(ds.lastNames)
		result = append(result, ReportRow{
			No:    i + 1,
			Name:  name,
			City:  faker.RandomString(ds.cities),
			Email: strings.ToLower(strings.ReplaceAll(name, " ", ".")) + "@example.com",
			Value: ds.report.fake(faker),
			Date:  date.AddDate(0, 0, faker.IntRange(0, 365)),
		})
	}
	return result
}

func newParagraphs(faker *gofakeit.Faker, title string, ds dataset) []string {
	first := fmt.Sprintf(ds.intro, title)

	return []string{
		first + " " + faker.RandomString(ds.prose) + " " + faker.RandomString(ds.prose),
		faker.RandomString(ds.prose) + " " + faker.RandomString(ds.prose) + " " + faker.RandomString(ds.prose),
		ds.fictional + " " + faker.RandomString(ds.prose) + " " + faker.RandomString(ds.prose),
	}
}
