package cli

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/ardizariat/generate-pdf/internal/config"
	"github.com/ardizariat/generate-pdf/internal/data"
	"github.com/ardizariat/generate-pdf/internal/logging"
	"github.com/ardizariat/generate-pdf/internal/prompt"
	"github.com/ardizariat/generate-pdf/internal/report"
)

// newLogger is logging.New; tests replace it to observe the verbose logger.
var newLogger = logging.New

// runGenerate executes the explicit generate sequence; rendering lives in report.
func runGenerate(cmd *cobra.Command, logger *zap.Logger) error {
	cfg, err := config.Load(cmd)
	if err != nil {
		return fmt.Errorf("muat config: %w", err)
	}
	if err := config.Validate(cfg); err != nil {
		return fmt.Errorf("validasi config: %w", err)
	}

	// main builds a warn-level logger before flags are parsed, so verbose
	// needs its own logger once the final config is known.
	if cfg.Verbose {
		logger = newLogger(true)
	}

	cfg, err = prompt.Fill(cmd, cfg)
	if err != nil {
		return fmt.Errorf("isi prompt: %w", err)
	}
	if err := config.Validate(cfg); err != nil {
		return fmt.Errorf("validasi config akhir: %w", err)
	}

	if _, err := os.Stat(cfg.Out); err == nil {
		if !cfg.Force {
			return fmt.Errorf("file %s sudah ada; gunakan --force untuk menimpa", cfg.Out)
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("cek %s: %w", cfg.Out, err)
	}

	start := time.Now()

	payload, err := data.Generate(cfg.Type, cfg.Rows, cfg.Seed)
	if err != nil {
		return fmt.Errorf("buat data %s: %w", cfg.Type, err)
	}

	pages, size, err := report.Generate(payload, cfg.Out, cfg.Verify)
	if err != nil {
		return fmt.Errorf("buat pdf %s: %w", cfg.Out, err)
	}

	logger.Info("generate selesai",
		zap.String("type", cfg.Type),
		zap.String("out", cfg.Out),
		zap.Int("rows", cfg.Rows),
		zap.Int("pages", pages),
		zap.Int64("bytes", size),
		zap.Duration("duration", time.Since(start)),
	)
	return nil
}
