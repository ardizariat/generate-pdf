# 19 — Prompt batch interaktif

Status: done
Blocked by: 17
Type: task

## Scope

`internal/prompt/prompt.go`, `internal/prompt/prompt_test.go`.

## Requirements

- Alur interaktif baru: tanya jumlah file (default 1, minimal 1), lalu satu form per file berisi jenis, jumlah halaman (default 1), seed, path output.
- Nilai global yang `Provided` (tipe/seed/title) menjadi default dan tidak ditanyakan per file; `pages` dan `out` selalu ditanyakan karena berbeda per file. Default `out` = `<type>-<index>-<seed>.pdf`.
- Setelah semua file diisi, tampilkan ringkasan daftar job dan minta konfirmasi sebelum generate; batal = error "input dibatalkan".
- Pemisahan tanggung jawab: fungsi murni pembentuk `[]config.Job` dari jawaban (mudah diuji tanpa form) + fungsi yang menjalankan form.
- `--no-input`/stdin bukan terminal: tidak ada prompt. `--rows` tidak ditanyakan. Tidak ada rendering atau pembacaan konfigurasi global di package ini.

## Acceptance

- Table-driven test untuk pembentuk job: jumlah file 1 dan N, warisan nilai global, penamaan output berurutan, `pages` per file, validasi input jumlah halaman (≥ 1).
- Test pembatalan memetakan `huh.ErrUserAborted` ke error eksplisit.
- `go test ./internal/prompt -count=1` lulus; `go vet ./...` dan `golangci-lint run` bersih.

## Out of scope

Registrasi flag, pemanggilan render, ringkasan hasil generate (ticket 20).

## Comments

- Direvisi setelah review UX: interaksi disederhanakan menjadi dua form (jumlah file + multi-select jenis, lalu jumlah halaman per file). Tidak ada lagi prompt seed/output/judul, ringkasan konfirmasi dihapus, jenis diputar berurutan, dan jenis tanpa tabel dipaksa 1 halaman.

## Comments

- Revisi lanjutan: prompt halaman per file diganti satu input "Jumlah halaman per file" yang dipakai untuk semua file, dan semuanya masuk satu form (jumlah file, multi-select jenis, jumlah halaman). Tidak ada lagi layar kedua.
