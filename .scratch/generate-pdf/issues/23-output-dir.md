# 23 — Folder output yang bisa dikustom

Status: done
Blocked by: none
Type: task

## Scope

`internal/config` (key/field `out-dir`, resolusi path), `internal/cli` (flag dan pembuatan folder), `internal/prompt` (nama default), dokumentasi.

## Requirements

- Default output pindah ke folder `output/` relatif cwd.
- Flag baru `--out-dir` (env `GENPDF_OUT_DIR`, key config `out-dir`) menjadi folder dasar untuk nama default dan path relatif.
- `--out` tetap berarti path file persis dan menang atas `--out-dir`; path absolut dipakai apa adanya.
- `out` relatif pada manifest batch diselesaikan terhadap `--out-dir`.
- Folder output dan folder induk path output dibuat otomatis dengan `MkdirAll`; error hanya bila path ada tapi bukan direktori atau gagal dibuat.
- Mode interaktif memakai `--out-dir`; `--out` tetap tidak dipakai di mode itu.

## Acceptance

- `pdfgen -t invoice --no-input` menulis `output/invoice-42.pdf` dan membuat folder `output`.
- `pdfgen --out-dir hasil ...` menulis ke `hasil/...`; `--out /abs/x.pdf` tetap persis.
- Manifest dengan `out: sub/a.pdf` menulis `--out-dir/sub/a.pdf`.
- Semua test lama lulus setelah ekspektasi path default diperbarui.
