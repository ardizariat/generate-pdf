// Package cli berisi perintah cobra. Stub minimal; digantikan ticket 03.
package cli

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// New mengembalikan root command pdfgen.
func New(logger *zap.Logger) *cobra.Command {
	return &cobra.Command{
		Use:           "pdfgen",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(*cobra.Command, []string) error {
			return nil
		},
	}
}
