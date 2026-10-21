# 17 — Config batch: model job dan manifest YAML

Status: ready-for-agent
Blocked by: none
Type: task

## Scope

`internal/config/config.go`, `internal/config/config_test.go`.

## Requirements

- Tambah key `batch` ke `configKeys`; `Config.Batch` berisi path manifest.
- Tipe `Job` dengan field `Type string`, `Pages int`, `Seed uint64`, `Out`, `Title`.
- `LoadBatch(path string, base Config) ([]Job, error)`: baca YAML bentuk `jobs: [...]`, abaikan field tak dikenal, error jelas bila file tidak ada/rusak atau `jobs` kosong.
- Resolusi nilai per job: field job menang; field kosong mewarisi nilai global (`base`, hasil flag/env/file/default). `pages` default 1; `type` default global; `seed` default global; `title` default per jenis; `out` kosong → `<type>-<index>-<seed>.pdf` dengan index mulai 1.
- `ValidateJobs(jobs []Job) error`: minimal satu job; `pages >= 1`; `pages > 1` hanya untuk `invoice`, `laporan`, `surat-jalan`; `type` salah satu jenis; `out` wajib terisi setelah default. Pesan menyebut index job (mulai 1), nama field, dan nilai.
- Tidak ada akses viper/cobra dari logika job; tidak ada global mutable state.

## Acceptance

- Table-driven test (`t.TempDir()`) untuk: manifest valid, `jobs` kosong, file hilang, YAML rusak, jenis tak dikenal, `pages: 0`, `pages > 1` pada `surat-resmi`/`kuitansi`/`sertifikat`, warisan nilai global vs override job, penamaan output default, field tak dikenal diabaikan.
- `go test ./internal/config -count=1` lulus; `go vet ./...` dan `golangci-lint run` bersih.

## Out of scope

Registrasi flag `--batch` (ticket 20), prompt, rendering, perubahan perilaku flag yang sudah ada.
