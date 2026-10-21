# 15 — README, Makefile, lint, gitignore

Status: ready-for-agent
Blocked by: 01, 03
Type: task

## Scope

`README.md`, `Makefile`, `.golangci.yml`, `.gitignore`.

## Requirements

- `README.md` wajib memuat: tujuan, versi Go & instalasi (`go install ./cmd/pdfgen`), cara pakai, enam jenis dokumen, tabel flag + default, sumber config + urutan prioritas, mode interaktif vs noninteraktif, perilaku output (default menolak menimpa, `--force` mengizinkan), cara menjalankan test.
- README menyebut bahwa output deterministik untuk seed sama dan seed 0 berarti acak.
- `Makefile`: `build`, `run`, `test`, `fmt`, `lint`, `race`; target memakai `go build ./...`, `go test ./... -count=1`, `gofmt -w .`, `golangci-lint run`, `go test -race ./... -count=1`.
- `.golangci.yml`: aktifkan minimal `govet`, `errcheck`, `staticcheck`, `gofmt`, `unused`.
- `.gitignore`: `/pdfgen`, `*.pdf` di root, `/dist`.
- Dokumentasi tidak boleh menyebut perilaku yang belum ada di kode.

## Acceptance

- `make build`, `make test`, `make fmt` jalan tanpa error.
- `golangci-lint run` (bila terpasang) melaporkan 0 issue; bila tidak terpasang, catat di README bagian test.
- README menyebut setiap flag di spec; verifikasi dengan pencocokan nama flag.

## Out of scope

CI pipeline, Dockerfile, dokumentasi arsitektur terpisah.
