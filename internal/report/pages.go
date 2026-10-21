package report

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ardizariat/generate-pdf/internal/data"
)

// maxFitRows bounds the page search. Targets needing more rows than this fail
// instead of rendering forever; raise it if a layout ever needs denser fills.
const maxFitRows = 2000

// GeneratePages renders the document built for a row count so the saved file
// has exactly targetPages pages, using the largest row count that still fits
// the target. The page count grows monotonically with rows, so a binary search
// finds that maximum without guessing row density.
func GeneratePages(build func(rows int) (data.Report, error), out string, targetPages int, verify bool) (rows int, pages int, size int64, err error) {
	if targetPages < 1 {
		return 0, 0, 0, fmt.Errorf("target pages must be >= 1, got %d", targetPages)
	}

	tmpDir, err := os.MkdirTemp("", "pdfgen-pages-*")
	if err != nil {
		return 0, 0, 0, fmt.Errorf("create temporary directory: %w", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	counts := make(map[int]int)
	pagesAt := func(rowCount int) (int, error) {
		if count, ok := counts[rowCount]; ok {
			return count, nil
		}
		r, err := build(rowCount)
		if err != nil {
			return 0, fmt.Errorf("build data for %d rows: %w", rowCount, err)
		}
		path := filepath.Join(tmpDir, fmt.Sprintf("rows-%d.pdf", rowCount))
		count, _, err := Generate(r, path, false)
		if err != nil {
			return 0, fmt.Errorf("render %d rows: %w", rowCount, err)
		}
		counts[rowCount] = count
		return count, nil
	}

	// Grow the upper bound until it crosses the target, so small targets stay cheap.
	lo, hi := 0, 1
	for {
		count, err := pagesAt(hi)
		if err != nil {
			return 0, 0, 0, err
		}
		if count > targetPages {
			break
		}
		lo = hi
		if hi == maxFitRows {
			if count == targetPages {
				break
			}
			return 0, 0, 0, fmt.Errorf("target %d pages not reached (maximum %d pages at %d rows)", targetPages, count, hi)
		}
		if hi > maxFitRows/2 {
			hi = maxFitRows
		} else {
			hi *= 2
		}
	}

	for lo+1 < hi {
		mid := (lo + hi) / 2
		count, err := pagesAt(mid)
		if err != nil {
			return 0, 0, 0, err
		}
		if count <= targetPages {
			lo = mid
		} else {
			hi = mid
		}
	}

	count, err := pagesAt(lo)
	if err != nil {
		return 0, 0, 0, err
	}
	if count != targetPages {
		return 0, 0, 0, fmt.Errorf("target %d pages not reached (best %d pages at %d rows)", targetPages, count, lo)
	}

	r, err := build(lo)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("build data for %d rows: %w", lo, err)
	}
	staged, err := stagePath(out)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("create temporary file for %s: %w", out, err)
	}
	defer func() { _ = os.Remove(staged) }()

	pages, size, err = generate(r, staged, out, verify)
	if err != nil {
		return 0, 0, 0, err
	}
	if pages != targetPages {
		return 0, 0, 0, fmt.Errorf("target %d pages not reached (%s contains %d pages)", targetPages, out, pages)
	}
	if err := commit(staged, out); err != nil {
		return 0, 0, 0, fmt.Errorf("save %s: %w", out, err)
	}
	return lo, pages, size, nil
}
