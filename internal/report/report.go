package report

import (
	"bytes"
	"fmt"
	"os"

	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/phpdave11/gofpdf"

	"github.com/ardizariat/generate-pdf/internal/data"
)

// verifyPDF is Verify; tests replace it to exercise a failed verification.
var verifyPDF = Verify

// Generate renders r, saves the PDF to out, and reports its page count and
// size. Verification always runs to obtain the page count; with verify set, a
// failed verification also fails the generate.
func Generate(r data.Report, out string, verify bool) (pages int, size int64, err error) {
	// gofpdf walks font maps in random order and stamps ModDate with the
	// wall clock, so identical payloads would otherwise differ. maroto does
	// not expose either knob.
	gofpdf.SetDefaultCatalogSort(true)
	gofpdf.SetDefaultModificationDate(r.Date)

	m, err := newDocument(r)
	if err != nil {
		return 0, 0, fmt.Errorf("buat dokumen: %w", err)
	}
	if err := addContent(m, r); err != nil {
		return 0, 0, err
	}

	doc, err := m.Generate()
	if err != nil {
		return 0, 0, fmt.Errorf("render %s: %w", r.Type, err)
	}
	if err := doc.Save(out); err != nil {
		return 0, 0, fmt.Errorf("simpan %s: %w", out, err)
	}

	info, err := os.Stat(out)
	if err != nil {
		return 0, 0, fmt.Errorf("stat %s: %w", out, err)
	}
	size = info.Size()

	pages, err = verifyPDF(out)
	if err != nil && verify {
		return 0, size, fmt.Errorf("verifikasi %s: %w", out, err)
	}
	return pages, size, nil
}

// addContent dispatches on the document type; there is no fallback layout.
func addContent(m core.Maroto, r data.Report) error {
	switch r.Type {
	case "invoice":
		return AddInvoice(m, r)
	case "surat-resmi":
		return AddLetter(m, r)
	case "laporan":
		return AddReportDocument(m, r)
	case "kuitansi":
		return AddReceipt(m, r)
	case "surat-jalan":
		return AddDelivery(m, r)
	case "sertifikat":
		return AddCertificate(m, r)
	default:
		return fmt.Errorf("jenis dokumen tidak didukung: %q", r.Type)
	}
}

// Verify validates the PDF at path with pdfcpu and returns its page count.
func Verify(path string) (pages int, err error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return 0, fmt.Errorf("baca %s: %w", path, err)
	}

	// pdfcpu repairs a truncated tail, so the EOF marker is checked here.
	if !bytes.HasSuffix(bytes.TrimSpace(content), []byte("%%EOF")) {
		return 0, fmt.Errorf("validasi %s: penanda %%EOF tidak ditemukan", path)
	}

	if err := api.Validate(bytes.NewReader(content), model.NewDefaultConfiguration()); err != nil {
		return 0, fmt.Errorf("validasi %s: %w", path, err)
	}

	pages, err = api.PageCountFile(path)
	if err != nil {
		return 0, fmt.Errorf("hitung halaman %s: %w", path, err)
	}
	return pages, nil
}
