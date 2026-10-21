// Package prompt renders the interactive batch form that fills a job list.
package prompt

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"charm.land/huh/v2"

	"github.com/ardizariat/generate-pdf/internal/config"
)

// typeLabels maps type slugs to the labels shown in the form.
var typeLabels = map[string]string{
	config.TypeInvoice:     "Invoice",
	config.TypeLetter:      "Surat Resmi",
	config.TypeReport:      "Laporan",
	config.TypeReceipt:     "Kuitansi",
	config.TypeDelivery:    "Surat Jalan",
	config.TypeCertificate: "Sertifikat",
}

// IsInteractive reports whether stdin is a terminal, so the CLI can decide
// whether to show the batch forms.
func IsInteractive() bool {
	return isTerminal()
}

// FillBatch asks for the file count and document types, then one page count per
// file. Output names, seed, and title come from the configuration defaults.
func FillBatch(cfg config.Config) ([]config.Job, error) {
	if cfg.NoInput || !isTerminal() {
		return nil, errors.New("prompt batch membutuhkan terminal interaktif")
	}

	count := "1"
	types := []string{cfg.Type}
	fields := make([]huh.Field, 0, 2)
	fields = append(fields, huh.NewInput().
		Title("Jumlah file").
		Value(&count).
		Validate(validateCount))
	if !cfg.Provided("type") {
		fields = append(fields, huh.NewMultiSelect[string]().
			Title("Jenis file (bisa lebih dari satu, dipakai berurutan)").
			Options(typeOptions()...).
			Value(&types).
			Validate(validateTypes))
	}
	if err := huh.NewForm(huh.NewGroup(fields...)).Run(); err != nil {
		return nil, mapAbort(err)
	}

	total, err := strconv.Atoi(count)
	if err != nil {
		return nil, fmt.Errorf("parse jumlah file: %w", err)
	}

	pages := make([]string, total)
	pageFields := make([]huh.Field, 0, total)
	for i := range pages {
		pages[i] = "1"
		pageFields = append(pageFields, huh.NewInput().
			Title(fmt.Sprintf("File %d: jumlah halaman", i+1)).
			Value(&pages[i]).
			Validate(validatePages))
	}
	if err := huh.NewForm(huh.NewGroup(pageFields...)).Run(); err != nil {
		return nil, mapAbort(err)
	}

	pageCounts := make([]int, 0, total)
	for _, value := range pages {
		count, err := strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("parse pages: %w", err)
		}
		pageCounts = append(pageCounts, count)
	}

	return buildJobs(cfg, types, pageCounts)
}

// buildJobs cycles the selected types over the requested page counts. Types
// without a table always get one page.
func buildJobs(base config.Config, types []string, pages []int) ([]config.Job, error) {
	if len(types) == 0 {
		return nil, errors.New("pilih minimal satu jenis dokumen")
	}
	if len(pages) == 0 {
		return nil, errors.New("minimal satu file")
	}

	jobs := make([]config.Job, 0, len(pages))
	for i, pageCount := range pages {
		docType := types[i%len(types)]
		title := typeLabels[docType]
		if base.Provided("title") {
			title = base.Title
		}
		if !config.SupportsMultiPage(docType) {
			pageCount = 1
		}
		jobs = append(jobs, config.Job{
			Type:  docType,
			Pages: pageCount,
			Seed:  base.Seed,
			Out:   fmt.Sprintf("%s-%d-%d.pdf", docType, i+1, base.Seed),
			Title: title,
		})
	}

	if err := config.ValidateJobs(jobs); err != nil {
		return nil, err
	}
	return jobs, nil
}

func typeOptions() []huh.Option[string] {
	types := config.Types()
	options := make([]huh.Option[string], 0, len(types))
	for _, docType := range types {
		options = append(options, huh.NewOption(typeLabels[docType], docType))
	}
	return options
}

func validateCount(value string) error {
	count, err := strconv.Atoi(value)
	if err != nil {
		return errors.New("masukkan angka bulat")
	}
	if count < 1 {
		return errors.New("minimal 1 file")
	}
	return nil
}

func validatePages(value string) error {
	pages, err := strconv.Atoi(value)
	if err != nil {
		return errors.New("masukkan angka bulat")
	}
	if pages < 1 {
		return errors.New("minimal 1 halaman")
	}
	return nil
}

func validateTypes(types []string) error {
	if len(types) == 0 {
		return errors.New("pilih minimal satu jenis")
	}
	return nil
}

// mapAbort turns a user abort into an explicit error instead of letting it
// look like a crash.
func mapAbort(err error) error {
	if errors.Is(err, huh.ErrUserAborted) {
		return fmt.Errorf("input dibatalkan: %w", err)
	}
	return err
}

func isTerminal() bool {
	info, err := os.Stdin.Stat()
	return err == nil && info.Mode()&os.ModeCharDevice != 0
}
