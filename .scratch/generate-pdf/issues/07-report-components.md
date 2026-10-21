# 07 — Komponen bersama: header, footer, tabel

Status: ready-for-agent
Blocked by: 05, 08
Type: task

## Scope

`internal/report/components.go`, `internal/report/components_test.go`.

## Requirements

- `func newDocument(r data.Report) (core.Maroto, error)`: config maroto A4 portrait, margin atas/kiri/kanan 10, bawah 15, `WithCreationDate(r.Date)`, `WithTitle/Author("generate-pdf")/WithCreator/WithSubject(r.Type)`, `WithPageNumber(props.PageNumber{Pattern: "Halaman {current} dari {total}", Place: props.Bottom, Size: 7})`; sertifikat → `WithOrientation(orientation.Horizontal)`. Lalu pasang header dan footer. Inilah satu-satunya tempat config maroto dibuat.
- `addHeader(m core.Maroto, r data.Report) error`: baris judul (bold, size 12), baris berisi kolom QR (`code.NewQrCol`, payload `https://example.com/report?doc=<DocID>&seed=<seed>&type=<type>`), kolom barcode dari `BarcodePNG(r.DocID)` (08), kolom ID dokumen (size 7), baris garis. Total tinggi ≤ 22mm. Error dari barcode diteruskan.
- `addFooter(m core.Maroto) error`: teks `generate-pdf` kiri, size 7; nomor halaman sudah dari config, jangan duplikasi.
- `addTableHeader(m core.Maroto, titles []string, widths []int)`: satu baris border penuh, teks bold, latar abu muda.
- `addTableRow(m core.Maroto, cells []string, widths []int)`: baris border penuh, tinggi 8mm, align per kolom bila diperlukan.
- `addParagraph(m core.Maroto, text string)`: `text.NewAutoRow` size 9.
- Semua helper menerima parameter eksplisit; tidak membaca config global, tidak mengimpor `cli`/`cobra`/`viper`/`huh`.
- `widths` wajib berjumlah 12 (grid maroto); helper mengembalikan error bila tidak.

## Acceptance

- Test: `widths` yang jumlahnya bukan 12 → error.
- Test: jumlah judul != jumlah lebar → error.
- Test integrasi ringan: header+footer pada dokumen kosong → `pdfcpu` validasi lolos.
- `go test ./internal/report -count=1 -run Components` lulus.

## Out of scope

Isi layout per jenis.
