# 12 — Layout kuitansi

Status: ready-for-agent
Blocked by: 05, 07, 08
Type: task

## Scope

`internal/report/receipt.go`, `internal/report/receipt_test.go`.

## Requirements

- `func AddReceipt(m core.Maroto, r data.Report) error`.
- Kotak berbingkai (`props.Cell{BorderType: border.Full, BorderThickness: 0.4}`) setinggi area isi berisi:
  - `KUITANSI` bold center size 14.
  - `No: <8 karakter pertama DocID>`.
  - `Telah terima dari` + perusahaan.
  - `Uang sejumlah` + `Rp<n>` bold, di bawahnya `data.Terbilang(n) + " rupiah"`.
  - `Untuk pembayaran` + keperluan.
  - Tanggal + kota.
- `signature.NewCol` kanan bawah (nama + jabatan penerima).
- Nominal dari payload (`r.Amount`, kelipatan 1.000 antara 100.000 dan 50.000.000) — jangan menghitung ulang di layout.
- Tanpa tabel.

## Acceptance

- Test: seed 42 → validasi lolos, halaman == 1.
- Test: `Terbilang(r.Amount)` dipanggil pada payload seed 42 dan hasilnya tidak kosong serta tidak memuat digit.
- `go test ./internal/report -count=1 -run Receipt` lulus.

## Out of scope

Materai, kuitansi multi-item, non-rupiah.
