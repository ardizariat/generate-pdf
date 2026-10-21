# 02 — Entry point dan logging

Status: ready-for-agent
Blocked by: 01
Type: task

## Scope

`cmd/pdfgen/main.go`, `internal/logging/logging.go`, `internal/logging/logging_test.go`.

## Requirements

- `internal/logging`: `func New(verbose bool) *zap.Logger` — console encoder, stderr, level `warn` default / `info` bila verbose. `Sync()` dipanggil sebelum exit (abaikan error sync stderr dengan `//nolint` bila perlu, atau `errors.Is`).
- Tanpa variabel package-level untuk logger: nilai dikembalikan dan diteruskan sebagai parameter.
- `cmd/pdfgen/main.go`: buat `cli.New(logger)` → `Execute()`; error → `logger.Error("generate gagal", zap.Error(err))` lalu `os.Exit(1)`; sukses → `os.Exit(0)` implisit. Tidak ada logika PDF, config, atau flag di file ini.
- Tidak ada `os.Exit` di luar `main.go`.

## Acceptance

- `go build ./...` lulus dan `go vet ./...` bersih.
- Test logger: `New(false)` tidak mengeluarkan entri `info`; `New(true)` mengeluarkannya (pakai `zaptest/observer` atau `zapcore` buffer).
- `go test ./internal/logging -count=1` lulus.

## Out of scope

Command, flag, prompt, report.
