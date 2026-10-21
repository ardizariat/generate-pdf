// Package prompt renders the interactive batch form that fills a job list.
package prompt

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"charm.land/huh/v2"

	"github.com/ardizariat/generate-pdf/internal/config"
	"github.com/ardizariat/generate-pdf/internal/locale"
)

// IsInteractive reports whether stdin is a terminal, so the CLI can decide
// whether to show the batch form.
func IsInteractive() bool {
	return isTerminal()
}

// FillBatch asks for the file count, document types, and one page count for
// all files, in a single form. Output names, seed, and title come from the
// configuration defaults.
func FillBatch(cfg config.Config) ([]config.Job, error) {
	if cfg.NoInput || !isTerminal() {
		return nil, errors.New("batch prompt requires an interactive terminal")
	}

	lang := locale.OrDefault(cfg.Language)
	labels := locale.Get(lang)
	count := "1"
	pages := "1"
	types := []string{cfg.Type}
	fields := make([]huh.Field, 0, 3)
	fields = append(fields, huh.NewInput().
		Title(labels.PromptFileCount).
		Value(&count).
		Validate(func(value string) error { return validateCount(value, labels) }))
	if !cfg.Provided("type") {
		fields = append(fields, huh.NewMultiSelect[string]().
			Title(labels.PromptTypes).
			Options(typeOptions(lang)...).
			Value(&types).
			Validate(func(types []string) error { return validateTypes(types, labels) }))
	}
	fields = append(fields, huh.NewInput().
		Title(labels.PromptPages).
		Value(&pages).
		Validate(func(value string) error { return validatePages(value, labels) }))

	if err := huh.NewForm(huh.NewGroup(fields...)).Run(); err != nil {
		return nil, mapAbort(err)
	}

	total, err := strconv.Atoi(count)
	if err != nil {
		return nil, fmt.Errorf("parse file count: %w", err)
	}
	pageCount, err := strconv.Atoi(pages)
	if err != nil {
		return nil, fmt.Errorf("parse page count: %w", err)
	}

	return buildJobs(cfg, types, total, pageCount)
}

// buildJobs cycles the selected types over the file count and applies the same
// page count to every file. Types without a table always get one page.
func buildJobs(base config.Config, types []string, files int, pages int) ([]config.Job, error) {
	if len(types) == 0 {
		return nil, errors.New("select at least one document type")
	}
	if files < 1 {
		return nil, errors.New("at least one file is required")
	}

	lang := locale.OrDefault(base.Language)
	jobs := make([]config.Job, 0, files)
	for i := 0; i < files; i++ {
		docType := types[i%len(types)]
		title := locale.Get(lang).TypeTitles[docType]
		if base.Provided("title") {
			title = base.Title
		}
		pageCount := pages
		if !config.SupportsMultiPage(docType) {
			pageCount = 1
		}
		name, err := config.DefaultOutputName(docType, time.Now())
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, config.Job{
			Type:  docType,
			Pages: pageCount,
			Seed:  base.Seed,
			Out:   config.OutputPath(base.OutDir, name),
			Title: title,
		})
	}

	if err := config.ValidateJobs(jobs); err != nil {
		return nil, err
	}
	return jobs, nil
}

func typeOptions(lang locale.Language) []huh.Option[string] {
	labels := locale.Get(lang)
	types := config.Types()
	options := make([]huh.Option[string], 0, len(types))
	for _, docType := range types {
		options = append(options, huh.NewOption(labels.TypeTitles[docType], docType))
	}
	return options
}

func validateCount(value string, labels locale.Labels) error {
	count, err := strconv.Atoi(value)
	if err != nil {
		return errors.New(labels.PromptInteger)
	}
	if count < 1 {
		return errors.New(labels.PromptMinOneFile)
	}
	return nil
}

func validatePages(value string, labels locale.Labels) error {
	pages, err := strconv.Atoi(value)
	if err != nil {
		return errors.New(labels.PromptInteger)
	}
	if pages < 1 {
		return errors.New(labels.PromptMinOnePage)
	}
	return nil
}

func validateTypes(types []string, labels locale.Labels) error {
	if len(types) == 0 {
		return errors.New(labels.PromptMinOneType)
	}
	return nil
}

// mapAbort turns a user abort into an explicit error instead of letting it
// look like a crash.
func mapAbort(err error) error {
	if errors.Is(err, huh.ErrUserAborted) {
		return fmt.Errorf("input aborted: %w", err)
	}
	return err
}

func isTerminal() bool {
	info, err := os.Stdin.Stat()
	return err == nil && info.Mode()&os.ModeCharDevice != 0
}
