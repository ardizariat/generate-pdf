package prompt

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"charm.land/huh/v2"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ardizariat/generate-pdf/internal/config"
)

// loadConfig builds a Config from the same flags package cli registers.
func loadConfig(t *testing.T, args ...string) config.Config {
	t.Helper()
	t.Chdir(t.TempDir())

	cmd := &cobra.Command{Use: "pdfgen"}
	flags := cmd.Flags()
	flags.StringP("type", "t", config.TypeInvoice, "")
	flags.StringP("out", "o", "", "")
	flags.String("out-dir", config.DefaultOutDir, "")
	flags.IntP("rows", "r", 0, "")
	flags.Uint64P("seed", "s", 42, "")
	flags.String("title", "", "")
	flags.String("batch", "", "")
	flags.Bool("no-input", false, "")
	flags.Bool("force", false, "")
	flags.Bool("verify", false, "")
	flags.BoolP("verbose", "v", false, "")
	flags.String("config", "", "")

	require.NoError(t, flags.Parse(args))
	cfg, err := config.Load(cmd)
	require.NoError(t, err)
	return cfg
}

func TestBuildJobsCyclesTypes(t *testing.T) {
	base := loadConfig(t)
	jobs, err := buildJobs(base, []string{config.TypeInvoice, config.TypeReport}, 3, 2)
	require.NoError(t, err)
	require.Len(t, jobs, 3)

	assert.Equal(t, config.Job{Type: config.TypeInvoice, Pages: 2, Seed: 42, Out: "output/invoice-1-42.pdf", Title: "Invoice"}, jobs[0])
	assert.Equal(t, config.Job{Type: config.TypeReport, Pages: 2, Seed: 42, Out: "output/laporan-2-42.pdf", Title: "Laporan"}, jobs[1])
	assert.Equal(t, config.Job{Type: config.TypeInvoice, Pages: 2, Seed: 42, Out: "output/invoice-3-42.pdf", Title: "Invoice"}, jobs[2])
}

func TestBuildJobsForcesSinglePageWithoutTable(t *testing.T) {
	base := loadConfig(t)
	jobs, err := buildJobs(base, []string{config.TypeLetter, config.TypeCertificate, config.TypeReceipt}, 3, 4)
	require.NoError(t, err)
	require.Len(t, jobs, 3)
	for _, job := range jobs {
		assert.Equal(t, 1, job.Pages, job.Type)
	}
}

func TestBuildJobsUsesProvidedTitleAndSeed(t *testing.T) {
	base := loadConfig(t, "--title", "Judul Khusus", "--seed", "9")
	jobs, err := buildJobs(base, []string{config.TypeReceipt}, 1, 1)
	require.NoError(t, err)
	assert.Equal(t, "Judul Khusus", jobs[0].Title)
	assert.Equal(t, uint64(9), jobs[0].Seed)
	assert.Equal(t, "output/kuitansi-1-9.pdf", jobs[0].Out)
}

func TestBuildJobsRejectsEmptySelection(t *testing.T) {
	base := loadConfig(t)

	_, err := buildJobs(base, nil, 1, 1)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "minimal satu jenis")

	_, err = buildJobs(base, []string{config.TypeInvoice}, 0, 1)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "minimal satu file")
}

func TestFillBatchRequiresTerminal(t *testing.T) {
	cfg := loadConfig(t, "--no-input")
	_, err := FillBatch(cfg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "terminal interaktif")
}

func TestMapAbort(t *testing.T) {
	aborted := mapAbort(huh.ErrUserAborted)
	require.Error(t, aborted)
	assert.Contains(t, aborted.Error(), "dibatalkan")
	assert.ErrorIs(t, aborted, huh.ErrUserAborted)

	wrapped := mapAbort(fmt.Errorf("run form: %w", huh.ErrUserAborted))
	require.Error(t, wrapped)
	assert.ErrorIs(t, wrapped, huh.ErrUserAborted)

	other := errors.New("boom")
	assert.Equal(t, other, mapAbort(other))
	assert.NoError(t, mapAbort(nil))
}

func TestValidateInput(t *testing.T) {
	for _, value := range []string{"1", "2", "10"} {
		assert.NoError(t, validateCount(value), value)
	}
	for _, value := range []string{"", "0", "-1", "abc"} {
		assert.Error(t, validateCount(value), value)
	}

	for _, value := range []string{"1", "3"} {
		assert.NoError(t, validatePages(value), value)
	}
	for _, value := range []string{"", "0", "-1", "abc"} {
		assert.Error(t, validatePages(value), value)
	}

	assert.NoError(t, validateTypes([]string{config.TypeInvoice}))
	assert.Error(t, validateTypes(nil))
}

func TestPipeStdinIsNotInteractive(t *testing.T) {
	read, write, err := os.Pipe()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = read.Close()
		_ = write.Close()
	})

	previous := os.Stdin
	os.Stdin = read
	t.Cleanup(func() { os.Stdin = previous })

	assert.False(t, IsInteractive())
}
