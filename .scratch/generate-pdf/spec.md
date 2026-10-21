# Spec — generate-pdf (dokumen dummy)

Modul: `github.com/ardizariat/generate-pdf`, Go 1.27.3.

Patuh penuh pada `AGENTS.md` (struktur `cmd/` + `internal/`, kualitas, testing, dokumentasi). Jika spec ini dan AGENTS.md berbeda, AGENTS.md menang.

## Tujuan

CLI yang menghasilkan PDF dummy untuk enam jenis dokumen kantor Indonesia. Tiga pemakaian: fixture test (harus valid dan deterministik), showcase layout, dan contoh data untuk stakeholder.

## Struktur

```
cmd/pdfgen/main.go          entry point; bangun logger, panggil cli, exit code
internal/cli/root.go        root command, flag, subcommand types
internal/cli/generate.go    alur generate (load → prompt → validasi → data → report → verify → log)
internal/config/config.go   Config, viper, default, Provided, Validate
internal/prompt/prompt.go   form huh v2
internal/logging/logging.go pembuatan *zap.Logger
internal/data/indonesia.go  data referensi Indonesia
internal/data/generator.go  pembentuk data dummy ber-seed
internal/data/terbilang.go  angka → kata
internal/report/report.go   Generate (dispatch + simpan + Verify)
internal/report/components.go header, footer, tabel, style bersama, pembangunan instance maroto
internal/report/barcode.go  Code128 → PNG bytes
internal/report/invoice.go
internal/report/letter.go
internal/report/report_document.go
internal/report/receipt.go
internal/report/delivery.go
internal/report/certificate.go
README.md, Makefile, .golangci.yml, .gitignore
```

`package data` tidak boleh mengimpor `report`; layout tidak boleh mengimpor `cli`, `viper`, `cobra`, atau `huh`.

## Perintah

```
generate-pdf [flags]
generate-pdf --batch jobs.yaml
generate-pdf types
```

| Flag | Default | Keterangan |
|---|---|---|
| `-t/--type` | `invoice` | `invoice`, `surat-resmi`, `laporan`, `kuitansi`, `surat-jalan`, `sertifikat` |
| `-o/--out` | `<type>-<seed>.pdf` | path output |
| `-r/--rows` | per jenis | invoice 8, laporan 70, surat-jalan 25, sisanya 0 (diabaikan) |
| `-s/--seed` | 42 | `uint64`; 0 berarti acak (determinisme tidak berlaku) |
| `--title` | per jenis | judul dokumen |
| `--batch` | kosong | path manifest YAML berisi daftar job; bila diisi, mode batch |
| `--no-input` | false | matikan prompt |
| `--force` | false | izinkan menimpa file output yang sudah ada |
| `--verify` | false | validasi pdfcpu + hitung halaman setelah tulis |
| `-v/--verbose` | false | level zap `info` (default `warn`) |
| `--config` | kosong | path file config eksplisit |

Exit code: `0` sukses, `1` kegagalan apa pun. Hanya `main.go` yang memanggil `os.Exit`.

## Konfigurasi

Presedensi: **flag > env `GENPDF_*` > file config > default**.

- File: `--config`, atau `generate-pdf.yaml` di cwd; tidak ada bukan error; tidak pernah ditulis.
- Env prefix `GENPDF_`, `-` → `_`. Semua key flag dapat datang dari env/file, termasuk `batch`.
- `config.Load(cmd)` melakukan penggabungan sumber, default per jenis, dan menghitung `Provided(key) bool` untuk flag/env/file sehingga package `prompt` tidak perlu menyentuh viper.
- `config.Validate(cfg)` diekspor dan dipanggil dua kali: setelah load, dan setelah prompt (final validation).
- Tag `validator/v10`: `type` `required,oneof=...`, `rows` `gte=0`, `out` `required`. Pesan menyebut field dan nilai.

## Interaktif

Aktif bila stdin character device dan `--no-input` false. `huh/v2`, satu form:

1. Jumlah file.
2. Multi-select jenis dokumen (boleh lebih dari satu).
3. Satu jumlah halaman per file, dipakai sama untuk semua file.

Jenis dipakai berputar sesuai urutan pilihan bila file lebih banyak daripada jenis.
Jenis tanpa tabel selalu 1 halaman. Output, seed, dan judul tidak ditanyakan: output
`<type>-<index>-<seed>.pdf`, seed dari flag/env/default, judul default per jenis atau
`--title`. Tidak ada konfirmasi; setelah form selesai batch langsung dijalankan.
`--rows` dan `--out` tidak dipakai pada mode interaktif. Pembatalan pengguna
dikembalikan sebagai error eksplisit ("dibatalkan"), bukan panic.

## Batch dan target halaman

Mode batch aktif bila `--batch <path>` diisi. Manifest YAML:

```yaml
jobs:
  - type: invoice
    pages: 2
    seed: 7
    out: out/invoice-a.pdf
    title: Invoice A
  - type: laporan
    pages: 3
```

- `jobs` wajib berisi minimal satu entri; field tidak dikenal bukan error.
- Field job (`type`, `pages`, `seed`, `out`, `title`) menimpa nilai global
  (flag/env/file/default); field kosong mewarisi global. `out` kosong →
  `<type>-<index>-<seed>.pdf`, index mulai 1.
- `pages` wajib ≥ 1, default 1. `pages` > 1 hanya untuk `invoice`, `laporan`,
  `surat-jalan`; jenis lain wajib 1 halaman.
- Seluruh daftar job divalidasi sebelum render; satu job tidak valid menggagalkan
  batch tanpa menulis file apa pun.
- Jumlah halaman harus tepat. Isi diisi semaksimal mungkin untuk jumlah halaman
  itu: jumlah baris dicari secara monoton (render, hitung halaman dengan pdfcpu,
  sesuaikan) sampai tepat. Target tidak tercapai → job gagal.
- Satu job gagal tidak menghentikan job lain. Ringkasan per file dicetak, exit code
  `1` bila ada yang gagal.
- `--force`, `--verify`, `-v` berlaku global untuk semua job. File yang sudah ada
  tanpa `--force` membuat job itu gagal, bukan ditimpa.
- Batch tidak pernah interaktif; manifest harus lengkap tanpa prompt.

## Perilaku output

- `os.Stat` pada path: jika file ada dan `--force` false → error yang menyebut path dan cara mengizinkannya.
- Direktori tujuan tidak ada → error jelas, bukan panic.
- `report.Generate` menulis file lalu, bila `--verify`, memvalidasi dengan pdfcpu; hanya setelah itu generate dinyatakan berhasil.

## Determinisme

- `WithCreationDate(time.Unix(int64(seed), 0).UTC())`.
- ID dokumen `uuid.NewSHA1(uuid.NameSpaceOID, []byte("<type>-<seed>"))`.
- Satu instance `gofakeit.New(seed)` yang diteruskan eksplisit; tanpa faker global.
- Tidak ada `time.Now()` di isi dokumen maupun metadata.
- Setiap job memakai seed-nya sendiri; job dengan input sama menghasilkan byte identik.

## Struktur dokumen

A4 portrait, margin 10mm (bawah 15mm); sertifikat landscape.

- Header berulang: judul, ID dokumen, QR (`https://example.com/report?doc=<uuid>&seed=<seed>&type=<type>`), Code128 barcode dari ID dokumen (dibuat di `barcode.go`), garis pemisah.
- Footer: teks `generate-pdf` kiri, nomor halaman kanan bawah `Halaman {current} dari {total}`.

### Isi per jenis

| Jenis | File | Isi |
|---|---|---|
| invoice | invoice.go | kop, no. invoice, tanggal, penjual & pembeli, tabel `Barang \| Qty \| Harga \| Subtotal`, Subtotal / PPN 11% / Total (decimal), tanda tangan |
| surat-resmi | letter.go | kop, nomor, tanggal, perihal, tujuan, 3 paragraf, tanda tangan |
| laporan | report_document.go | judul, 2 paragraf, tabel `No \| Nama \| Kota \| Email \| Nilai \| Tanggal` mengalir multi-halaman |
| kuitansi | receipt.go | kotak berbingkai, "Telah terima dari", "Uang sejumlah" (angka + terbilang), "Untuk pembayaran", tanggal, tanda tangan |
| surat-jalan | delivery.go | no. surat jalan, tujuan, tabel `Barang \| Qty \| Satuan \| Keterangan`, tanda tangan pengirim & penerima |
| sertifikat | certificate.go | landscape, bingkai, nama besar ≥24pt center, paragraf penghargaan, QR verifikasi |

Paragraf menyatakan dokumen dihasilkan otomatis dan datanya fiktif.

## Data

- `gofakeit/v7` untuk field generik; nama orang, kota, jalan, perusahaan, barang, satuan, jabatan, keperluan dari slice hardcoded (v7 tanpa locale).
- `shopspring/decimal` untuk seluruh aritmetika uang invoice, PPN 11%.
- `terbilang(int64) string`: 0 → "nol", negatif → prefix "minus ", dukungan sampai triliun; tanpa dependency.

## Logging

Zap console encoder ke stderr; `warn` default, `info` dengan `-v`. Satu ringkasan per generate (`type`, `out`, `rows`, `pages`, `bytes`, `duration`). Logger dibuat di `internal/logging` dan diteruskan sebagai parameter — tanpa variabel global. Error dicatat hanya di batas aplikasi (`main.go`), package lain mengembalikan error dengan `%w`.

## Dependency & pemakaian (tidak boleh menganggur)

| Modul | Dipakai untuk |
|---|---|
| `johnfercher/maroto/v2` | layout PDF |
| `brianvoe/gofakeit/v7` | data palsu ber-seed |
| `pdfcpu/pdfcpu` | `report.Verify` + test |
| `boombuler/barcode` | Code128 di `barcode.go` |
| `google/uuid` | ID dokumen (`NewSHA1`) |
| `go-playground/validator/v10` | validasi config |
| `shopspring/decimal` | uang invoice |
| `stretchr/testify` | assertion test |
| `spf13/cobra` | perintah & flag |
| `spf13/viper` | config file/env/flag |
| `go.yaml.in/yaml/v3` | parse manifest batch (`jobs:`) |
| `charm.land/huh/v2` | prompt |
| `charm.land/lipgloss/v2` | styling `types` dan ringkasan |
| `go.uber.org/zap` | logging |

## Ekspektasi halaman

- `laporan` 70 baris → ≥ 3 halaman; 5 baris → 1 halaman.
- `surat-resmi`, `kuitansi`, `sertifikat` → 1 halaman.
- `invoice` rows 8, `surat-jalan` rows 25 → 1 halaman; `surat-jalan` rows 200 → ≥ 2 halaman.

## Dokumentasi & tooling

- `README.md`: tujuan, instalasi, cara pakai (termasuk batch + contoh manifest), jenis dokumen, tabel flag + default (termasuk `--batch`), sumber config + prioritas, mode interaktif/noninteraktif, perilaku output (termasuk `--force`), target halaman, cara test.
- `Makefile`: `build`, `test`, `lint`, `fmt`, `run`.
- `.golangci.yml`: minimal `govet`, `errcheck`, `staticcheck`, `gofmt`.
- `.gitignore`: binary hasil build dan `*.pdf` di root.

## Definition of Done

```
gofmt -w .
go vet ./...
go test ./... -count=1
go test -race ./... -count=1
golangci-lint run
```

Plus: keenam jenis di-generate ke `/tmp/opencode/pdfcheck/` dengan `--no-input --verify -v`, semuanya lolos validasi, halaman sesuai tabel, dan `--force` salah satu file membuktikan perilaku overwrite.

Plus batch: manifest tiga job (`pages` 1, 2, 3) menghasilkan tiga file dengan jumlah
halaman tepat sesuai target, semua lolos `--verify`, dan satu job gagal tidak
menghentikan job lain (exit 1, ringkasan menyebut file yang gagal).

## Non-goals

HTTP server, gambar dari file, tanda tangan berbasis aset, Dockerfile, CI pipeline, flag di luar tabel di atas, generate paralel, override `verify`/`force` per job, dan halaman kosong dekoratif.
