// Package cli holds the cobra commands.
package cli

import (
	"fmt"
	"runtime/debug"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/ardizariat/generate-pdf/internal/config"
	"github.com/ardizariat/generate-pdf/internal/locale"
)

const (
	defaultRows = -1
	defaultSeed = 42
)

// version is set at build time by goreleaser; "dev" for local builds.
var version = "dev"

// New returns the pdfgen root command. It never prints errors or usage itself.
func New(logger *zap.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "pdfgen",
		Short:         "Generate dummy office PDF documents",
		Version:       effectiveVersion(),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runGenerate(cmd, logger)
		},
	}

	cmd.Flags().StringP("type", "t", config.TypeInvoice, "document type: invoice, letter, report, receipt, delivery-note, certificate")
	cmd.PersistentFlags().String("lang", string(locale.English), "output language: en or id")
	cmd.Flags().StringP("out", "o", "", "output PDF path (default <out-dir>/<type>_<UTC time>_<uuid v7>.pdf)")
	cmd.Flags().String("out-dir", config.DefaultOutDir, "base output directory; created automatically when missing")
	cmd.Flags().IntP("rows", "r", defaultRows, "table row count (-1 = per-type default)")
	cmd.Flags().Uint64P("seed", "s", defaultSeed, "data seed; the same seed produces identical output")
	cmd.Flags().String("title", "", "document title (default per type and language)")
	cmd.Flags().String("batch", "", "path to a YAML job manifest")
	cmd.Flags().Bool("no-input", false, "disable the interactive prompt")
	cmd.Flags().Bool("force", false, "allow overwriting an existing output file")
	cmd.Flags().Bool("verify", false, "validate the PDF with pdfcpu after writing")
	cmd.Flags().BoolP("verbose", "v", false, "info-level logging (default warn)")
	cmd.Flags().String("config", "", "explicit configuration file path")

	cmd.AddCommand(newTypesCommand())
	return cmd
}

// effectiveVersion prefers the goreleaser ldflag, then the module version that
// go install records in build info.
func effectiveVersion() string {
	if version != "dev" {
		return strings.TrimPrefix(version, "v")
	}
	info, ok := debug.ReadBuildInfo()
	if !ok || info.Main.Version == "" || info.Main.Version == "(devel)" {
		return version
	}
	return strings.TrimPrefix(info.Main.Version, "v")
}

// newTypesCommand prints the supported document types in the configured language.
func newTypesCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "types",
		Short: "Show the supported document types",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			lang := locale.English
			if cfg, err := config.Load(cmd); err == nil {
				lang = locale.OrDefault(cfg.Language)
			}
			labels := locale.Get(lang)

			name := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#5fafff"))
			desc := lipgloss.NewStyle().Faint(true)

			out := cmd.OutOrStdout()
			for _, docType := range config.Types() {
				if _, err := fmt.Fprintf(out, "%s  %s\n", name.Render(docType), desc.Render(labels.TypeDescriptions[docType])); err != nil {
					return fmt.Errorf("write %s: %w", docType, err)
				}
			}
			return nil
		},
	}
}
