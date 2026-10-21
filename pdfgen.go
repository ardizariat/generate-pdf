// Package pdfgen renders the same dummy office documents as the pdfgen CLI.
//
// The API is intentionally small: fill Options, then call Render for any
// io.Writer or Generate for a file path. Type slugs, row counts, and seeds keep
// output deterministic, and Language selects English (default) or Indonesian.
package pdfgen

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/ardizariat/generate-pdf/internal/config"
	"github.com/ardizariat/generate-pdf/internal/data"
	"github.com/ardizariat/generate-pdf/internal/locale"
	"github.com/ardizariat/generate-pdf/internal/report"
)

// Sentinel errors returned by Render and Generate.
var (
	ErrInvalidOptions  = errors.New("invalid options")
	ErrUnsupportedType = errors.New("unsupported document type")
	ErrOutputExists    = errors.New("output file already exists")
)

// Options controls one document. Zero values mean: Type=invoice,
// Language=en, and Title=the localized default for the type. Rows is the exact
// table row count; use -1 for the per-type default and 0 for no table rows.
// Seed is deterministic, including 0.
type Options struct {
	Type     string
	Title    string
	Language string
	Rows     int
	Seed     uint64
	Verify   bool
	Force    bool
}

// Result reports the generated document's page count and byte size.
type Result struct {
	Pages int
	Bytes int64
}

// Render writes the PDF to w and returns its page count and size.
func Render(w io.Writer, opts Options) (Result, error) {
	if w == nil {
		return Result{}, fmt.Errorf("%w: writer is nil", ErrInvalidOptions)
	}
	r, err := build(opts)
	if err != nil {
		return Result{}, err
	}

	tmp, err := os.CreateTemp("", "pdfgen-render-*.pdf")
	if err != nil {
		return Result{}, fmt.Errorf("create temporary file: %w", err)
	}
	path := tmp.Name()
	if err := tmp.Close(); err != nil {
		_ = os.Remove(path)
		return Result{}, fmt.Errorf("close temporary file: %w", err)
	}
	defer func() { _ = os.Remove(path) }()

	pages, size, err := report.Generate(r, path, opts.Verify)
	if err != nil {
		return Result{}, err
	}

	f, err := os.Open(path)
	if err != nil {
		return Result{}, fmt.Errorf("open rendered pdf: %w", err)
	}
	defer func() { _ = f.Close() }()

	if _, err := io.Copy(w, f); err != nil {
		return Result{Pages: pages, Bytes: size}, fmt.Errorf("write pdf: %w", err)
	}
	return Result{Pages: pages, Bytes: size}, nil
}

// Generate writes the PDF to path and returns its page count and size. Missing
// parent directories are created. An existing file is only replaced with
// Force=true.
func Generate(path string, opts Options) (Result, error) {
	if path == "" {
		return Result{}, fmt.Errorf("%w: path is empty", ErrInvalidOptions)
	}
	r, err := build(opts)
	if err != nil {
		return Result{}, err
	}

	if _, err := os.Stat(path); err == nil && !opts.Force {
		return Result{}, fmt.Errorf("%w: %s", ErrOutputExists, path)
	} else if err != nil && !errors.Is(err, os.ErrNotExist) {
		return Result{}, fmt.Errorf("check %s: %w", path, err)
	}

	dir := filepath.Dir(path)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return Result{}, fmt.Errorf("create directory %s: %w", dir, err)
		}
	}

	pages, size, err := report.Generate(r, path, opts.Verify)
	if err != nil {
		return Result{}, err
	}
	return Result{Pages: pages, Bytes: size}, nil
}

func build(opts Options) (data.Report, error) {
	docType := opts.Type
	if docType == "" {
		docType = config.TypeInvoice
	}
	if !supportedType(docType) {
		return data.Report{}, fmt.Errorf("%w: %q", ErrUnsupportedType, docType)
	}

	lang, err := locale.Parse(opts.Language)
	if err != nil {
		return data.Report{}, fmt.Errorf("%w: %v", ErrInvalidOptions, err)
	}
	if opts.Rows < -1 {
		return data.Report{}, fmt.Errorf("%w: rows must be >= -1", ErrInvalidOptions)
	}
	rows := opts.Rows
	if rows == -1 {
		rows = config.DefaultRows(docType)
	}

	r, err := data.Generate(docType, rows, opts.Seed, lang)
	if err != nil {
		return data.Report{}, err
	}
	if opts.Title != "" {
		r.Title = opts.Title
	}
	return r, nil
}

func supportedType(docType string) bool {
	for _, candidate := range config.Types() {
		if candidate == docType {
			return true
		}
	}
	return false
}
