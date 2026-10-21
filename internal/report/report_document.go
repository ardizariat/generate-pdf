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
	"github.com/ardizariat/generate-pdf/internal/locale"
)

var reportWidths = []int{1, 3, 2, 3, 2, 2}

// reportGridSize gives the report table 13 grid units, matching its column widths.
const reportGridSize = 13

// AddReportDocument renders the report layout using only the data payload.
// Table rows flow onto new pages by themselves; no manual AddPage is used.
func AddReportDocument(m core.Maroto, r data.Report) error {
	labels := locale.Get(r.Language)
	addReportHeading(m, r, labels)

	rows := make([][]string, 0, len(r.ReportRows))
	for _, record := range r.ReportRows {
		rows = append(rows, []string{
			strconv.Itoa(record.No),
			record.Name,
			record.City,
			record.Email,
			locale.FormatMoney(r.Language, record.Value),
			locale.FormatDate(r.Language, record.Date),
		})
	}
	titles := []string{labels.No, labels.Name, labels.City, labels.Email, labels.Value, labels.Date}
	if err := addTable(m, titles, reportWidths, rows); err != nil {
		return fmt.Errorf("report table: %w", err)
	}
	return nil
}

// addReportHeading adds the title, date, and two summary paragraphs.
func addReportHeading(m core.Maroto, r data.Report, labels locale.Labels) {
	m.AddRows(row.New(9).Add(
		col.New().Add(text.New(r.Title, props.Text{Style: fontstyle.Bold, Size: 14})),
	))
	m.AddRows(row.New(5).Add(
		col.New().Add(text.New(labels.Date+": "+locale.FormatDate(r.Language, r.Date), props.Text{Size: 9})),
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
