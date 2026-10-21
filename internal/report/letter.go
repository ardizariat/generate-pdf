package report

import (
	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/line"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/signature"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"

	"github.com/ardizariat/generate-pdf/internal/data"
)

// AddLetter renders the formal letter layout using only the data payload.
func AddLetter(m core.Maroto, r data.Report) error {
	addLetterHeader(m, r.From)
	addLetterMeta(m, r)
	addLetterRecipient(m, r.To)

	addParagraph(m, "Dengan hormat,")
	for _, paragraph := range r.Paragraphs {
		addParagraph(m, paragraph)
		m.AddRows(row.New(2))
	}

	addLetterClosing(m, r.From)
	return nil
}

// addLetterHeader adds the centered letterhead and its rule.
func addLetterHeader(m core.Maroto, from data.Party) {
	m.AddRows(row.New(10).Add(
		col.New(gridSize).Add(text.New(from.Company, props.Text{
			Style: fontstyle.Bold,
			Size:  14,
			Align: align.Center,
		})),
	))
	m.AddRows(row.New(5).Add(
		col.New(gridSize).Add(text.New(from.Address+" • "+from.Email, props.Text{
			Size:  8,
			Align: align.Center,
		})),
	))
	m.AddRows(row.New(3).Add(line.NewCol(gridSize)))
}

// addLetterMeta adds the letter reference lines.
func addLetterMeta(m core.Maroto, r data.Report) {
	metaRows := []string{
		"Nomor: " + r.Number,
		"Tanggal: " + r.Date.Format("02 Jan 2006"),
		"Perihal: " + r.Purpose,
	}
	for _, meta := range metaRows {
		m.AddRows(row.New(5).Add(
			col.New(gridSize).Add(text.New(meta, props.Text{Size: 9})),
		))
	}
	m.AddRows(row.New(4))
}

// addLetterRecipient adds the addressee block.
func addLetterRecipient(m core.Maroto, to data.Party) {
	m.AddRows(row.New(5).Add(
		col.New(gridSize).Add(text.New("Kepada Yth.", props.Text{Size: 9})),
	))
	m.AddRows(row.New(5).Add(
		col.New(gridSize).Add(text.New(to.Name, props.Text{Size: 9, Style: fontstyle.Bold})),
	))
	m.AddRows(row.New(5).Add(
		col.New(gridSize).Add(text.New(to.Position+" - "+to.Company, props.Text{Size: 9})),
	))
	m.AddRows(row.New(5).Add(
		col.New(gridSize).Add(text.New(to.City, props.Text{Size: 9})),
	))
	m.AddRows(row.New(4))
}

// addLetterClosing adds the right-aligned signature block.
func addLetterClosing(m core.Maroto, from data.Party) {
	m.AddRows(row.New(6).Add(
		col.New(6),
		col.New(6).Add(text.New("Hormat kami,", props.Text{Size: 9})),
	))
	m.AddRows(row.New(18).Add(
		col.New(6),
		signature.NewCol(6, from.Name+"\n"+from.Position),
	))
}
