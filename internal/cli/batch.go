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
	"github.com/ardizariat/generate-pdf/internal/report"
)

// runBatch generates every job, keeps going after failures, and prints a
// per-file summary. It returns an error when at least one job failed.
func runBatch(cmd *cobra.Command, logger *zap.Logger, cfg config.Config, jobs []config.Job) error {
	if err := config.ValidateJobs(jobs); err != nil {
		return err
	}

	out := cmd.OutOrStdout()
	failures := 0
	for _, job := range jobs {
		start := time.Now()
		pages, size, err := generateJob(job, cfg)
		if err != nil {
			failures++
			if _, writeErr := fmt.Fprintf(out, "gagal  %s: %v\n", job.Out, err); writeErr != nil {
				return fmt.Errorf("tulis ringkasan: %w", writeErr)
			}
			continue
		}

		logger.Info("generate selesai",
			zap.String("type", job.Type),
			zap.String("out", job.Out),
			zap.Int("pages", pages),
			zap.Int64("bytes", size),
			zap.Duration("duration", time.Since(start)),
		)
		if _, err := fmt.Fprintf(out, "ok     %s  %d halaman  %d bytes\n", job.Out, pages, size); err != nil {
			return fmt.Errorf("tulis ringkasan: %w", err)
		}
	}

	if _, err := fmt.Fprintf(out, "ringkasan: %d sukses, %d gagal\n", len(jobs)-failures, failures); err != nil {
		return fmt.Errorf("tulis ringkasan: %w", err)
	}
	if failures > 0 {
		return fmt.Errorf("%d dari %d file gagal", failures, len(jobs))
	}
	return nil
}

// generateJob writes one batch job after the overwrite check.
func generateJob(job config.Job, cfg config.Config) (pages int, size int64, err error) {
	if _, err := os.Stat(job.Out); err == nil {
		if !cfg.Force {
			return 0, 0, fmt.Errorf("file %s sudah ada; gunakan --force untuk menimpa", job.Out)
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return 0, 0, fmt.Errorf("cek %s: %w", job.Out, err)
	}

	build := func(rows int) (data.Report, error) {
		r, err := data.Generate(job.Type, rows, job.Seed)
		if err != nil {
			return data.Report{}, err
		}
		r.Title = job.Title
		return r, nil
	}

	_, pages, size, err = report.GeneratePages(build, job.Out, job.Pages, cfg.Verify)
	if err != nil {
		return 0, 0, fmt.Errorf("buat pdf %s: %w", job.Out, err)
	}
	return pages, size, nil
}
