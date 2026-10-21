package report

import (
	"fmt"
	"strconv"

	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/signature"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
	"github.com/shopspring/decimal"

	"github.com/ardizariat/generate-pdf/internal/data"
	"github.com/ardizariat/generate-pdf/internal/locale"
)

var invoiceWidths = []int{5, 1, 3, 3}

// AddInvoice renders the invoice layout using only the data payload.
func AddInvoice(m core.Maroto, r data.Report) error {
	labels := locale.Get(r.Language)
	addInvoiceHeader(m, r, labels)

	if err := addTableHeader(m, []string{labels.Item, labels.Qty, labels.Price, labels.Subtotal}, invoiceWidths); err != nil {
		return fmt.Errorf("invoice table header: %w", err)
	}
	for _, item := range r.InvoiceItems {
		cells := []string{
			item.Name,
			strconv.Itoa(item.Qty),
			locale.FormatMoney(r.Language, item.Price),
			locale.FormatMoney(r.Language, item.Subtotal),
		}
		if err := addTableRow(m, cells, invoiceWidths); err != nil {
			return fmt.Errorf("invoice table row: %w", err)
		}
	}

	addInvoiceTotals(m, r, labels)
	addInvoiceSignature(m, r.From)
	return nil
}

// addInvoiceHeader adds the seller letterhead, document meta, and both parties.
func addInvoiceHeader(m core.Maroto, r data.Report, labels locale.Labels) {
	m.AddRows(row.New(8).Add(
		col.New(gridSize).Add(text.New(r.From.Company, props.Text{Style: fontstyle.Bold, Size: 12})),
	))
	m.AddRows(row.New(5).Add(
		col.New(gridSize).Add(text.New(r.From.Address+" • "+r.From.Email, props.Text{Size: 8})),
	))
	m.AddRows(row.New(4))

	m.AddRows(row.New(5).Add(
		col.New(6).Add(text.New(labels.InvoiceNo+": "+r.DocID, props.Text{Size: 8})),
		col.New(6).Add(text.New(labels.Date+": "+locale.FormatDate(r.Language, r.Date), props.Text{Size: 8, Align: align.Right})),
	))
	m.AddRows(row.New(4))

	m.AddRows(row.New(5).Add(
		col.New(6).Add(text.New(labels.Seller, props.Text{Style: fontstyle.Bold, Size: 8})),
		col.New(6).Add(text.New(labels.Buyer, props.Text{Style: fontstyle.Bold, Size: 8})),
	))
	addInvoicePartyRow(m, r.From.Company, r.To.Company)
	addInvoicePartyRow(m, r.From.Address, r.To.Address)
	addInvoicePartyRow(m, r.From.Email, r.To.Email)
	m.AddRows(row.New(4))
}

func addInvoicePartyRow(m core.Maroto, left, right string) {
	m.AddRows(row.New(4).Add(
		col.New(6).Add(text.New(left, props.Text{Size: 8})),
		col.New(6).Add(text.New(right, props.Text{Size: 8})),
	))
}

// addInvoiceTotals adds the right-aligned amount summary.
func addInvoiceTotals(m core.Maroto, r data.Report, labels locale.Labels) {
	addInvoiceTotalRow(m, r.Language, labels.Subtotal, r.InvoiceTotals.Subtotal, false)
	addInvoiceTotalRow(m, r.Language, labels.Tax, r.InvoiceTotals.Tax, false)
	addInvoiceTotalRow(m, r.Language, labels.Total, r.InvoiceTotals.Total, true)
	m.AddRows(row.New(4))
}

func addInvoiceTotalRow(m core.Maroto, lang locale.Language, label string, value decimal.Decimal, bold bool) {
	style := fontstyle.Normal
	if bold {
		style = fontstyle.Bold
	}

	m.AddRows(row.New(6).Add(
		col.New(6),
		col.New(3).Add(text.New(label, props.Text{Size: 9, Style: style, Align: align.Right})),
		col.New(3).Add(text.New(locale.FormatMoney(lang, value), props.Text{Size: 9, Style: style, Align: align.Right})),
	))
}

// addInvoiceSignature adds the issuer signature on the right side.
func addInvoiceSignature(m core.Maroto, issuer data.Party) {
	m.AddRows(row.New(18).Add(
		col.New(6),
		signature.NewCol(6, issuer.Name+"\n"+issuer.Position),
	))
}
