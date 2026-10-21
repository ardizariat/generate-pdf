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

func TestBuildJobsDefaults(t *testing.T) {
	base := loadConfig(t)
	answers := []jobAnswers{
		{Type: config.TypeReport, Pages: 2, Seed: 7},
		{Pages: 1, Seed: 9, Out: "custom.pdf"},
	}

	jobs, err := buildJobs(base, answers)
	require.NoError(t, err)
	require.Len(t, jobs, 2)

	assert.Equal(t, config.Job{Type: config.TypeReport, Pages: 2, Seed: 7, Out: "laporan-1-7.pdf", Title: "Laporan"}, jobs[0])
	assert.Equal(t, config.Job{Type: config.TypeInvoice, Pages: 1, Seed: 9, Out: "custom.pdf", Title: "Invoice"}, jobs[1])
}

func TestBuildJobsUsesProvidedTitle(t *testing.T) {
	base := loadConfig(t, "--title", "Judul Khusus")
	jobs, err := buildJobs(base, []jobAnswers{{Type: config.TypeReceipt, Pages: 1, Seed: 42}})
	require.NoError(t, err)
	assert.Equal(t, "Judul Khusus", jobs[0].Title)
}

func TestBuildJobsRejectsInvalidBatch(t *testing.T) {
	base := loadConfig(t)

	_, err := buildJobs(base, []jobAnswers{{Type: config.TypeLetter, Pages: 2, Seed: 42}})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "surat-resmi")

	_, err = buildJobs(base, []jobAnswers{{Type: config.TypeInvoice, Pages: 0, Seed: 42}})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Pages=0")
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

	assert.NoError(t, validateSeed("0"))
	assert.NoError(t, validateSeed("18446744073709551615"))
	for _, value := range []string{"", "-1", "abc"} {
		assert.Error(t, validateSeed(value), value)
	}
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
