// Package cli holds the cobra commands. Minimal stub; replaced by ticket 03.
package cli

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// New returns the pdfgen root command.
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
