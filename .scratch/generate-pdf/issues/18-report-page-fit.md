# 18 — Report: isi menyesuaikan target halaman

Status: done
Blocked by: none
Type: task

## Scope

`internal/report` (file baru untuk pencarian halaman), test di package yang sama.

## Requirements

- Fungsi baru, mis. `GeneratePages(build func(rows int) (data.Report, error), out string, targetPages int, verify bool) (rows int, pages int, size int64, err error)`.
- Pencarian monoton: jumlah halaman tidak menurun saat rows naik. Cari jumlah baris **maksimum** yang masih menghasilkan tepat `targetPages` halaman, sehingga halaman terisi penuh, bukan halaman kosong tambahan.
- Setiap kandidat dirender ke direktori sementara (`t.TempDir()` di test), jumlah halaman dihitung lewat `Verify` (pdfcpu). File final ditulis hanya setelah target tercapai; file sementara dibersihkan.
- `targetPages < 1` → error validasi. Target tidak tercapai dalam batas atas pencarian → error yang menyebut target dan batas rows; tidak ada file output tertinggal.
- Tidak mengubah perilaku `Generate` yang ada; determinisme tetap (input sama → byte sama).

## Acceptance

- Table-driven test lintas jenis: `invoice`, `laporan`, `surat-jalan` target 1, 2, 3 halaman → jumlah halaman tepat (diverifikasi ulang dengan `Verify`), file ada dan bukan kosong; `surat-resmi`, `kuitansi`, `sertifikat` target 1 → 1 halaman.
- Test jalur gagal: target 0, dan target yang mustahil dalam batas rows → error, pesan menyebut target, tidak ada file output.
- Test determinisme: dua pemanggilan dengan input sama → byte identik.
- `go test ./internal/report -count=1` lulus; `go vet ./...` dan `golangci-lint run` bersih.

## Out of scope

Perubahan layout, halaman kosong dekoratif, generate paralel, perubahan signature `Generate`/`Verify`.
