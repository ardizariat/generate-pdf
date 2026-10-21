package report

import (
	"fmt"
	"strconv"

	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"

	"github.com/ardizariat/generate-pdf/internal/data"
)

var (
	reportTitles = []string{"No", "Nama", "Kota", "Email", "Nilai", "Tanggal"}
	reportWidths = []int{1, 3, 2, 3, 2, 2}
)

// reportGridSize gives the report table 13 grid units, matching its column widths.
const reportGridSize = 13

// AddReportDocument renders the report layout using only the data payload.
// Table rows flow onto new pages by themselves; no manual AddPage is used.
func AddReportDocument(m core.Maroto, r data.Report) error {
	addReportHeading(m, r)

	if err := addTableHeader(m, reportTitles, reportWidths); err != nil {
		return fmt.Errorf("report table header: %w", err)
	}
	for _, record := range r.ReportRows {
		cells := []string{
			strconv.Itoa(record.No),
			record.Name,
			record.City,
			record.Email,
			formatRupiah(record.Value),
			record.Date.Format("02 Jan 2006"),
		}
		if err := addTableRow(m, cells, reportWidths); err != nil {
			return fmt.Errorf("report table row: %w", err)
		}
	}
	return nil
}

// addReportHeading adds the title, date, and two summary paragraphs.
func addReportHeading(m core.Maroto, r data.Report) {
	m.AddRows(row.New(9).Add(
		col.New().Add(text.New(r.Title, props.Text{Style: fontstyle.Bold, Size: 14})),
	))
	m.AddRows(row.New(5).Add(
		col.New().Add(text.New("Tanggal: "+r.Date.Format("02 Jan 2006"), props.Text{Size: 9})),
	))
	m.AddRows(row.New(4))

	paragraphs := r.Paragraphs
	if len(paragraphs) > 0 {
		// First paragraph states the document is generated; the last states the
		// data is fictitious.
		addParagraph(m, paragraphs[0])
		if len(paragraphs) > 1 {
			addParagraph(m, paragraphs[len(paragraphs)-1])
		}
	}
	m.AddRows(row.New(4))
}
