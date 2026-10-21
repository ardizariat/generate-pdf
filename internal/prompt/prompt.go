// Package prompt renders the interactive batch forms that fill a job list.
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

// jobAnswers holds one file's raw form answers before defaults are resolved.
type jobAnswers struct {
	Type  string
	Pages int
	Seed  uint64
	Out   string
}

// IsInteractive reports whether stdin is a terminal, so the CLI can decide
// whether to show the batch forms.
func IsInteractive() bool {
	return isTerminal()
}

// FillBatch asks for the number of files, one form per file, then a final
// confirmation. It must run on an interactive terminal.
func FillBatch(cfg config.Config) ([]config.Job, error) {
	if cfg.NoInput || !isTerminal() {
		return nil, errors.New("prompt batch membutuhkan terminal interaktif")
	}

	count := "1"
	if err := huh.NewForm(huh.NewGroup(
		huh.NewInput().Title("Jumlah file").Value(&count).Validate(validateCount),
	)).Run(); err != nil {
		return nil, mapAbort(err)
	}
	total, err := strconv.Atoi(count)
	if err != nil {
		return nil, fmt.Errorf("parse jumlah file: %w", err)
	}

	answers := make([]jobAnswers, 0, total)
	for i := 0; i < total; i++ {
		answer, err := askJob(cfg, i+1)
		if err != nil {
			return nil, err
		}
		answers = append(answers, answer)
	}

	jobs, err := buildJobs(cfg, answers)
	if err != nil {
		return nil, err
	}

	confirmed := true
	if err := huh.NewForm(huh.NewGroup(
		huh.NewConfirm().
			Title(fmt.Sprintf("Generate %d file?", len(jobs))).
			Description(formatJobs(jobs)).
			Affirmative("Ya").
			Negative("Batal").
			Value(&confirmed),
	)).Run(); err != nil {
		return nil, mapAbort(err)
	}
	if !confirmed {
		return nil, errors.New("input dibatalkan")
	}
	return jobs, nil
}

// askJob renders the form for one file. Values already provided globally are
// defaults and are not asked again; pages and output are always asked.
func askJob(cfg config.Config, index int) (jobAnswers, error) {
	answer := jobAnswers{Type: cfg.Type, Pages: 1, Seed: cfg.Seed}
	pages := strconv.Itoa(answer.Pages)
	seed := strconv.FormatUint(answer.Seed, 10)
	out := ""

	fields := make([]huh.Field, 0, 4)
	if !cfg.Provided("type") {
		fields = append(fields, huh.NewSelect[string]().
			Title(fmt.Sprintf("File %d: jenis dokumen", index)).
			Options(typeOptions()...).
			Value(&answer.Type))
	}
	fields = append(fields, huh.NewInput().
		Title(fmt.Sprintf("File %d: jumlah halaman", index)).
		Value(&pages).
		Validate(validatePages))
	if !cfg.Provided("seed") {
		fields = append(fields, huh.NewInput().
			Title(fmt.Sprintf("File %d: seed", index)).
			Value(&seed).
			Validate(validateSeed))
	}
	fields = append(fields, huh.NewInput().
		Title(fmt.Sprintf("File %d: path output (kosong = default)", index)).
		Value(&out))

	if err := huh.NewForm(huh.NewGroup(fields...)).Run(); err != nil {
		return jobAnswers{}, mapAbort(err)
	}

	value, err := strconv.Atoi(pages)
	if err != nil {
		return jobAnswers{}, fmt.Errorf("parse pages: %w", err)
	}
	answer.Pages = value
	if !cfg.Provided("seed") {
		parsed, err := strconv.ParseUint(seed, 10, 64)
		if err != nil {
			return jobAnswers{}, fmt.Errorf("parse seed: %w", err)
		}
		answer.Seed = parsed
	}
	answer.Out = out
	return answer, nil
}

// buildJobs resolves defaults and validates the whole batch before generating.
func buildJobs(base config.Config, answers []jobAnswers) ([]config.Job, error) {
	jobs := make([]config.Job, 0, len(answers))
	for i, answer := range answers {
		job := config.Job{
			Type:  answer.Type,
			Pages: answer.Pages,
			Seed:  answer.Seed,
			Out:   answer.Out,
			Title: base.Title,
		}
		if job.Type == "" {
			job.Type = base.Type
		}
		if !base.Provided("title") {
			job.Title = typeLabels[job.Type]
		}
		if job.Out == "" {
			job.Out = fmt.Sprintf("%s-%d-%d.pdf", job.Type, i+1, job.Seed)
		}
		jobs = append(jobs, job)
	}

	if err := config.ValidateJobs(jobs); err != nil {
		return nil, err
	}
	return jobs, nil
}

func formatJobs(jobs []config.Job) string {
	result := ""
	for i, job := range jobs {
		if i > 0 {
			result += "\n"
		}
		result += fmt.Sprintf("%d. %s  %d halaman  %s", i+1, job.Type, job.Pages, job.Out)
	}
	return result
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

func validateSeed(value string) error {
	if _, err := strconv.ParseUint(value, 10, 64); err != nil {
		return errors.New("masukkan angka bulat 0 atau lebih")
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
