# 11 — Layout laporan multi-halaman

Status: ready-for-agent
Blocked by: 05, 07, 08
Type: task

## Scope

`internal/report/report_document.go`, `internal/report/report_document_test.go`.

## Requirements

- `func AddReportDocument(m core.Maroto, r data.Report) error`.
- Judul (bold 14) + tanggal, lalu dua paragraf ringkasan (helper `addParagraph`) yang menyatakan data fiktif.
- Tabel 6 kolom border `No | Nama | Kota | Email | Nilai | Tanggal`, lebar 1/3/2/3/2/2, tinggi baris 8mm, header memakai `addTableHeader`.
- `r.Rows` baris dari payload, ditambahkan per baris lewat `addTableRow`; biarkan maroto memutus halaman sendiri (tanpa `AddPage`).
- Jangan menambah footer kedua.

## Acceptance

- Test: seed 42 rows 70 → validasi lolos, halaman ≥ 3.
- Test: rows 5 → halaman == 1.
- Test: header tabel muncul di halaman kedua (indikator: `pages ≥ 3` untuk 70 baris sudah membuktikan page-break; cek visual manual dicatat di komentar test).
- `go test ./internal/report -count=1 -run ReportDocument` lulus.

## Out of scope

Grafik, subtotal per halaman, daftar isi.
