// Package prompt renders the interactive forms that fill missing configuration.
package prompt

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"charm.land/huh/v2"
	"github.com/spf13/cobra"

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

// defaultRows mirrors the per-type row defaults in package config. Types
// missing from this map have no table, so rows is never asked for them.
// shortcut: duplicated because config does not export per-type defaults;
// export them from config if a second consumer appears.
var defaultRows = map[string]int{
	config.TypeInvoice:  8,
	config.TypeReport:   70,
	config.TypeDelivery: 25,
}

// Fill asks for the configuration fields that were not supplied by flag, env,
// or config file. It is a no-op when input is disabled or stdin is not a
// terminal. The returned Config still needs final validation by the caller.
func Fill(cmd *cobra.Command, cfg config.Config) (config.Config, error) {
	if cfg.NoInput || !isTerminal() {
		return cfg, nil
	}

	fields := fieldsToAsk(cmd, cfg)
	previousType := cfg.Type
	if contains(fields, "type") {
		if err := askType(&cfg); err != nil {
			return cfg, err
		}
	}
	applyTypeDefaults(&cfg, previousType)
	if !cfg.Provided("out") {
		cfg.Out = fmt.Sprintf("%s-%d.pdf", cfg.Type, cfg.Seed)
	}

	fields = fieldsToAsk(cmd, cfg)
	if err := askDetails(&cfg, fields); err != nil {
		return cfg, err
	}
	return cfg, nil
}

// fieldsToAsk lists the config keys the forms must ask for. cmd is accepted
// for symmetry with config.Load; cfg.Provided already covers flags, env, and
// the config file.
func fieldsToAsk(_ *cobra.Command, cfg config.Config) []string {
	fields := make([]string, 0, 4)
	if !cfg.Provided("type") {
		fields = append(fields, "type")
	}
	if !cfg.Provided("rows") && hasTable(cfg.Type) {
		fields = append(fields, "rows")
	}
	if !cfg.Provided("seed") {
		fields = append(fields, "seed")
	}
	if !cfg.Provided("out") {
		fields = append(fields, "out")
	}
	return fields
}

// applyTypeDefaults refreshes defaults that depend on the selected type when
// the user changes it in the form.
func applyTypeDefaults(cfg *config.Config, previousType string) {
	if cfg.Type == previousType {
		return
	}
	if !cfg.Provided("rows") {
		cfg.Rows = defaultRows[cfg.Type]
	}
	if !cfg.Provided("title") {
		cfg.Title = typeLabels[cfg.Type]
	}
}

func askType(cfg *config.Config) error {
	types := config.Types()
	options := make([]huh.Option[string], 0, len(types))
	for _, typ := range types {
		options = append(options, huh.NewOption(typeLabels[typ], typ))
	}

	form := huh.NewForm(huh.NewGroup(
		huh.NewSelect[string]().
			Title("Jenis dokumen").
			Options(options...).
			Value(&cfg.Type),
	))
	if err := form.Run(); err != nil {
		return mapAbort(err)
	}
	return nil
}

func askDetails(cfg *config.Config, fields []string) error {
	rows := strconv.Itoa(cfg.Rows)
	seed := strconv.FormatUint(cfg.Seed, 10)
	out := cfg.Out

	inputs := make([]huh.Field, 0, 3)
	if contains(fields, "rows") {
		inputs = append(inputs, huh.NewInput().
			Title("Jumlah baris").
			Value(&rows).
			Validate(validateRows))
	}
	if contains(fields, "seed") {
		inputs = append(inputs, huh.NewInput().
			Title("Seed").
			Value(&seed).
			Validate(validateSeed))
	}
	if contains(fields, "out") {
		inputs = append(inputs, huh.NewInput().
			Title("Path output").
			Value(&out).
			Validate(validateOut))
	}
	if len(inputs) == 0 {
		return nil
	}

	if err := huh.NewForm(huh.NewGroup(inputs...)).Run(); err != nil {
		return mapAbort(err)
	}

	if contains(fields, "rows") {
		value, err := strconv.Atoi(rows)
		if err != nil {
			return fmt.Errorf("parse rows: %w", err)
		}
		cfg.Rows = value
	}
	if contains(fields, "seed") {
		value, err := strconv.ParseUint(seed, 10, 64)
		if err != nil {
			return fmt.Errorf("parse seed: %w", err)
		}
		cfg.Seed = value
	}
	if contains(fields, "out") {
		cfg.Out = out
	}
	return nil
}

func validateRows(value string) error {
	rows, err := strconv.Atoi(value)
	if err != nil {
		return errors.New("masukkan angka bulat")
	}
	if rows < 0 {
		return errors.New("minimal 0")
	}
	return nil
}

func validateSeed(value string) error {
	if _, err := strconv.ParseUint(value, 10, 64); err != nil {
		return errors.New("masukkan angka bulat 0 atau lebih")
	}
	return nil
}

func validateOut(value string) error {
	if value == "" {
		return errors.New("path output wajib diisi")
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

func hasTable(docType string) bool {
	_, ok := defaultRows[docType]
	return ok
}

func contains(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
