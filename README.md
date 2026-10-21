# generate-pdf (pdfgen)

CLI Go untuk menghasilkan PDF **dummy** bagi enam jenis dokumen kantor Indonesia.
Hasilnya valid (dapat diverifikasi dengan pdfcpu), dan deterministik untuk seed yang
sama, sehingga cocok dipakai sebagai fixture test, showcase layout, maupun contoh data
untuk stakeholder. Seluruh data dan tanda tangan bersifat fiktif.

## Persyaratan dan instalasi

- Go 1.27.3 atau lebih baru.
- Instalasi ke `$GOBIN`/`$GOPATH/bin` (pastikan direktori tersebut ada di `PATH`):

```bash
go install ./cmd/pdfgen
```

Alternatif tanpa instalasi: `make build` lalu jalankan `./pdfgen`.

## Cara pakai

```bash
pdfgen [flags]       # hasilkan dokumen
pdfgen types         # tampilkan jenis dokumen yang didukung
```

Tanpa argumen dan tanpa prompt (mis. stdin bukan terminal, atau dengan `--no-input`),
`pdfgen` membuat invoice default ke `invoice-42.pdf`. Contoh:

```bash
pdfgen -t laporan -s 7 --no-input --verify
pdfgen -t kuitansi --out /tmp/opencode/kuitansi.pdf --force -v
pdfgen -t sertifikat --title "Sertifikat Pelatihan" --no-input
```

## Jenis dokumen

| `--type` | Default `--rows` | Default `--title` | Isi |
|---|---|---|---|
| `invoice` | 8 | Invoice | Kop, penjual/pembeli, tabel barang, subtotal/PPN 11%/total, tanda tangan |
| `surat-resmi` | 0 (tanpa tabel) | Surat Resmi | Kop, nomor, perihal, tujuan, tiga paragraf, tanda tangan |
| `laporan` | 70 | Laporan | Dua paragraf dan tabel yang mengalir multi-halaman |
| `kuitansi` | 0 (tanpa tabel) | Kuitansi | Kotak berbingkai, nominal angka dan terbilang, tanda tangan |
| `surat-jalan` | 25 | Surat Jalan | Tujuan, tabel barang, tanda tangan pengirim dan penerima |
| `sertifikat` | 0 (tanpa tabel) | Sertifikat | Landscape, nama besar, paragraf penghargaan, QR verifikasi |

`--rows` hanya berpengaruh pada `invoice`, `laporan`, dan `surat-jalan`; untuk jenis
lain nilainya diabaikan.

## Flag

| Flag | Default | Keterangan |
|---|---|---|
| `-t, --type` | `invoice` | Jenis dokumen: `invoice`, `surat-resmi`, `laporan`, `kuitansi`, `surat-jalan`, `sertifikat` |
| `-o, --out` | `<type>-<seed>.pdf` | Path output PDF |
| `-r, --rows` | per jenis (lihat tabel di atas) | Jumlah baris tabel; diabaikan bila jenis tidak bertabel |
| `-s, --seed` | `42` | Seed data (`uint64`); `0` berarti acak |
| `--title` | per jenis | Judul dokumen |
| `--no-input` | `false` | Matikan prompt interaktif |
| `--force` | `false` | Izinkan menimpa file output yang sudah ada |
| `--verify` | `false` | Validasi PDF dengan pdfcpu setelah file ditulis |
| `-v, --verbose` | `false` | Log level `info` (default `warn`) |
| `--config` | kosong | Path file konfigurasi eksplisit |

Exit code: `0` sukses, `1` untuk kegagalan apa pun (config tidak valid, file tujuan
sudah ada, prompt dibatalkan, render/verifikasi gagal).

## Konfigurasi

Urutan prioritas (tertinggi ke terendah):

1. Flag CLI
2. Environment variable `GENPDF_*`
3. File konfigurasi
4. Default (umum dan per jenis)

Sumber file, berurutan sesuai prioritas:

- `--config <path>` untuk file eksplisit, atau
- environment variable `GENPDF_CONFIG=<path>`, atau
- `generate-pdf.yaml` di direktori kerja saat ini.

File konfigurasi bersifat opsional; bila tidak ada, bukan error. Aplikasi tidak pernah
menulis file konfigurasi. Nama key sama dengan nama flag.

Environment variable memakai prefix `GENPDF_` dengan `-` diganti `_`:
`GENPDF_TYPE`, `GENPDF_OUT`, `GENPDF_ROWS`, `GENPDF_SEED`, `GENPDF_TITLE`,
`GENPDF_NO_INPUT`, `GENPDF_FORCE`, `GENPDF_VERIFY`, `GENPDF_VERBOSE`,
`GENPDF_CONFIG`.

Config divalidasi setelah semua sumber digabung dan sekali lagi setelah prompt:
`type` wajib salah satu jenis di atas, `rows` minimal `0`, `out` wajib terisi.

## Mode interaktif dan noninteraktif

- **Interaktif** aktif bila stdin adalah terminal dan `--no-input` tidak diberikan.
  Form menanyakan jenis dokumen, jumlah baris (hanya untuk jenis bertabel), seed, dan
  path output. Field yang sudah diisi lewat flag, env, atau file konfigurasi tidak
  ditanyakan lagi. Membatalkan form menghasilkan error `input dibatalkan` dan exit
  code `1`.
- **Noninteraktif** berlaku bila `--no-input` diberikan atau stdin bukan terminal.
  Semua nilai diambil dari flag, env, file konfigurasi, atau default; tidak ada
  pertanyaan dan tidak ada pembacaan stdin.

## Perilaku output

- Nama default `<type>-<seed>.pdf` di direktori kerja.
- Bila file tujuan sudah ada, generate **ditolak** dengan error yang menyebut path dan
  cara mengizinkannya. `--force` mengizinkan penimpaan.
- Direktori tujuan harus sudah ada; direktori yang tidak ada menghasilkan error, bukan
  panic.
- `--verify` memvalidasi hasil dengan pdfcpu setelah file ditulis; verifikasi gagal
  berarti generate gagal. Jumlah halaman selalu dihitung (lewat pdfcpu) sebagai bagian
  dari proses generate, terlepas dari `--verify`.
- Generate baru dinyatakan berhasil setelah file selesai ditulis.

## Determinisme

- Seed, jenis, jumlah baris, dan judul yang sama menghasilkan PDF dengan isi dan
  metadata identik. Tanggal dokumen diturunkan dari seed
  (`time.Unix(seed, 0).UTC()`) dan ID dokumen adalah UUID v5 dari `<type>-<seed>`.
- `--seed 0` berarti data acak; output tidak deterministik.
- Default seed adalah `42`, sehingga pemakaian tanpa flag tetap dapat direproduksi.

## Menjalankan test

```bash
make build   # go build ./...
make test    # go test ./... -count=1
make fmt     # gofmt -w .
make lint    # golangci-lint run (govet, errcheck, staticcheck, gofmt, unused)
make race    # go test -race ./... -count=1
make run     # go run ./cmd/pdfgen
```

`make lint` membutuhkan `golangci-lint` v2 terpasang. Bila belum tersedia, jalankan
pemeriksaan setara secara manual:

```bash
go vet ./...
gofmt -l .
```
