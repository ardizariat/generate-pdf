package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newCommand meniru flag yang didaftarkan package cli, tanpa menjalankan perintah.
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
	assert.Equal(t, "invoice-42.pdf", cfg.Out)
	assert.Equal(t, uint64(42), cfg.Seed)
	assert.Empty(t, cfg.ConfigPath)
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
			wantOut:  "invoice-42.pdf",
		},
		{
			name:     "file over default",
			file:     "type: laporan\nrows: 5\nout: from-file.pdf\n",
			wantType: TypeReport,
			wantRows: 5,
			wantOut:  "from-file.pdf",
		},
		{
			name:     "env over file",
			file:     "type: laporan\nrows: 5\nout: from-file.pdf\n",
			env:      map[string]string{"GENPDF_TYPE": TypeDelivery, "GENPDF_ROWS": "9", "GENPDF_OUT": "from-env.pdf"},
			wantType: TypeDelivery,
			wantRows: 9,
			wantOut:  "from-env.pdf",
		},
		{
			name:     "flag over env",
			file:     "type: laporan\nrows: 5\nout: from-file.pdf\n",
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
			assert.Equal(t, tc.wantOut, cfg.Out)
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
		{typ: TypeLetter, wantRows: 0, wantTitle: "Surat Resmi"},
		{typ: TypeReport, wantRows: 70, wantTitle: "Laporan"},
		{typ: TypeReceipt, wantRows: 0, wantTitle: "Kuitansi"},
		{typ: TypeDelivery, wantRows: 25, wantTitle: "Surat Jalan"},
		{typ: TypeCertificate, wantRows: 0, wantTitle: "Sertifikat"},
	}

	for _, tc := range tests {
		t.Run(tc.typ, func(t *testing.T) {
			t.Chdir(t.TempDir())

			cfg, err := Load(newCommand(t, "--type", tc.typ, "--seed", "7"))
			require.NoError(t, err)

			assert.Equal(t, tc.wantRows, cfg.Rows)
			assert.Equal(t, tc.wantTitle, cfg.Title)
			assert.Equal(t, tc.typ+"-7.pdf", cfg.Out)
		})
	}
}

func TestProvided(t *testing.T) {
	t.Chdir(t.TempDir())
	writeConfigFile(t, ".", "rows: 5\n")
	t.Setenv("GENPDF_TITLE", "Judul Env")

	cfg, err := Load(newCommand(t, "--force"))
	require.NoError(t, err)

	assert.True(t, cfg.Provided("rows"), "file config")
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
	require.NoError(t, os.WriteFile(path, []byte("type: kuitansi\n"), 0o644))

	cfg, err := Load(newCommand(t, "--config", path))
	require.NoError(t, err)

	assert.True(t, cfg.Provided("config"))
	assert.True(t, cfg.Provided("type"))
	assert.Equal(t, path, cfg.ConfigPath)
	assert.Equal(t, TypeReceipt, cfg.Type)
}

func TestLoadMissingConfigFileIsNotError(t *testing.T) {
	t.Chdir(t.TempDir())

	_, err := Load(newCommand(t))
	assert.NoError(t, err)
}

func TestLoadExplicitMissingConfigFileIsError(t *testing.T) {
	t.Chdir(t.TempDir())

	_, err := Load(newCommand(t, "--config", "tidak-ada.yaml"))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "baca config")
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
			name:    "jenis tak dikenal",
			cfg:     Config{Type: "faktur", Out: "faktur-42.pdf"},
			wantErr: []string{"Type", "faktur", "oneof"},
		},
		{
			name:    "rows negatif",
			cfg:     Config{Type: TypeInvoice, Out: "invoice-42.pdf", Rows: -1},
			wantErr: []string{"Rows", "-1", "gte"},
		},
		{
			name:    "out kosong",
			cfg:     Config{Type: TypeInvoice},
			wantErr: []string{"Out", "required"},
		},
		{
			name: "seed nol diizinkan",
			cfg:  Config{Type: TypeReceipt, Out: "kuitansi-0.pdf", Seed: 0},
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
