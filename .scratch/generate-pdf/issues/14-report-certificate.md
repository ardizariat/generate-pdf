# 14 — Layout sertifikat landscape

Status: ready-for-agent
Blocked by: 05, 07, 08
Type: task

## Scope

`internal/report/certificate.go`, `internal/report/certificate_test.go`.

## Requirements

- `func AddCertificate(m core.Maroto, r data.Report) error`; orientasi landscape sudah diatur di 06, jangan diubah di sini.
- Bingkai: `props.Cell{BorderType: border.Full, BorderThickness: 1}` pada satu baris pembungkus, atau dua baris garis atas/bawah — pilih yang paling sederhana dan pastikan tidak menabrak header.
- `SERTIFIKAT` bold center size 18.
- Nama penerima size 28 bold center dari payload.
- Paragraf penghargaan (center, size 10) yang menyatakan data fiktif.
- Baris bawah: tanggal + `signature.NewCol` pemberi di kiri, `code.NewQrCol` di kanan dengan payload `https://example.com/cert/<DocID>`.
- Tanpa tabel.

## Acceptance

- Test: seed 42 → validasi lolos, halaman == 1.
- Test: dimensi halaman landscape; bila API pdfcpu untuk dimensi tidak tersedia, assert halaman == 1 dengan komentar `shortcut:` yang menyebut alasan.
- `go test ./internal/report -count=1 -run Certificate` lulus.

## Out of scope

Font TTF kustom, latar gambar, ornamen.
