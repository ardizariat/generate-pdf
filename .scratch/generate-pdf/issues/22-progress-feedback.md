# 22 — Progress selama generate

Status: done
Blocked by: none
Type: task

## Scope

`internal/cli` (helper progress baru dan pemakaian di batch serta satu file), test helper.

## Requirements

- Selama generate, tampilkan progress agar user tahu proses berjalan.
- Bila stdin dan stderr adalah terminal dan `-v` tidak diberikan: spinner animasi (`charm.land/huh/v2/spinner`) dengan label `[i/N] <output> · target <pages> halaman`.
- Bila bukan terminal: satu baris teks yang sama ke stderr sebelum generate, tanpa karakter ANSI.
- `-v`: tanpa spinner/baris progress; log zap sudah menampilkan tahapan.
- Progress ditulis ke stderr; stdout tetap hanya ringkasan (`ok`, `gagal`, `ringkasan`) sehingga skrip yang mem-parse tetap aman.
- Berlaku untuk mode satu file dan batch (interaktif maupun manifest).
- Kegagalan job tetap dicatat, spinner berhenti, job berikutnya berjalan seperti sebelumnya.

## Acceptance

- Generate interaktif di terminal menampilkan spinner dan label file yang sedang diproses.
- Pipe/CI tidak menampilkan escape code; hanya baris teks progress.
- `-v` tidak menampilkan progress tambahan.
- Stdout tidak berubah; semua test lama tetap lulus.
