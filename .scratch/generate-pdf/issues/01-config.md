# 01 — Config, viper, dan validasi

Status: ready-for-agent
Blocked by: none
Type: task

## Goal

Satu representasi konfigurasi final, dibaca di satu tempat, tanpa ada package lain yang menyentuh viper.

## Scope

`internal/config/config.go`, `internal/config/config_test.go`.

## Requirements

- `type Config struct` berisi `Type`, `Out`, `Title string`, `Rows int`, `Seed uint64`, `NoInput`, `Verify`, `Verbose`, `Force bool`, `ConfigPath string`, plus state internal untuk `Provided`.
- `func (c Config) Provided(key string) bool` — true bila nilai datang dari flag yang berubah, env `GENPDF_*`, atau file config. Dipakai `internal/prompt`; jangan ekspos viper.
- `func Load(cmd *cobra.Command) (Config, error)`: viper instance lokal, `SetEnvPrefix("GENPDF")`, `SetEnvKeyReplacer("-", "_")`, file dari `--config` atau `generate-pdf.yaml` di cwd (tidak ada = bukan error), `BindPFlag` semua flag.
- Default per jenis dipusatkan di package ini: `rows` (invoice 8, laporan 70, surat-jalan 25, lainnya 0), `title`, dan `out` (`<type>-<seed>.pdf`).
- `func Validate(c Config) error` diekspor (dipakai ulang setelah prompt): `validator/v10`, tag `required,oneof=...` untuk `Type`, `gte=0` untuk `Rows`, `required` untuk `Out`. Seed 0 diizinkan. Pesan error menyebut field + nilai.
- Tidak ada global mutable state; konstanta daftar jenis dan peta default bersifat read-only.

## Acceptance

- `go test ./internal/config -count=1` lulus.
- Test presedensi: flag > env > file > default untuk minimal `type`, `rows`, `out`.
- Test default per jenis dan `out` turunan `<type>-<seed>.pdf`.
- Test `Provided` untuk tiga sumber.
- Test `Validate`: jenis tak dikenal, rows negatif, out kosong — semua error dengan pesan yang menyebut field.
- Test `Load` dengan file config hilang tidak error.

## Out of scope

Tidak ada cobra command, prompt, maroto, atau logging.
