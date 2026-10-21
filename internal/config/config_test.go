package config

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newCommand mimics the flags registered by package cli without running the command.
func newCommand(t *testing.T, args ...string) *cobra.Command {
	t.Helper()

	cmd := &cobra.Command{
		Use: "generate-pdf",
		RunE: func(*cobra.Command, []string) error {
			return nil
		},
	}
	flags := cmd.Flags()
	flags.StringP("type", "t", TypeInvoice, "")
	flags.String("lang", "en", "")
	flags.StringP("out", "o", "", "")
	flags.String("out-dir", DefaultOutDir, "")
	flags.IntP("rows", "r", 0, "")
	flags.Uint64P("seed", "s", 42, "")
	flags.String("title", "", "")
	flags.Bool("no-input", false, "")
	flags.Bool("force", false, "")
	flags.Bool("verify", false, "")
	flags.BoolP("verbose", "v", false, "")
	flags.String("config", "", "")

	require.NoError(t, flags.Parse(args))
	return cmd
}

func writeConfigFile(t *testing.T, dir, content string) {
	t.Helper()
	require.NoError(t, os.WriteFile(filepath.Join(dir, configName+".yaml"), []byte(content), 0o644))
}

func TestLoadDefaults(t *testing.T) {
	t.Chdir(t.TempDir())

	cfg, err := Load(newCommand(t))
	require.NoError(t, err)

	assert.Equal(t, TypeInvoice, cfg.Type)
	assert.Equal(t, 8, cfg.Rows)
	assert.Equal(t, "Invoice", cfg.Title)
	requireDefaultName(t, cfg.Out, "output", TypeInvoice)
	assert.Equal(t, DefaultOutDir, cfg.OutDir)
	assert.Equal(t, uint64(42), cfg.Seed)
	assert.Empty(t, cfg.ConfigPath)
}

func TestLoadOutDir(t *testing.T) {
	t.Chdir(t.TempDir())

	cfg, err := Load(newCommand(t, "--out-dir", "hasil"))
	require.NoError(t, err)
	requireDefaultName(t, cfg.Out, "hasil", TypeInvoice)

	cfg, err = Load(newCommand(t, "--out", "custom.pdf", "--out-dir", "hasil"))
	require.NoError(t, err)
	assert.Equal(t, "custom.pdf", cfg.Out)

	cfg, err = Load(newCommand(t, "--out", "hasil/custom.pdf"))
	require.NoError(t, err)
	assert.Equal(t, "hasil/custom.pdf", cfg.Out)
}

func TestLoadPrecedence(t *testing.T) {
	tests := []struct {
		name     string
		file     string
		env      map[string]string
		args     []string
		wantType string
		wantRows int
		wantOut  string
	}{
		{
			name:     "default",
			wantType: TypeInvoice,
			wantRows: 8,
		},
		{
			name:     "file over default",
			file:     "type: report\nrows: 5\nout: from-file.pdf\n",
			wantType: TypeReport,
			wantRows: 5,
			wantOut:  "from-file.pdf",
		},
		{
			name:     "env over file",
			file:     "type: report\nrows: 5\nout: from-file.pdf\n",
			env:      map[string]string{"GENPDF_TYPE": TypeDelivery, "GENPDF_ROWS": "9", "GENPDF_OUT": "from-env.pdf"},
			wantType: TypeDelivery,
			wantRows: 9,
			wantOut:  "from-env.pdf",
		},
		{
			name:     "flag over env",
			file:     "type: report\nrows: 5\nout: from-file.pdf\n",
			env:      map[string]string{"GENPDF_TYPE": TypeDelivery, "GENPDF_ROWS": "9", "GENPDF_OUT": "from-env.pdf"},
			args:     []string{"--type", TypeCertificate, "--rows", "3", "--out", "from-flag.pdf"},
			wantType: TypeCertificate,
			wantRows: 3,
			wantOut:  "from-flag.pdf",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Chdir(t.TempDir())
			for key, value := range tc.env {
				t.Setenv(key, value)
			}
			if tc.file != "" {
				writeConfigFile(t, ".", tc.file)
			}

			cfg, err := Load(newCommand(t, tc.args...))
			require.NoError(t, err)

			assert.Equal(t, tc.wantType, cfg.Type)
			assert.Equal(t, tc.wantRows, cfg.Rows)
			if tc.wantOut == "" {
				requireDefaultName(t, cfg.Out, "output", cfg.Type)
			} else {
				assert.Equal(t, tc.wantOut, cfg.Out)
			}
		})
	}
}

func TestLoadDefaultsPerType(t *testing.T) {
	tests := []struct {
		typ       string
		wantRows  int
		wantTitle string
	}{
		{typ: TypeInvoice, wantRows: 8, wantTitle: "Invoice"},
		{typ: TypeLetter, wantRows: 0, wantTitle: "Letter"},
		{typ: TypeReport, wantRows: 70, wantTitle: "Report"},
		{typ: TypeReceipt, wantRows: 0, wantTitle: "Receipt"},
		{typ: TypeDelivery, wantRows: 25, wantTitle: "Delivery Note"},
		{typ: TypeCertificate, wantRows: 0, wantTitle: "Certificate"},
	}

	for _, tc := range tests {
		t.Run(tc.typ, func(t *testing.T) {
			t.Chdir(t.TempDir())

			cfg, err := Load(newCommand(t, "--type", tc.typ, "--seed", "7"))
			require.NoError(t, err)

			assert.Equal(t, tc.wantRows, cfg.Rows)
			assert.Equal(t, tc.wantTitle, cfg.Title)
			requireDefaultName(t, cfg.Out, "output", tc.typ)
		})
	}
}

func TestProvided(t *testing.T) {
	t.Chdir(t.TempDir())
	writeConfigFile(t, ".", "rows: 5\n")
	t.Setenv("GENPDF_TITLE", "Env Title")

	cfg, err := Load(newCommand(t, "--force"))
	require.NoError(t, err)

	assert.True(t, cfg.Provided("rows"), "config file")
	assert.True(t, cfg.Provided("title"), "env")
	assert.True(t, cfg.Provided("force"), "flag")
	assert.False(t, cfg.Provided("out"))

	plainDir := t.TempDir()
	t.Chdir(plainDir)
	plain, err := Load(newCommand(t))
	require.NoError(t, err)
	assert.False(t, plain.Provided("rows"))
	assert.False(t, plain.Provided("type"))
}

func TestProvidedConfigFlag(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	path := filepath.Join(dir, "custom.yaml")
	require.NoError(t, os.WriteFile(path, []byte("type: receipt\n"), 0o644))

	cfg, err := Load(newCommand(t, "--config", path))
	require.NoError(t, err)

	assert.True(t, cfg.Provided("config"))
	assert.True(t, cfg.Provided("type"))
	assert.Equal(t, path, cfg.ConfigPath)
	assert.Equal(t, TypeReceipt, cfg.Type)
}

func TestLoadConfigPathFromEnv(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	envPath := filepath.Join(dir, "from-env.yaml")
	require.NoError(t, os.WriteFile(envPath, []byte("type: receipt\n"), 0o644))
	t.Setenv("GENPDF_CONFIG", envPath)

	cfg, err := Load(newCommand(t))
	require.NoError(t, err)
	assert.Equal(t, envPath, cfg.ConfigPath)
	assert.Equal(t, TypeReceipt, cfg.Type)

	flagPath := filepath.Join(dir, "from-flag.yaml")
	require.NoError(t, os.WriteFile(flagPath, []byte("type: report\n"), 0o644))

	cfg, err = Load(newCommand(t, "--config", flagPath))
	require.NoError(t, err)
	assert.Equal(t, flagPath, cfg.ConfigPath, "flag beats GENPDF_CONFIG")
	assert.Equal(t, TypeReport, cfg.Type)
}

func TestLoadMissingConfigFileIsNotError(t *testing.T) {
	t.Chdir(t.TempDir())

	_, err := Load(newCommand(t))
	assert.NoError(t, err)
}

func TestLoadExplicitMissingConfigFileIsError(t *testing.T) {
	t.Chdir(t.TempDir())

	_, err := Load(newCommand(t, "--config", "missing.yaml"))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "read config")
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr []string
	}{
		{
			name: "valid",
			cfg:  Config{Type: TypeInvoice, Out: "invoice-42.pdf"},
		},
		{
			name:    "unknown type",
			cfg:     Config{Type: "faktur", Out: "faktur-42.pdf"},
			wantErr: []string{"Type", "faktur", "oneof"},
		},
		{
			name:    "negative rows",
			cfg:     Config{Type: TypeInvoice, Out: "invoice-42.pdf", Rows: -1},
			wantErr: []string{"Rows", "-1", "gte"},
		},
		{
			name:    "empty out",
			cfg:     Config{Type: TypeInvoice},
			wantErr: []string{"Out", "required"},
		},
		{
			name: "zero seed is allowed",
			cfg:  Config{Type: TypeReceipt, Out: "receipt-0.pdf", Seed: 0},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := Validate(tc.cfg)
			if len(tc.wantErr) == 0 {
				assert.NoError(t, err)
				return
			}
			require.Error(t, err)
			for _, want := range tc.wantErr {
				assert.Contains(t, err.Error(), want)
			}
		})
	}
}

func TestValidateAcceptsEveryType(t *testing.T) {
	for _, typ := range Types() {
		assert.NoError(t, Validate(Config{Type: typ, Out: typ + "-42.pdf"}), typ)
	}
}

func TestLoadLanguage(t *testing.T) {
	t.Chdir(t.TempDir())

	t.Run("flag", func(t *testing.T) {
		cfg, err := Load(newCommand(t, "--type", TypeReceipt, "--lang", "id"))
		require.NoError(t, err)
		assert.Equal(t, "id", cfg.Language)
		assert.Equal(t, "Kuitansi", cfg.Title)
	})

	t.Run("env", func(t *testing.T) {
		t.Setenv("GENPDF_LANG", "id")
		cfg, err := Load(newCommand(t, "--type", TypeReceipt))
		require.NoError(t, err)
		assert.Equal(t, "id", cfg.Language)
		assert.Equal(t, "Kuitansi", cfg.Title)
	})

	t.Run("file", func(t *testing.T) {
		writeConfigFile(t, ".", "lang: id\n")
		cfg, err := Load(newCommand(t, "--type", TypeReceipt))
		require.NoError(t, err)
		assert.Equal(t, "id", cfg.Language)
		assert.Equal(t, "Kuitansi", cfg.Title)
	})
}

func TestValidateRejectsUnknownLanguage(t *testing.T) {
	err := Validate(Config{Type: TypeInvoice, Language: "fr", Out: "invoice.pdf"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Language")
}

// requireDefaultName checks the <type>_<UTC>_<uuid v7>.pdf default name.
func requireDefaultName(t *testing.T, path, dir, docType string) {
	t.Helper()

	pattern := "^" + regexp.QuoteMeta(docType) +
		`_\d{8}T\d{6}Z_[0-9a-f]{8}-[0-9a-f]{4}-7[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}\.pdf$`
	if dir == "" {
		dir = "."
	}
	assert.Equal(t, dir, filepath.Dir(path))
	assert.Regexp(t, pattern, filepath.Base(path))
}
