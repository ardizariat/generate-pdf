package report

import (
	"github.com/johnfercher/maroto/v2/pkg/components/code"
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

// AddCertificate renders the landscape certificate using only the data payload.
// The caller already configured the landscape page.
func AddCertificate(m core.Maroto, r data.Report) error {
	m.AddRows(row.New(6))
	m.AddRows(line.NewRow(1))
	m.AddRows(row.New(6))

	m.AddRows(row.New(10).Add(
		col.New().Add(text.New("SERTIFIKAT", props.Text{
			Style: fontstyle.Bold,
			Size:  18,
			Align: align.Center,
		})),
	))
	m.AddRows(row.New(4))
	m.AddRows(row.New(14).Add(
		col.New().Add(text.New(r.To.Name, props.Text{
			Style: fontstyle.Bold,
			Size:  28,
			Align: align.Center,
		})),
	))
	m.AddRows(row.New(4))
	m.AddRows(text.NewAutoRow(certificateParagraph(r), props.Text{
		Size:  10,
		Align: align.Center,
	}))
	m.AddRows(row.New(8))

	m.AddRows(row.New(6).Add(
		col.New(6).Add(text.New(r.Date.Format("02 Jan 2006"), props.Text{Size: 10})),
		col.New(6),
	))
	m.AddRows(row.New(28).Add(
		signature.NewCol(6, r.From.Name),
		code.NewQrCol(6, "https://example.com/cert/"+r.DocID, props.Rect{Center: true, Percent: 80}),
	))
	m.AddRows(row.New(6))
	m.AddRows(line.NewRow(1))
	return nil
}

// certificateParagraph picks the payload paragraph that states the data is fictitious.
func certificateParagraph(r data.Report) string {
	if len(r.Paragraphs) > 0 {
		return r.Paragraphs[len(r.Paragraphs)-1]
	}
	return "Dokumen ini dihasilkan secara otomatis dan seluruh datanya fiktif."
}
