# 05 — Data dummy Indonesia dan terbilang

Status: ready-for-agent
Blocked by: none
Type: task

## Scope

`internal/data/indonesia.go`, `internal/data/generator.go`, `internal/data/terbilang.go`, test untuk generator dan terbilang.

## Requirements

- `indonesia.go`: slice read-only (masing-masing ≥ 12 entri) untuk nama depan, nama belakang, kota, jalan, perusahaan, barang, satuan, jabatan, keperluan pembayaran. Tidak dimodifikasi saat runtime.
- `generator.go`:
  - `func Generate(docType string, rows int, seed uint64) (Report, error)` menghasilkan satu struct payload berisi semua data yang dibutuhkan layout: `Type`, `DocID` (uuid SHA1 dari `docType`+`seed`), `Date time.Time` (dari seed), `Rows`, pihak-pihak (`Party`), `InvoiceItems` + `InvoiceTotals` (Subtotal, PPN, Total sebagai `decimal.Decimal`), `DeliveryItems`, `ReportRows`, `Amount` (nominal kuitansi), serta `Paragraphs []string` (prosa Indonesia yang menyebut "data fiktif").
  - Tipe data per jenis: `ReportRow`, `InvoiceItem` (dengan `decimal.Decimal` untuk harga), `DeliveryItem`, `Party`.
  - Semua randomisasi lewat satu `*gofakeit.Faker` yang dibuat di fungsi ini dari seed; tanpa faker global, tanpa `math/rand` global.
  - Format tanggal `02 Jan 2006`, rupiah `Rp1.234.567`.
  - Jenis tak dikenal → error.
- `terbilang.go`: `func Terbilang(n int64) string` — 0 "nol", negatif prefix "minus ", satuan sampai triliun, lowercase, tanpa kata "rupiah". Tidak mengimpor `report`.

## Acceptance

- Table-driven `Terbilang`: 0, 7, 11, 25, 100, 1000, 1_500_000, 2_000_000_000, 1_000_000_000_000, -250.
- `Generate` seed sama → payload `reflect.DeepEqual`; seed beda → berbeda; jenis tak dikenal → error.
- `Generate` tidak menyentuh filesystem/waktu sistem.
- Tidak ada import `internal/report` (cek dengan `go list -deps` atau review).
- `go test ./internal/data -count=1` lulus.

## Out of scope

Layout, PDF, cobra, viper, zap.
