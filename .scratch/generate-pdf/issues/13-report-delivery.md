# 13 — Layout surat jalan

Status: ready-for-agent
Blocked by: 05, 07, 08
Type: task

## Scope

`internal/report/delivery.go`, `internal/report/delivery_test.go`.

## Requirements

- `func AddDelivery(m core.Maroto, r data.Report) error`.
- Kop perusahaan pengirim (bold 12) + alamat/telepon.
- Meta: `No. Surat Jalan: SJ-<tahun>-<6 karakter DocID>`, `Tanggal: <date>`.
- Tujuan: `Kepada:` perusahaan + alamat + kota penerima.
- Tabel border `Barang | Qty | Satuan | Keterangan` (lebar 5/1/2/4) memakai helper 07, `len(r.DeliveryItems)` baris.
- Dua blok tanda tangan bersebelahan (kolom 6+6): `Pengirim` dan `Penerima`, masing-masing `signature.NewCol` + nama.
- Biarkan mengalir bila melebihi halaman; jangan `AddPage` manual.

## Acceptance

- Test: seed 42 rows 25 → validasi lolos, halaman == 1.
- Test: rows 200 → halaman ≥ 2.
- `go test ./internal/report -count=1 -run Delivery` lulus.

## Out of scope

Nomor kendaraan, jam berangkat, barcode tambahan.
