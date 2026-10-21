package report

import (
	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/signature"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/border"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"

	"github.com/ardizariat/generate-pdf/internal/data"
)

const receiptBoxHeight = 100

// AddReceipt renders the receipt layout using only the data payload.
func AddReceipt(m core.Maroto, r data.Report) error {
	m.AddRows(
		row.New(receiptBoxHeight).
			WithStyle(&props.Cell{BorderType: border.Full, BorderThickness: 0.4}).
			Add(receiptContent(r), receiptSignature(r.From)),
	)
	return nil
}

// receiptContent builds the framed body; the box takes its content height.
func receiptContent(r data.Report) core.Col {
	return col.New(8).Add(
		text.New("KUITANSI", props.Text{Style: fontstyle.Bold, Size: 14, Align: align.Center, Top: 4}),
		text.New("No: "+r.DocID[:8], props.Text{Size: 9, Align: align.Center, Top: 16}),
		text.New("Telah terima dari", props.Text{Style: fontstyle.Bold, Size: 9, Top: 28}),
		text.New(r.From.Company, props.Text{Size: 10, Top: 34}),
		text.New("Uang sejumlah", props.Text{Style: fontstyle.Bold, Size: 9, Top: 46}),
		text.New(formatRupiah(r.Amount), props.Text{Style: fontstyle.Bold, Size: 12, Top: 52}),
		text.New(data.Terbilang(r.Amount.IntPart())+" rupiah", props.Text{Size: 9, Top: 60}),
		text.New("Untuk pembayaran", props.Text{Style: fontstyle.Bold, Size: 9, Top: 74}),
		text.New(r.Purpose, props.Text{Size: 10, Top: 80}),
		text.New(r.From.City+", "+r.Date.Format("02 Jan 2006"), props.Text{Size: 9, Top: 90}),
	)
}

// receiptSignature places the recipient signature at the bottom right, inside
// the same framed row. The signature line renders at the bottom of its cell, so
// the position is a separate line above the frame.
func receiptSignature(issuer data.Party) core.Col {
	return signature.NewCol(4, issuer.Name).Add(
		text.New(issuer.Position, props.Text{Size: 8, Align: align.Center, Top: 92}),
	)
}
