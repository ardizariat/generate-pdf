package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/ardizariat/generate-pdf/internal/config"
	"github.com/ardizariat/generate-pdf/internal/data"
	"github.com/ardizariat/generate-pdf/internal/locale"
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
		return fmt.Errorf("load config: %w", err)
	}
	if err := config.Validate(cfg); err != nil {
		return fmt.Errorf("validate config: %w", err)
	}

	// main builds a warn-level logger before flags are parsed, so verbose
	// needs its own logger once the final config is known.
	if cfg.Verbose {
		logger = newLogger(true)
	}

	if cfg.Batch != "" {
		jobs, err := config.LoadBatch(cfg.Batch, cfg)
		if err != nil {
			return fmt.Errorf("load batch: %w", err)
		}
		return runBatch(cmd, logger, cfg, jobs)
	}

	if !cfg.NoInput && prompt.IsInteractive() {
		jobs, err := prompt.FillBatch(cfg)
		if err != nil {
			return fmt.Errorf("fill prompt: %w", err)
		}
		return runBatch(cmd, logger, cfg, jobs)
	}

	if _, err := os.Stat(cfg.Out); err == nil {
		if !cfg.Force {
			return fmt.Errorf("file %s already exists; use --force to overwrite", cfg.Out)
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("check %s: %w", cfg.Out, err)
	}

	if err := ensureParentDir(cfg.Out); err != nil {
		return err
	}

	start := time.Now()
	labels := locale.Get(locale.OrDefault(cfg.Language))

	var pages int
	var size int64
	var genErr error
	progress(cfg, fmt.Sprintf(labels.ProgressRows, 1, 1, cfg.Out, cfg.Rows), func() {
		payload, err := data.Generate(cfg.Type, cfg.Rows, cfg.Seed, locale.OrDefault(cfg.Language))
		if err != nil {
			genErr = fmt.Errorf("build data %s: %w", cfg.Type, err)
			return
		}
		payload.Title = cfg.Title

		pages, size, genErr = report.Generate(payload, cfg.Out, cfg.Verify)
		if genErr != nil {
			genErr = fmt.Errorf("build pdf %s: %w", cfg.Out, genErr)
		}
	})
	if genErr != nil {
		return genErr
	}

	logger.Info("generate complete",
		zap.String("type", cfg.Type),
		zap.String("out", cfg.Out),
		zap.Int("rows", cfg.Rows),
		zap.Int("pages", pages),
		zap.Int64("bytes", size),
		zap.Duration("duration", time.Since(start)),
	)
	return nil
}

// ensureParentDir creates the output directory tree so a fresh checkout can
// generate into output/ without a manual mkdir.
func ensureParentDir(path string) error {
	dir := filepath.Dir(path)
	if dir == "." || dir == "" {
		return nil
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create directory %s: %w", dir, err)
	}
	return nil
}
