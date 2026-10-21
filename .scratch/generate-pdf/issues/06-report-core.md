# 06 — Report core: dispatch, penyimpanan, verifikasi

Status: ready-for-agent
Blocked by: 01, 05, 09, 10, 11, 12, 13, 14
Type: task

## Scope

`internal/report/report.go`, `internal/report/report_test.go`.

## Requirements

- `func Generate(r data.Report, out string, verify bool) (pages int, size int64, err error)`:
  1. `newDocument(r)` (07) untuk instance maroto yang sudah berisi config, header, dan footer.
  2. dispatch eksplisit `switch r.Type` ke `AddInvoice`, `AddLetter`, `AddReportDocument`, `AddReceipt`, `AddDelivery`, `AddCertificate`; `default` → `fmt.Errorf("jenis dokumen tidak didukung: %q", r.Type)`.
  3. `Generate()` maroto → `doc.Save(out)`.
  4. `os.Stat(out)` untuk ukuran; `Verify(out)` untuk jumlah halaman (selalu, tidak hanya saat flag).
  5. Bila `verify` dan `Verify` gagal → error menyebut path; file boleh tetap ada tetapi generate dinyatakan gagal.
- `func Verify(path string) (pages int, err error)`: buka file, `api.Validate(rs, model.NewDefaultConfiguration())`, lalu `api.PageCountFile(path)`; error dibungkus `%w` dengan path.
- Pakai helper dari `components.go` (07) dan `barcode.go` (08); file ini tidak berisi detail layout per jenis maupun config maroto.

## Acceptance

- Test: jenis tak dikenal → error, tidak ada file tercipta.
- Test: dua generate payload seed 42 jenis sama → byte file identik (`os.ReadFile`), dan `pages >= 1`.
- Test: `Verify` pada file yang dipotong 50 byte terakhir → error.
- Test: `--verify` gagal tidak mengembalikan pages 0 tanpa error.
- `go test ./internal/report -count=1 -run Core` lulus.

## Out of scope

Layout per jenis (09–14), helper komponen (07), barcode (08).
