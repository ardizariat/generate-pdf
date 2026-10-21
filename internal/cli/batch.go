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
	"github.com/ardizariat/generate-pdf/internal/locale"
	"github.com/ardizariat/generate-pdf/internal/report"
)

// runBatch generates every job, keeps going after failures, and prints a
// per-file summary. It returns an error when at least one job failed.
func runBatch(cmd *cobra.Command, logger *zap.Logger, cfg config.Config, jobs []config.Job) error {
	if err := config.ValidateJobs(jobs); err != nil {
		return err
	}

	out := cmd.OutOrStdout()
	labels := locale.Get(locale.OrDefault(cfg.Language))
	failures := 0
	for i, job := range jobs {
		start := time.Now()
		var pages int
		var size int64
		var jobErr error
		progress(cfg, fmt.Sprintf(labels.ProgressPages, i+1, len(jobs), job.Out, job.Pages), func() {
			pages, size, jobErr = generateJob(job, cfg)
		})
		if jobErr != nil {
			failures++
			if _, writeErr := fmt.Fprintf(out, labels.BatchFailed, job.Out, jobErr); writeErr != nil {
				return fmt.Errorf("write summary: %w", writeErr)
			}
			continue
		}

		logger.Info("generate complete",
			zap.String("type", job.Type),
			zap.String("out", job.Out),
			zap.Int("pages", pages),
			zap.Int64("bytes", size),
			zap.Duration("duration", time.Since(start)),
		)
		if _, err := fmt.Fprintf(out, labels.BatchOK, job.Out, pages, size); err != nil {
			return fmt.Errorf("write summary: %w", err)
		}
	}

	if _, err := fmt.Fprintf(out, labels.BatchSummary, len(jobs)-failures, failures); err != nil {
		return fmt.Errorf("write summary: %w", err)
	}
	if failures > 0 {
		return fmt.Errorf("%d of %d files failed", failures, len(jobs))
	}
	return nil
}

// generateJob writes one batch job after the overwrite check.
func generateJob(job config.Job, cfg config.Config) (pages int, size int64, err error) {
	if _, err := os.Stat(job.Out); err == nil {
		if !cfg.Force {
			return 0, 0, fmt.Errorf("file %s already exists; use --force to overwrite", job.Out)
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return 0, 0, fmt.Errorf("check %s: %w", job.Out, err)
	}

	if err := ensureParentDir(job.Out); err != nil {
		return 0, 0, err
	}

	build := func(rows int) (data.Report, error) {
		r, err := data.Generate(job.Type, rows, job.Seed, locale.OrDefault(cfg.Language))
		if err != nil {
			return data.Report{}, err
		}
		r.Title = job.Title
		return r, nil
	}

	_, pages, size, err = report.GeneratePages(build, job.Out, job.Pages, cfg.Verify)
	if err != nil {
		return 0, 0, fmt.Errorf("build pdf %s: %w", job.Out, err)
	}
	return pages, size, nil
}
