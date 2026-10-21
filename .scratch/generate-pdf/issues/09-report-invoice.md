# 09 — Layout invoice

Status: ready-for-agent
Blocked by: 05, 07, 08
Type: task

## Scope

`internal/report/invoice.go`, `internal/report/invoice_test.go`.

## Requirements

- `func AddInvoice(m core.Maroto, r data.Report) error` — hanya layout, data dari payload.
- Kop: perusahaan penjual (bold 12) + alamat/email (size 8).
- Meta: `No. Invoice: <DocID>`, `Tanggal: <formatted date>`.
- Dua kolom pihak: Penjual dan Pembeli (nama perusahaan, alamat, email).
- Tabel border `Barang | Qty | Harga | Subtotal` (lebar 5/1/3/3) memakai helper 07, `len(r.InvoiceItems)` baris.
- Perhitungan di `data` (05): subtotal per baris, subtotal total, PPN 11%, total; layout hanya menampilkan `r.InvoiceTotals` (decimal yang sudah diformat).
- Baris ringkasan rata kanan: Subtotal, PPN 11%, Total (bold).
- `signature.NewCol` rata kanan (nama + jabatan penerbit), tinggi 18mm.
- Error dari helper/barcode diteruskan; tidak boleh ada `errcheck` yang diabaikan.

## Acceptance

- Test: payload seed 42 rows 8 → file valid, halaman == 1.
- Test: total yang tampil = `subtotal + PPN` dari payload (uji `data` totals, bukan parsing PDF).
- `go test ./internal/report -count=1 -run Invoice` lulus.

## Out of scope

Diskon, pajak lain, invoice multi-halaman.
