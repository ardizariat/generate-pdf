package config

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeManifest(t *testing.T, content string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "jobs.yaml")
	require.NoError(t, os.WriteFile(path, []byte(content), 0o644))
	return path
}

func TestLoadBatchResolvesDefaults(t *testing.T) {
	path := writeManifest(t, `jobs:
  - type: laporan
  - out: custom.pdf
`)

	base := Config{Type: TypeInvoice, Seed: 42}
	jobs, err := LoadBatch(path, base)
	require.NoError(t, err)

	require.Len(t, jobs, 2)
	assert.Equal(t, TypeReport, jobs[0].Type)
	assert.Equal(t, 1, jobs[0].Pages)
	assert.Equal(t, uint64(42), jobs[0].Seed)
	assert.Equal(t, "Laporan", jobs[0].Title)
	requireDefaultName(t, jobs[0].Out, "", TypeReport)
	assert.Equal(t, Job{Type: TypeInvoice, Pages: 1, Seed: 42, Out: "custom.pdf", Title: "Invoice"}, jobs[1])
}

func TestLoadBatchJobOverridesGlobal(t *testing.T) {
	path := writeManifest(t, `jobs:
  - type: surat-jalan
    pages: 3
    seed: 7
    out: a.pdf
    title: Pengiriman
`)

	base := Config{Type: TypeInvoice, Seed: 42}
	jobs, err := LoadBatch(path, base)
	require.NoError(t, err)

	require.Len(t, jobs, 1)
	assert.Equal(t, Job{Type: TypeDelivery, Pages: 3, Seed: 7, Out: "a.pdf", Title: "Pengiriman"}, jobs[0])
}

func TestLoadBatchInheritsFlags(t *testing.T) {
	t.Chdir(t.TempDir())
	cmd := newCommand(t, "-t", TypeReceipt, "-s", "9")
	base, err := Load(cmd)
	require.NoError(t, err)

	path := writeManifest(t, "jobs:\n  - pages: 1\n")
	jobs, err := LoadBatch(path, base)
	require.NoError(t, err)

	require.Len(t, jobs, 1)
	assert.Equal(t, TypeReceipt, jobs[0].Type)
	assert.Equal(t, uint64(9), jobs[0].Seed)
	requireDefaultName(t, jobs[0].Out, "output", TypeReceipt)
	assert.Equal(t, "Kuitansi", jobs[0].Title)
}

func TestLoadBatchRejectsUnknownFields(t *testing.T) {
	path := writeManifest(t, "jobs:\n  - type: invoice\n    page: 2\n")
	_, err := LoadBatch(path, Config{Type: TypeInvoice, Seed: 42})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "field page not found")
}

func TestLoadBatchErrors(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    []string
	}{
		{name: "empty jobs", content: "jobs: []\n", want: []string{"jobs kosong"}},
		{name: "invalid yaml", content: "jobs: [\n", want: []string{"parse manifest"}},
		{name: "zero pages", content: "jobs:\n  - type: invoice\n    pages: 0\n", want: []string{"job 1", "Pages=0"}},
		{name: "unknown type", content: "jobs:\n  - type: faktur\n", want: []string{"job 1", "faktur"}},
		{name: "letter multi page", content: "jobs:\n  - type: surat-resmi\n    pages: 2\n", want: []string{"surat-resmi", "maksimal 1"}},
		{name: "receipt multi page", content: "jobs:\n  - type: kuitansi\n    pages: 2\n", want: []string{"kuitansi", "maksimal 1"}},
		{name: "certificate multi page", content: "jobs:\n  - type: sertifikat\n    pages: 2\n", want: []string{"sertifikat", "maksimal 1"}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			path := writeManifest(t, tc.content)
			_, err := LoadBatch(path, Config{Type: TypeInvoice, Seed: 42})
			require.Error(t, err)
			for _, want := range tc.want {
				assert.Contains(t, err.Error(), want)
			}
		})
	}
}

func TestLoadBatchMissingFile(t *testing.T) {
	_, err := LoadBatch(filepath.Join(t.TempDir(), "missing.yaml"), Config{Type: TypeInvoice})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "baca manifest")
}

func TestLoadBatchAcceptsMultiPageTableTypes(t *testing.T) {
	path := writeManifest(t, `jobs:
  - type: invoice
    pages: 2
  - type: laporan
    pages: 3
  - type: surat-jalan
    pages: 2
`)
	jobs, err := LoadBatch(path, Config{Type: TypeInvoice, Seed: 42})
	require.NoError(t, err)
	require.Len(t, jobs, 3)
}

func TestValidateJobsEmpty(t *testing.T) {
	assert.Error(t, ValidateJobs(nil))
}

func TestSupportsMultiPage(t *testing.T) {
	for _, docType := range []string{TypeInvoice, TypeReport, TypeDelivery} {
		assert.True(t, SupportsMultiPage(docType), docType)
	}
	for _, docType := range []string{TypeLetter, TypeReceipt, TypeCertificate} {
		assert.False(t, SupportsMultiPage(docType), docType)
	}
}

func TestLoadBatchResolvesOutDir(t *testing.T) {
	abs := filepath.Join(t.TempDir(), "custom.pdf")
	path := writeManifest(t, fmt.Sprintf(`jobs:
  - type: invoice
  - type: invoice
    out: sub/custom.pdf
  - type: invoice
    out: %q
`, abs))

	base := Config{Type: TypeInvoice, Seed: 42, OutDir: "output"}
	jobs, err := LoadBatch(path, base)
	require.NoError(t, err)
	require.Len(t, jobs, 3)

	requireDefaultName(t, jobs[0].Out, "output", TypeInvoice)
	assert.Equal(t, filepath.Join("output", "sub", "custom.pdf"), jobs[1].Out)
	assert.Equal(t, abs, jobs[2].Out)
}

func TestOutputPath(t *testing.T) {
	abs := filepath.Join(t.TempDir(), "a.pdf")
	assert.Equal(t, filepath.Join("output", "a.pdf"), OutputPath("output", "a.pdf"))
	assert.Equal(t, "", OutputPath("output", ""))
	assert.Equal(t, "a.pdf", OutputPath("", "a.pdf"))
	assert.Equal(t, abs, OutputPath("output", abs))
}

func TestDefaultOutputName(t *testing.T) {
	at := time.Date(2026, 10, 21, 14, 15, 30, 0, time.UTC)
	name, err := DefaultOutputName(TypeInvoice, at)
	require.NoError(t, err)
	assert.Regexp(t, `^invoice_20261021T141530Z_[0-9a-f]{8}-[0-9a-f]{4}-7[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}\.pdf$`, name)

	other, err := DefaultOutputName(TypeInvoice, at)
	require.NoError(t, err)
	assert.NotEqual(t, name, other, "uuid v7 harus unik")
}

func TestLoadBatchInheritsProvidedTitle(t *testing.T) {
	t.Chdir(t.TempDir())
	base, err := Load(newCommand(t, "--type", TypeReceipt, "--seed", "9", "--title", "Judul Global"))
	require.NoError(t, err)

	path := writeManifest(t, "jobs:\n  - pages: 1\n")
	jobs, err := LoadBatch(path, base)
	require.NoError(t, err)
	require.Len(t, jobs, 1)
	assert.Equal(t, "Judul Global", jobs[0].Title)
}
