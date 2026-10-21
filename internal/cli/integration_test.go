package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ardizariat/generate-pdf/internal/report"
)

// docTypeCase describes one full CLI generate run and its expected page count.
type docTypeCase struct {
	name     string
	args     []string
	minPages int
	maxPages int // 0 = no upper bound
}

// docTypeCases covers every supported type with the defaults from the spec plus
// the multi-page and single-page row variants named there.
func docTypeCases() []docTypeCase {
	return []docTypeCase{
		{name: "invoice default rows", args: []string{"-t", "invoice"}, minPages: 1, maxPages: 1},
		{name: "surat-resmi", args: []string{"-t", "surat-resmi"}, minPages: 1, maxPages: 1},
		{name: "laporan default rows", args: []string{"-t", "laporan"}, minPages: 3},
		{name: "laporan few rows", args: []string{"-t", "laporan", "-r", "5"}, minPages: 1, maxPages: 1},
		{name: "kuitansi", args: []string{"-t", "kuitansi"}, minPages: 1, maxPages: 1},
		{name: "surat-jalan default rows", args: []string{"-t", "surat-jalan"}, minPages: 1, maxPages: 1},
		{name: "surat-jalan many rows", args: []string{"-t", "surat-jalan", "-r", "200"}, minPages: 2},
		{name: "sertifikat", args: []string{"-t", "sertifikat"}, minPages: 1, maxPages: 1},
	}
}

func TestGenerateEveryTypeThroughCLI(t *testing.T) {
	for _, tc := range docTypeCases() {
		t.Run(tc.name, func(t *testing.T) {
			out := filepath.Join(t.TempDir(), "out.pdf")
			args := append(append([]string{}, tc.args...), "-o", out, "--no-input", "-s", "42")

			_, err := execute(t, args...)
			require.NoError(t, err)

			info, err := os.Stat(out)
			require.NoError(t, err)
			assert.Positive(t, info.Size())

			pages, err := report.Verify(out)
			require.NoError(t, err)
			assert.GreaterOrEqual(t, pages, tc.minPages)
			if tc.maxPages > 0 {
				assert.LessOrEqual(t, pages, tc.maxPages)
			}
		})
	}
}

func TestGenerateIsDeterministicThroughCLI(t *testing.T) {
	for _, docType := range []string{"invoice", "laporan"} {
		t.Run(docType, func(t *testing.T) {
			dir := t.TempDir()
			first := filepath.Join(dir, "first.pdf")
			second := filepath.Join(dir, "second.pdf")

			for _, out := range []string{first, second} {
				_, err := execute(t, "-t", docType, "-s", "42", "-o", out, "--no-input")
				require.NoError(t, err)
			}

			a, err := os.ReadFile(first)
			require.NoError(t, err)
			b, err := os.ReadFile(second)
			require.NoError(t, err)
			assert.True(t, strings.HasPrefix(string(a), "%PDF-"))
			assert.Equal(t, a, b, "dua run CLI dengan input sama harus byte-identik")
		})
	}
}
