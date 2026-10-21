package report

import (
	"fmt"
	"strings"

	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/code"
	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/image"
	"github.com/johnfercher/maroto/v2/pkg/components/line"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/border"
	"github.com/johnfercher/maroto/v2/pkg/consts/extension"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/consts/orientation"
	"github.com/johnfercher/maroto/v2/pkg/consts/pagesize"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
	"github.com/shopspring/decimal"

	"github.com/ardizariat/generate-pdf/internal/data"
)

const (
	gridSize = 12

	headerTitleHeight = 7
	headerInfoHeight  = 10
	headerLineHeight  = 1

	tableRowHeight     = 8
	tableFontSize      = 9
	documentIDFontSize = 7
	footerFontSize     = 7
	pageNumberFontSize = 7
)

var tableHeaderColor = props.Color{Red: 240, Green: 240, Blue: 240}

// newDocument builds the maroto instance with page config, header, and footer.
// It is the only place document configuration is created.
func newDocument(r data.Report) (core.Maroto, error) {
	pageOrientation := orientation.Vertical
	if r.Type == "sertifikat" {
		pageOrientation = orientation.Horizontal
	}

	cfg := config.NewBuilder().
		WithPageSize(pagesize.A4).
		WithOrientation(pageOrientation).
		WithTopMargin(10).
		WithLeftMargin(10).
		WithRightMargin(10).
		WithBottomMargin(15).
		WithCreationDate(r.Date).
		WithTitle("generate-pdf", true).
		WithAuthor("generate-pdf", true).
		WithCreator("generate-pdf", true).
		WithSubject(r.Type, true).
		WithPageNumber(props.PageNumber{
			Pattern: "Halaman {current} dari {total}",
			Place:   props.Bottom,
			Size:    pageNumberFontSize,
		}).
		Build()

	m := maroto.New(cfg)
	if err := addHeader(m, r); err != nil {
		return nil, fmt.Errorf("add header: %w", err)
	}
	if err := addFooter(m); err != nil {
		return nil, fmt.Errorf("add footer: %w", err)
	}
	return m, nil
}

// addHeader registers the repeating header: title, QR, barcode, document ID, and rule.
func addHeader(m core.Maroto, r data.Report) error {
	barcode, err := BarcodePNG(r.DocID)
	if err != nil {
		return fmt.Errorf("header barcode: %w", err)
	}

	qrPayload := fmt.Sprintf("https://example.com/report?doc=%s&seed=%d&type=%s", r.DocID, r.Seed, r.Type)

	titleRow := row.New(headerTitleHeight).Add(
		col.New(gridSize).Add(text.New(r.Title, props.Text{Style: fontstyle.Bold, Size: 12})),
	)
	infoRow := row.New(headerInfoHeight).Add(
		code.NewQrCol(2, qrPayload, props.Rect{Center: true, Percent: 100}),
		image.NewFromBytesCol(3, barcode, extension.Png, props.Rect{Center: true, Percent: 100}),
		col.New(7).Add(text.New(r.DocID, props.Text{
			Size:  documentIDFontSize,
			Align: align.Right,
			Top:   3,
		})),
	)
	lineRow := row.New(headerLineHeight).Add(line.NewCol(gridSize))

	return m.RegisterHeader(titleRow, infoRow, lineRow)
}

// addFooter registers the repeating footer text. Page numbers come from the
// page config, so they are not repeated here.
func addFooter(m core.Maroto) error {
	footerRow := row.New(5).Add(
		col.New(gridSize).Add(text.New("generate-pdf", props.Text{Size: footerFontSize})),
	)
	return m.RegisterFooter(footerRow)
}

// addTableHeader adds one bordered header row with a light background.
func addTableHeader(m core.Maroto, titles []string, widths []int) error {
	if err := validateTable(titles, widths); err != nil {
		return err
	}

	cols := make([]core.Col, 0, len(widths))
	for i, width := range widths {
		cols = append(cols, col.New(width).
			WithStyle(&props.Cell{
				BackgroundColor: &tableHeaderColor,
				BorderType:      border.Full,
				BorderThickness: 0.1,
			}).
			Add(text.New(titles[i], props.Text{
				Style: fontstyle.Bold,
				Size:  tableFontSize,
				Align: align.Center,
			})))
	}
	m.AddRows(row.New(tableRowHeight).Add(cols...))
	return nil
}

// addTableRow adds one bordered data row.
func addTableRow(m core.Maroto, cells []string, widths []int) error {
	if err := validateTable(cells, widths); err != nil {
		return err
	}

	cols := make([]core.Col, 0, len(widths))
	for i, width := range widths {
		cols = append(cols, col.New(width).
			WithStyle(&props.Cell{
				BorderType:      border.Full,
				BorderThickness: 0.1,
			}).
			Add(text.New(cells[i], props.Text{Size: tableFontSize})))
	}
	m.AddRows(row.New(tableRowHeight).Add(cols...))
	return nil
}

// addParagraph adds a wrapped paragraph.
func addParagraph(m core.Maroto, value string) {
	m.AddRows(text.NewAutoRow(value, props.Text{Size: 9}))
}

// formatRupiah renders a decimal as Indonesian rupiah, for example Rp1.234.567,50.
func formatRupiah(value decimal.Decimal) string {
	negative := value.IsNegative()
	parts := strings.SplitN(value.Abs().StringFixed(2), ".", 2)
	whole := parts[0]

	var grouped strings.Builder
	for i, digit := range whole {
		if i > 0 && (len(whole)-i)%3 == 0 {
			grouped.WriteByte('.')
		}
		grouped.WriteRune(digit)
	}

	result := "Rp" + grouped.String()
	if parts[1] != "00" {
		result += "," + parts[1]
	}
	if negative {
		result = "-" + result
	}
	return result
}

// validateTable checks that each cell has a width and that the widths fill the
// 12-column maroto grid.
func validateTable(cells []string, widths []int) error {
	if len(cells) != len(widths) {
		return fmt.Errorf("table: %d cells for %d widths", len(cells), len(widths))
	}
	sum := 0
	for _, width := range widths {
		sum += width
	}
	if sum != gridSize {
		return fmt.Errorf("table: widths sum to %d, want %d", sum, gridSize)
	}
	return nil
}
