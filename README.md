# generate-pdf (pdfgen)

CLI Go untuk menghasilkan PDF **dummy** bagi enam jenis dokumen kantor Indonesia.
Hasilnya valid (dapat diverifikasi dengan pdfcpu) dan isinya deterministik untuk seed
yang sama, sehingga cocok dipakai sebagai fixture test, showcase layout, maupun contoh
data untuk stakeholder. Nama file default tetap unik per generate (timestamp + UUID).
Seluruh data dan tanda tangan bersifat fiktif.

## Persyaratan dan instalasi

- Go 1.27.3 atau lebih baru.
- Instalasi ke `$GOBIN`/`$GOPATH/bin` (pastikan direktori tersebut ada di `PATH`):

```bash
go install ./cmd/pdfgen
```

Alternatif tanpa instalasi: `make build` lalu jalankan `./pdfgen`.

## Cara pakai

### Interaktif

Jalankan di terminal tanpa `--no-input`:

```bash
pdfgen
```

Satu form menanyakan **jumlah file**, **jenis dokumen** (boleh lebih dari satu), dan
**satu jumlah halaman** yang dipakai untuk semua file. Setelah form selesai generate
langsung berjalan dan ringkasan dicetak ke stdout.

### Satu dokumen (noninteraktif)

```bash
pdfgen --no-input                                  # invoice default
pdfgen -t laporan -s 7 --no-input --verify
pdfgen -t sertifikat --title "Sertifikat Pelatihan" --no-input
pdfgen -t kuitansi --out output/kuitansi.pdf
pdfgen -t laporan --out-dir hasil --no-input       # tulis ke hasil/
```

Pada mode ini `--rows` mengatur jumlah baris tabel; tidak ada opsi jumlah halaman pasti.

### Batch

```bash
pdfgen --batch jobs.yaml --no-input
```

Manifest dan aturan halaman dijelaskan di [Batch dan target halaman](#batch-dan-target-halaman).

### Perintah lain

```bash
pdfgen types    # tampilkan jenis dokumen yang didukung
```

Bila stdin bukan terminal (mis. pipe/CI), `pdfgen` otomatis memakai mode satu dokumen
noninteraktif. Nama file default selalu
`<out-dir>/<type>_<timestamp UTC>_<uuid v7>.pdf`, dan folder output dibuat otomatis.
Menimpa file yang sudah ada hanya dengan `--force`.

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
| `-o, --out` | `<out-dir>/<type>_<UTC>_<uuid v7>.pdf` | Path output PDF (persis) |
| `--out-dir` | `output` | Folder dasar output default dan path relatif; dibuat otomatis |
| `-r, --rows` | `-1` (efektif per jenis) | Jumlah baris tabel mode satu dokumen; diabaikan bila jenis tidak bertabel |
| `-s, --seed` | `42` | Seed data (`uint64`); `0` berarti acak |
| `--title` | per jenis | Judul dokumen |
| `--batch` | kosong | Path manifest YAML berisi daftar job; bila diisi, mode batch dan prompt dilewati |
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
menulis file konfigurasi. Nama key sama dengan nama flag:

```yaml
# generate-pdf.yaml
type: laporan
rows: 20
seed: 7
out-dir: hasil
verify: true
```

Environment variable memakai prefix `GENPDF_` dengan `-` diganti `_`:
`GENPDF_TYPE`, `GENPDF_OUT`, `GENPDF_OUT_DIR`, `GENPDF_ROWS`, `GENPDF_SEED`,
`GENPDF_TITLE`, `GENPDF_BATCH`, `GENPDF_NO_INPUT`, `GENPDF_FORCE`, `GENPDF_VERIFY`,
`GENPDF_VERBOSE`, `GENPDF_CONFIG`.

Config divalidasi setelah semua sumber digabung dan sekali lagi setelah prompt:
`type` wajib salah satu jenis di atas dan `rows` minimal `0`. Bila `--out` tidak
diberikan, path dihitung dari `--out-dir` dan nama file otomatis.

## Mode interaktif dan noninteraktif

- **Interaktif** aktif bila stdin adalah terminal dan `--no-input` tidak diberikan.
  Satu form saja: jumlah file, jenis dokumen (multi-select, boleh lebih dari satu), dan
  satu jumlah halaman yang dipakai untuk semua file. Jenis dipakai berurutan secara
  berputar bila file lebih banyak daripada jenis.
  - Bila `--type`/env/file sudah menentukan jenis, field multi-select dilewati dan
    semua file memakai jenis itu.
  - `--rows` dan `--out` tidak dipakai di mode ini; lokasi diatur `--out-dir` dan nama
    file selalu otomatis.
  - Jumlah halaman > 1 hanya berlaku untuk `invoice`, `laporan`, dan `surat-jalan`;
    jenis tanpa tabel selalu 1 halaman.
  - Path output, seed, dan judul tidak ditanyakan: output otomatis
    `<out-dir>/<type>_<UTC>_<uuid v7>.pdf`, seed dari flag/env/default, judul default
    per jenis atau `--title`.
  - Tidak ada konfirmasi; generate langsung berjalan dan ringkasan dicetak setelahnya.
    Membatalkan form menghasilkan error `input dibatalkan` dan exit code `1`.
- **Noninteraktif** berlaku bila `--no-input` diberikan, stdin bukan terminal, atau
  `--batch` digunakan. Semua nilai diambil dari flag, env, file konfigurasi, atau
  default; tidak ada pertanyaan dan tidak ada pembacaan stdin.

## Batch dan target halaman

`--batch jobs.yaml` menghasilkan setiap entri manifest sebagai satu PDF:

```yaml
jobs:
  - type: invoice
    pages: 2
    seed: 7
    out: invoice-a.pdf        # relatif -> <out-dir>/invoice-a.pdf
  - type: laporan
    pages: 3
    seed: 8                   # out kosong -> nama unik otomatis
```

- `jobs` wajib berisi minimal satu entri; field yang tidak dikenal diabaikan. Batch
  tidak pernah menampilkan prompt.
- `--rows` tidak dipakai di mode batch; isi tabel diatur `pages`.
- Field job (`type`, `pages`, `seed`, `out`, `title`) menimpa nilai global; field yang
  kosong mewarisi flag/env/file/default. `title` kosong memakai judul default per jenis
  kecuali `--title` diberikan. `out` kosong menjadi
  `<out-dir>/<type>_<timestamp UTC>_<uuid v7>.pdf` (timestamp `YYYYMMDDTHHMMSSZ`,
  UUID dibuat per file). `out` relatif diselesaikan terhadap `--out-dir`; path
  absolut dipakai apa adanya.
- `pages` default `1` dan harus tepat. Isi tabel diisi semaksimal mungkin untuk jumlah
  halaman itu; jumlah halaman akhir dihitung pdfcpu dan target yang tidak tercapai
  menggagalkan job tersebut.
- `pages` > 1 hanya untuk `invoice`, `laporan`, dan `surat-jalan`. Jenis lain wajib
  1 halaman.
- Seluruh job divalidasi sebelum render: satu entri tidak valid menggagalkan batch
  tanpa menulis file apa pun.
- Satu job gagal tidak menghentikan job lain. Ringkasan per file dicetak ke stdout,
  mis. `ok     output/laporan_20261021T141530Z_0192f3a4-....pdf  3 halaman  178697 bytes`,
  dan exit code `1` bila ada kegagalan.
- `--force`, `--verify`, dan `-v` berlaku untuk semua job.
- Selama generate (batch maupun satu file), progress ditulis ke stderr: spinner
  `[2/3] file.pdf · target 2 halaman` bila stdin dan stderr terminal; bila bukan
  terminal, satu baris teks. `-v` mematikan progress karena log zap sudah menampilkan
  tahapan. stdout tetap hanya berisi ringkasan.

## Perilaku output

- Nama default `<out-dir>/<type>_<timestamp UTC>_<uuid v7>.pdf`, mis.
  `output/invoice_20261021T141530Z_0192f3a4-....pdf`. `--out-dir` default `output`.
  Timestamp memakai `YYYYMMDDTHHMMSSZ`; UUID v7 dibuat per file.
- `--out` menang atas `--out-dir` untuk nama file; isinya dipakai persis, termasuk
  path absolut.
- Folder output dibuat otomatis (`MkdirAll`), termasuk folder induk dari `--out` atau
  `out` relatif di manifest. Error hanya bila path ada tapi bukan direktori atau gagal
  dibuat.
- Bila file tujuan sudah ada, generate **ditolak** dengan error yang menyebut path dan
  cara mengizinkannya. `--force` mengizinkan penimpaan. Pada batch, hanya job itu yang
  gagal dan job lain tetap berjalan.
- `--verify` memvalidasi hasil dengan pdfcpu setelah file ditulis; verifikasi gagal
  berarti generate gagal. Jumlah halaman selalu dihitung (lewat pdfcpu) sebagai bagian
  dari proses generate, terlepas dari `--verify`.
- Generate baru dinyatakan berhasil setelah file selesai ditulis.

## Determinisme

- Seed, jenis, jumlah baris, dan judul yang sama menghasilkan PDF dengan isi dan
  metadata identik. Nama file **tidak** deterministik karena memakai timestamp dan
  UUID v7; pakai `--out`/`out` untuk path yang dapat diprediksi.
- Tanggal dokumen diturunkan dari seed (`time.Unix(seed, 0).UTC()`) dan ID dokumen
  adalah UUID v5 dari `<type>-<seed>`. Pada batch, setiap job memakai seed-nya sendiri.
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
