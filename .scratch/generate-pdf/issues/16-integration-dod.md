# 16 — Integrasi, race test, dan Definition of Done

Status: done
Blocked by: 03, 04, 06, 09, 10, 11, 12, 13, 14, 15
Type: task

## Scope

`internal/cli` (test integrasi), `internal/report` (test table lintas jenis), `go.mod`, `go.sum`.

## Requirements

- Test table lintas jenis lewat jalur CLI penuh (`cmd.SetArgs`, `t.TempDir()`): setiap jenis → exit 0, file ada, `report.Verify` lolos, halaman sesuai tabel spec.
- Test determinisme lintas dua jenis (invoice, laporan): dua run CLI input sama → byte file identik.
- Test `types` memuat keenam slug dan cocok dengan daftar di `internal/config`.
- Overwrite: run kedua tanpa `--force` → exit non-zero dan file lama tidak berubah; dengan `--force` → sukses.
- Audit dependency: setiap modul di tabel spec diimpor minimal satu file non-test (kecuali `testify`); `go mod tidy` idempoten (hash `go.mod`/`go.sum` sama sebelum dan sesudah run kedua).

## Acceptance

Lampirkan output:

```
gofmt -w .
go vet ./...
go test ./... -count=1
go test -race ./... -count=1
golangci-lint run
go mod tidy
mkdir -p /tmp/opencode/pdfcheck
for t in invoice surat-resmi laporan kuitansi surat-jalan sertifikat; do
  go run ./cmd/pdfgen -t "$t" -o "/tmp/opencode/pdfcheck/$t.pdf" --no-input --verify -v
done
ls -l /tmp/opencode/pdfcheck
```

## Out of scope

Fitur baru, refactor kosmetik, perubahan perilaku flag.

## Comments

- Selesai. Test baru: `internal/cli/integration_test.go` (tabel lintas jenis lewat CLI penuh, determinisme invoice+laporan), `internal/cli/deps_test.go` (audit dependency non-test + `go mod tidy -diff`), `internal/report/cross_type_test.go` (dispatch/simpan/verify enam jenis). `TestTypesListsAllDocumentTypes` kini membandingkan slug persis dengan `config.Types()`.
- Verifikasi: `gofmt -w .`, `go vet ./...`, `go test ./... -count=1` (121 tes), `go test -race ./... -count=1` (121 tes), `golangci-lint run` 0 issue, `go mod tidy` idempoten (hash `go.mod`/`go.sum` sama pada run kedua).
- Enam jenis di-generate ke `/tmp/opencode/pdfcheck` dengan `--no-input --verify -v`: invoice 1, surat-resmi 1, laporan 3, kuitansi 1, surat-jalan 1, sertifikat 1 halaman; semua lolos validasi pdfcpu.

