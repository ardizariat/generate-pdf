# 21 — Dokumentasi batch dan Definition of Done

Status: ready-for-agent
Blocked by: 20
Type: task

## Scope

`README.md`, test integrasi akhir di `internal/cli`, `go.mod`/`go.sum` hasil `go mod tidy`.

## Requirements

- README: bagian batch berisi contoh manifest (`jobs:`), tabel flag menambah `--batch`, mode interaktif menjelaskan jumlah file + halaman per file, semantik target halaman (tepat, diverifikasi pdfcpu, hanya `invoice`/`laporan`/`surat-jalan` boleh > 1), perilaku lanjut-setelah-gagal + exit 1, determinisme per job.
- README tidak boleh menyebut perilaku yang tidak ada di kode.
- Test integrasi akhir: manifest tiga job halaman 1/2/3 → tiga file, halaman tepat, `--verify` lolos; satu job gagal (path sudah ada tanpa `--force`) → job lain tetap jalan dan exit non-nil; `--batch` dengan manifest tidak valid tidak menulis file.
- `go mod tidy` tidak mengubah `go.mod`/`go.sum`.

## Acceptance

Lampirkan output:

```
gofmt -w .
go vet ./...
go test ./... -count=1
go test -race ./... -count=1
golangci-lint run
go mod tidy && git diff --exit-code go.mod go.sum
```

plus generate batch nyata ke `/tmp/opencode/pdfcheck/batch/` dengan manifest halaman
1/2/3, `ls -l`, dan jumlah halaman tiap PDF.

## Out of scope

Fitur baru, refactor kosmetik, perubahan perilaku mode satu file.
