package cli

import (
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// moduleRoot returns the repository root from this package directory.
func moduleRoot(t *testing.T) string {
	t.Helper()

	root, err := filepath.Abs(filepath.Join("..", ".."))
	require.NoError(t, err)
	return root
}

func runGo(t *testing.T, args ...string) (string, error) {
	t.Helper()

	cmd := exec.Command("go", args...)
	cmd.Dir = moduleRoot(t)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// TestSpecDependenciesAreImported checks the spec's dependency table: each
// module must be imported by non-test code, except testify which is test-only.
func TestSpecDependenciesAreImported(t *testing.T) {
	// go list .Imports excludes test imports, so this proves non-test usage.
	out, err := runGo(t, "list", "-f", `{{join .Imports "\n"}}`, "./...")
	require.NoError(t, err, out)

	imported := make(map[string]bool)
	for _, line := range strings.Split(out, "\n") {
		if line != "" {
			imported[line] = true
		}
	}
	importsModule := func(module string) bool {
		for path := range imported {
			if path == module || strings.HasPrefix(path, module+"/") {
				return true
			}
		}
		return false
	}

	modules := []string{
		"github.com/johnfercher/maroto/v2",
		"github.com/brianvoe/gofakeit/v7",
		"github.com/pdfcpu/pdfcpu",
		"github.com/boombuler/barcode",
		"github.com/google/uuid",
		"github.com/go-playground/validator/v10",
		"github.com/shopspring/decimal",
		"github.com/spf13/cobra",
		"github.com/spf13/viper",
		"charm.land/huh/v2",
		"charm.land/lipgloss/v2",
		"go.uber.org/zap",
	}
	for _, module := range modules {
		assert.True(t, importsModule(module), "%s harus diimpor kode non-test", module)
	}
}

// TestGoModTidyIsIdempotent runs tidy in diff mode, which fails when the module
// files are not already tidy and never rewrites them.
func TestGoModTidyIsIdempotent(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("go mod tidy -diff runs in the Linux CI job")
	}
	out, err := runGo(t, "mod", "tidy", "-diff")
	require.NoError(t, err, out)
}
