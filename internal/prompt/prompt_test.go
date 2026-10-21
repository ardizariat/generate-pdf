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

func TestFieldsToAsk(t *testing.T) {
	tests := []struct {
		name string
		args []string
		env  map[string]string
		want []string
	}{
		{
			name: "nothing provided",
			want: []string{"type", "rows", "seed", "out"},
		},
		{
			name: "everything provided",
			args: []string{"--type", config.TypeInvoice, "--rows", "8", "--seed", "7", "--out", "out.pdf"},
			want: []string{},
		},
		{
			name: "type without table skips rows",
			args: []string{"--type", config.TypeLetter},
			want: []string{"seed", "out"},
		},
		{
			name: "env provides seed",
			env:  map[string]string{"GENPDF_SEED": "7"},
			want: []string{"type", "rows", "out"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			for key, value := range tc.env {
				t.Setenv(key, value)
			}

			cfg := loadConfig(t, tc.args...)
			assert.Equal(t, tc.want, fieldsToAsk(nil, cfg))
		})
	}
}

func TestApplyTypeDefaults(t *testing.T) {
	cfg := loadConfig(t)
	cfg.Type = config.TypeLetter
	applyTypeDefaults(&cfg, config.TypeInvoice)
	assert.Equal(t, 0, cfg.Rows, "type without table")
	assert.Equal(t, "Surat Resmi", cfg.Title)

	provided := loadConfig(t, "--rows", "5")
	provided.Type = config.TypeCertificate
	applyTypeDefaults(&provided, config.TypeInvoice)
	assert.Equal(t, 5, provided.Rows, "explicit rows must not be overwritten")
}

func TestFillSkipsWhenStdinIsNotTerminal(t *testing.T) {
	cfg := loadConfig(t)
	swapStdin(t, pipeStdin(t))

	got, err := Fill(nil, cfg)
	require.NoError(t, err)
	assert.Equal(t, cfg, got)
}

func TestFillSkipsWhenInputDisabled(t *testing.T) {
	cfg := loadConfig(t, "--no-input")

	got, err := Fill(nil, cfg)
	require.NoError(t, err)
	assert.Equal(t, cfg, got)
}

func TestMapAbort(t *testing.T) {
	aborted := mapAbort(huh.ErrUserAborted)
	require.Error(t, aborted)
	assert.Contains(t, aborted.Error(), "dibatalkan")
	assert.True(t, errors.Is(aborted, huh.ErrUserAborted))

	wrapped := mapAbort(fmt.Errorf("run form: %w", huh.ErrUserAborted))
	require.Error(t, wrapped)
	assert.True(t, errors.Is(wrapped, huh.ErrUserAborted))

	other := errors.New("boom")
	assert.Equal(t, other, mapAbort(other))
	assert.NoError(t, mapAbort(nil))
}

func TestValidateInput(t *testing.T) {
	rows := []struct {
		value string
		valid bool
	}{
		{value: "0", valid: true},
		{value: "8", valid: true},
		{value: "-1", valid: false},
		{value: "abc", valid: false},
	}
	for _, tc := range rows {
		err := validateRows(tc.value)
		if tc.valid {
			assert.NoError(t, err, tc.value)
		} else {
			assert.Error(t, err, tc.value)
		}
	}

	assert.NoError(t, validateSeed("0"))
	assert.NoError(t, validateSeed("18446744073709551615"))
	for _, value := range []string{"", "-1", "abc"} {
		assert.Error(t, validateSeed(value), value)
	}

	assert.NoError(t, validateOut("invoice-42.pdf"))
	assert.Error(t, validateOut(""))
}

// pipeStdin returns the read end of a pipe, which is not a character device.
func pipeStdin(t *testing.T) *os.File {
	t.Helper()
	read, write, err := os.Pipe()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = read.Close()
		_ = write.Close()
	})
	return read
}

func swapStdin(t *testing.T, stdin *os.File) {
	t.Helper()
	previous := os.Stdin
	os.Stdin = stdin
	t.Cleanup(func() { os.Stdin = previous })
}
