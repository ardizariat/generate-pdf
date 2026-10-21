package data

import (
	"github.com/ardizariat/generate-pdf/internal/locale"
	"strings"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateDeterministic(t *testing.T) {
	first, err := Generate("invoice", 8, 42, locale.English)
	require.NoError(t, err)
	second, err := Generate("invoice", 8, 42, locale.English)
	require.NoError(t, err)
	assert.Equal(t, first, second)

	other, err := Generate("invoice", 8, 43, locale.English)
	require.NoError(t, err)
	assert.NotEqual(t, first.DocID, other.DocID)
	assert.NotEqual(t, first.InvoiceItems, other.InvoiceItems)
}

func TestGenerateUnknownType(t *testing.T) {
	_, err := Generate("faktur", 3, 42, locale.English)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "faktur")
}

func TestGenerateDateFromSeed(t *testing.T) {
	report, err := Generate("letter", 0, 42, locale.English)
	require.NoError(t, err)
	assert.Equal(t, time.Unix(42, 0).UTC(), report.Date)
}

func TestGenerateRowCount(t *testing.T) {
	invoice, err := Generate("invoice", 8, 42, locale.English)
	require.NoError(t, err)
	assert.Len(t, invoice.InvoiceItems, 8)
	assert.Equal(t, 8, invoice.Rows)

	delivery, err := Generate("delivery-note", 25, 42, locale.English)
	require.NoError(t, err)
	assert.Len(t, delivery.DeliveryItems, 25)

	reportDoc, err := Generate("report", 70, 42, locale.English)
	require.NoError(t, err)
	assert.Len(t, reportDoc.ReportRows, 70)
	assert.Equal(t, 1, reportDoc.ReportRows[0].No)
	assert.Equal(t, 70, reportDoc.ReportRows[69].No)

	for _, typ := range []string{"letter", "receipt", "certificate"} {
		report, err := Generate(typ, 5, 42, locale.English)
		require.NoError(t, err)
		assert.Empty(t, report.InvoiceItems, typ)
		assert.Empty(t, report.DeliveryItems, typ)
		assert.Empty(t, report.ReportRows, typ)
	}
}

func TestGenerateInvoiceTotals(t *testing.T) {
	report, err := Generate("invoice", 8, 42, locale.English)
	require.NoError(t, err)

	subtotal := decimal.Zero
	for _, item := range report.InvoiceItems {
		assert.True(t, item.Subtotal.Equal(item.Price.Mul(decimal.NewFromInt(int64(item.Qty)))))
		subtotal = subtotal.Add(item.Subtotal)
	}

	wantTax := subtotal.Mul(decimal.NewFromInt(11)).Div(decimal.NewFromInt(100)).Round(2)
	assert.True(t, report.InvoiceTotals.Subtotal.Equal(subtotal))
	assert.True(t, report.InvoiceTotals.Tax.Equal(wantTax))
	assert.True(t, report.InvoiceTotals.Total.Equal(subtotal.Add(wantTax)))
}

func TestGenerateReceiptAmount(t *testing.T) {
	report, err := Generate("receipt", 0, 42, locale.English)
	require.NoError(t, err)

	assert.True(t, report.Amount.GreaterThanOrEqual(decimal.NewFromInt(100)))
	assert.True(t, report.Amount.LessThanOrEqual(decimal.NewFromInt(50_000)))
	assert.True(t, report.Amount.Equal(report.Amount.Round(2)))
	assert.NotEmpty(t, NumberWords(locale.English, report.Amount.IntPart()))
}

func TestGenerateParagraphs(t *testing.T) {
	for _, typ := range []string{"invoice", "letter", "report", "receipt", "delivery-note", "certificate"} {
		report, err := Generate(typ, 0, 42, locale.English)
		require.NoError(t, err)

		require.GreaterOrEqual(t, len(report.Paragraphs), 3, typ)
		assert.Contains(t, strings.Join(report.Paragraphs, " "), "fictional data", typ)
		for _, paragraph := range report.Paragraphs {
			assert.NotEmpty(t, paragraph, typ)
		}
	}
}

func TestGenerateCompleteParty(t *testing.T) {
	report, err := Generate("letter", 0, 42, locale.English)
	require.NoError(t, err)

	for _, party := range []Party{report.From, report.To} {
		assert.NotEmpty(t, party.Name)
		assert.NotEmpty(t, party.Company)
		assert.NotEmpty(t, party.Address)
		assert.NotEmpty(t, party.City)
		assert.Contains(t, party.Email, "@example.com")
		assert.NotEmpty(t, party.Phone)
		assert.NotEmpty(t, party.Position)
	}
}
