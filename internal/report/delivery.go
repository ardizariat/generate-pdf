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

	"github.com/ardizariat/generate-pdf/internal/data"
	"github.com/ardizariat/generate-pdf/internal/locale"
)

var deliveryWidths = []int{5, 1, 2, 4}

// AddDelivery renders the delivery note layout using only the data payload.
// Table rows flow onto new pages by themselves; no manual AddPage is used.
func AddDelivery(m core.Maroto, r data.Report) error {
	labels := locale.Get(r.Language)
	addDeliveryHeader(m, r, labels)
	addDeliveryRecipient(m, r.To, labels)

	rows := make([][]string, 0, len(r.DeliveryItems))
	for _, item := range r.DeliveryItems {
		rows = append(rows, []string{
			item.Name,
			strconv.Itoa(item.Qty),
			item.Unit,
			item.Note,
		})
	}
	titles := []string{labels.Item, labels.Qty, labels.Unit, labels.Note}
	if err := addTable(m, titles, deliveryWidths, rows); err != nil {
		return fmt.Errorf("delivery table: %w", err)
	}

	addDeliverySignatures(m, r, labels)
	return nil
}

// addDeliveryHeader adds the sender letterhead and document meta. Heights stay
// compact so 25 table rows still fit on one page.
func addDeliveryHeader(m core.Maroto, r data.Report, labels locale.Labels) {
	m.AddRows(row.New(6).Add(
		col.New().Add(text.New(r.From.Company, props.Text{Style: fontstyle.Bold, Size: 12})),
	))
	m.AddRows(row.New(4).Add(
		col.New().Add(text.New(r.From.Address+" • "+r.From.Phone, props.Text{Size: 8})),
	))
	m.AddRows(row.New(4).Add(
		col.New(6).Add(text.New(labels.DeliveryNo+": DN-"+r.DocID[:6], props.Text{Size: 8})),
		col.New(6).Add(text.New(labels.Date+": "+locale.FormatDate(r.Language, r.Date), props.Text{Size: 8, Align: align.Right})),
	))
}

// addDeliveryRecipient adds the addressee block.
func addDeliveryRecipient(m core.Maroto, to data.Party, labels locale.Labels) {
	m.AddRows(row.New(4).Add(
		col.New(2).Add(text.New(labels.To, props.Text{Size: 9})),
		col.New(10).Add(text.New(to.Company, props.Text{Size: 9, Style: fontstyle.Bold})),
	))
	m.AddRows(row.New(4).Add(
		col.New().Add(text.New(to.Address+", "+to.City, props.Text{Size: 9})),
	))
}

// addDeliverySignatures adds the sender and recipient signature blocks.
func addDeliverySignatures(m core.Maroto, r data.Report, labels locale.Labels) {
	m.AddRows(row.New(17).Add(
		signature.NewCol(6, r.From.Name).Add(text.New(labels.Sender, props.Text{Size: 9})),
		signature.NewCol(6, r.To.Name).Add(text.New(labels.Receiver, props.Text{Size: 9})),
	))
}
