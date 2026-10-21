// Package cli holds the cobra commands.
package cli

import (
	"fmt"

	"charm.land/lipgloss/v2"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/ardizariat/generate-pdf/internal/config"
)

const (
	defaultRows = -1
	defaultSeed = 42
)

// typeDescriptions is read-only: short description per document type for `types`.
var typeDescriptions = map[string]string{
	config.TypeInvoice:     "Invoice penjualan dengan tabel barang dan PPN",
	config.TypeLetter:      "Surat resmi dengan kop dan tanda tangan",
	config.TypeReport:      "Laporan tabel yang mengalir multi-halaman",
	config.TypeReceipt:     "Kuitansi pembayaran dengan nominal terbilang",
	config.TypeDelivery:    "Surat jalan pengiriman dengan dua tanda tangan",
	config.TypeCertificate: "Sertifikat landscape dengan QR verifikasi",
}

// New returns the pdfgen root command. It never prints errors or usage itself.
func New(logger *zap.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "pdfgen",
		Short:         "Generate PDF dokumen dummy Indonesia",
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runGenerate(cmd, logger)
		},
	}

	cmd.Flags().StringP("type", "t", config.TypeInvoice, "jenis dokumen: invoice, surat-resmi, laporan, kuitansi, surat-jalan, sertifikat")
	cmd.Flags().StringP("out", "o", "", "path output PDF (default <type>-<seed>.pdf)")
	cmd.Flags().IntP("rows", "r", defaultRows, "jumlah baris tabel (-1 = default per jenis)")
	cmd.Flags().Uint64P("seed", "s", defaultSeed, "seed data; 0 berarti acak")
	cmd.Flags().String("title", "", "judul dokumen (default per jenis)")
	cmd.Flags().String("batch", "", "path manifest YAML berisi daftar job")
	cmd.Flags().Bool("no-input", false, "matikan prompt interaktif")
	cmd.Flags().Bool("force", false, "izinkan menimpa file output yang sudah ada")
	cmd.Flags().Bool("verify", false, "validasi PDF dengan pdfcpu setelah ditulis")
	cmd.Flags().BoolP("verbose", "v", false, "log level info (default warn)")
	cmd.Flags().String("config", "", "path file konfigurasi eksplisit")

	cmd.AddCommand(newTypesCommand())
	return cmd
}

// newTypesCommand prints the supported document types; it does not load config.
func newTypesCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "types",
		Short: "Tampilkan jenis dokumen yang didukung",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			name := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#5fafff"))
			desc := lipgloss.NewStyle().Faint(true)

			out := cmd.OutOrStdout()
			for _, docType := range config.Types() {
				if _, err := fmt.Fprintf(out, "%s  %s\n", name.Render(docType), desc.Render(typeDescriptions[docType])); err != nil {
					return fmt.Errorf("tulis %s: %w", docType, err)
				}
			}
			return nil
		},
	}
}
