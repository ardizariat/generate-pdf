package report

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

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
// failed verification also fails the generate. The file is staged next to out
// and renamed only after rendering (and requested verification) succeeds, so a
// failure never leaves a partial file at out.
func Generate(r data.Report, out string, verify bool) (pages int, size int64, err error) {
	staged, err := stagePath(out)
	if err != nil {
		return 0, 0, fmt.Errorf("create temporary file for %s: %w", out, err)
	}
	defer func() { _ = os.Remove(staged) }()

	pages, size, err = generate(r, staged, out, verify)
	if err != nil {
		return 0, 0, err
	}
	if err := commit(staged, out); err != nil {
		return 0, 0, fmt.Errorf("save %s: %w", out, err)
	}
	return pages, size, nil
}

// generate renders and verifies the PDF at path. display is the user-facing
// path used in errors; path may be a staging file.
func generate(r data.Report, path, display string, verify bool) (pages int, size int64, err error) {
	// gofpdf walks font maps in random order and stamps ModDate with the
	// wall clock, so identical payloads would otherwise differ. maroto does
	// not expose either knob.
	gofpdf.SetDefaultCatalogSort(true)
	gofpdf.SetDefaultModificationDate(r.Date)

	m, err := newDocument(r)
	if err != nil {
		return 0, 0, fmt.Errorf("build document: %w", err)
	}
	if err := addContent(m, r); err != nil {
		return 0, 0, err
	}

	doc, err := m.Generate()
	if err != nil {
		return 0, 0, fmt.Errorf("render %s: %w", r.Type, err)
	}
	if err := doc.Save(path); err != nil {
		return 0, 0, fmt.Errorf("save %s: %w", display, err)
	}

	info, err := os.Stat(path)
	if err != nil {
		return 0, 0, fmt.Errorf("stat %s: %w", display, err)
	}
	size = info.Size()

	pages, err = verifyPDF(path)
	if err != nil && verify {
		return 0, size, fmt.Errorf("verify %s: %w", display, err)
	}
	return pages, size, nil
}

// stagePath creates an empty temporary file in out's directory so the final
// rename stays on the same filesystem.
func stagePath(out string) (string, error) {
	f, err := os.CreateTemp(filepath.Dir(out), ".pdfgen-*.pdf")
	if err != nil {
		return "", err
	}
	name := f.Name()
	if err := f.Close(); err != nil {
		_ = os.Remove(name)
		return "", err
	}
	return name, nil
}

// commit gives the staged file the permissions a regular output would have,
// then renames it over out. The rename is atomic on Unix; Go's Windows rename
// replaces the destination too, so --force keeps working there.
func commit(staged, out string) error {
	mode := os.FileMode(0o644)
	if info, err := os.Stat(out); err == nil {
		mode = info.Mode().Perm()
	}
	if err := os.Chmod(staged, mode); err != nil {
		return err
	}
	return os.Rename(staged, out)
}

// addContent dispatches on the document type; there is no fallback layout.
func addContent(m core.Maroto, r data.Report) error {
	switch r.Type {
	case "invoice":
		return AddInvoice(m, r)
	case "letter":
		return AddLetter(m, r)
	case "report":
		return AddReportDocument(m, r)
	case "receipt":
		return AddReceipt(m, r)
	case "delivery-note":
		return AddDelivery(m, r)
	case "certificate":
		return AddCertificate(m, r)
	default:
		return fmt.Errorf("unsupported document type: %q", r.Type)
	}
}

// Verify validates the PDF at path with pdfcpu and returns its page count.
func Verify(path string) (pages int, err error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return 0, fmt.Errorf("read %s: %w", path, err)
	}

	// pdfcpu repairs a truncated tail, so the EOF marker is checked here.
	if !bytes.HasSuffix(bytes.TrimSpace(content), []byte("%%EOF")) {
		return 0, fmt.Errorf("validate %s: %%EOF marker not found", path)
	}

	if err := api.Validate(bytes.NewReader(content), model.NewDefaultConfiguration()); err != nil {
		return 0, fmt.Errorf("validate %s: %w", path, err)
	}

	pages, err = api.PageCountFile(path)
	if err != nil {
		return 0, fmt.Errorf("count pages in %s: %w", path, err)
	}
	return pages, nil
}
