# 03 — CLI: root, flag, alur generate

Status: ready-for-agent
Blocked by: 01, 02, 05, 06
Type: task

## Scope

`internal/cli/root.go`, `internal/cli/generate.go`, `internal/cli/cli_test.go`.

## Requirements

- `func New(logger *zap.Logger) *cobra.Command` mengembalikan root command `pdfgen` dengan `SilenceErrors`/`SilenceUsage` true; tidak mencetak apa pun sendiri.
- Flag (lihat spec): `-t/--type`, `-o/--out`, `-r/--rows` (sentinel -1), `-s/--seed`, `--title`, `--no-input`, `--force`, `--verify`, `-v/--verbose`, `--config`. Setiap flag punya deskripsi singkat.
- Subcommand `types` mencetak enam jenis + deskripsi, di-styling `lipgloss/v2`; tidak memuat config.
- `root.go` hanya: definisi command, flag, registrasi subcommand, dan pemanggilan `runGenerate(cmd, logger)`.
- `generate.go` menjalankan urutan eksplisit (tanpa rendering PDF di sini):
  1. `config.Load(cmd)`
  2. `config.Validate(cfg)`
  3. `prompt.Fill(cmd, cfg)` bila TTY dan `!cfg.NoInput`
  4. `config.Validate(cfg)` ulang (final validation)
  5. cek overwrite: `os.Stat(cfg.Out)` ada dan `!cfg.Force` → error yang menyebut path dan `--force`
  6. `data.Generate(cfg.Type, cfg.Rows, cfg.Seed)` → payload dokumen
  7. `report.Generate(payload, cfg.Out, cfg.Verify)` → `(pages, size)`
  8. log ringkasan level `info` (type, out, rows, pages, bytes, duration)
- Setiap error dibungkus `fmt.Errorf("...: %w", err)`; tidak ada `os.Exit`, tidak ada log error di lapisan ini.
- Flag tak dikenal → error dari cobra, bukan panic.

## Acceptance

- `go run ./cmd/pdfgen types` mencetak enam jenis, exit 0.
- `go run ./cmd/pdfgen -t salah` exit 1 dengan satu pesan error (tidak dobel).
- Test alur memakai command yang dijalankan dengan argumen, `t.TempDir()` untuk output; test overwrite: tanpa `--force` gagal, dengan `--force` sukses.
- `go test ./internal/cli -count=1` lulus.

## Out of scope

Prompt (04), data (05), report (06–14).
