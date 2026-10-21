# 10 — Layout surat resmi

Status: ready-for-agent
Blocked by: 05, 07, 08
Type: task

## Scope

`internal/report/letter.go`, `internal/report/letter_test.go`.

## Requirements

- `func AddLetter(m core.Maroto, r data.Report) error`.
- Kop: nama instansi (bold 14, center), alamat + email (size 8, center), garis pemisah.
- Meta kiri: `Nomor: <nomor>`, `Tanggal: <date>`, `Perihal: <keperluan>`.
- Tujuan: `Kepada Yth.` + nama + jabatan + perusahaan + kota.
- `Dengan hormat,`, lalu tiga paragraf (masing-masing 3–5 kalimat) memakai helper `addParagraph`; salah satunya memuat frasa "data fiktif".
- Penutup `Hormat kami,` + `signature.NewCol` rata kanan.
- Tanpa tabel; `r.Rows` tidak dipakai.

## Acceptance

- Test: payload seed 42 → validasi lolos, halaman == 1.
- Test: `data` menyediakan minimal tiga paragraf dan salah satunya memuat "data fiktif" (uji payload, bukan PDF).
- `go test ./internal/report -count=1 -run Letter` lulus.

## Out of scope

Tembusan, lampiran, kop gambar.
