# generate-pdf (pdfgen)

`pdfgen` membuat file PDF **contoh (dummy)** untuk enam jenis dokumen kantor:
invoice, surat resmi, laporan, kuitansi, surat jalan, dan sertifikat.

Semua isinya data karangan, bukan data asli, jadi aman dipakai untuk demo, latihan,
tampilan ke klien, atau menguji aplikasi lain yang membutuhkan contoh PDF. File bisa
dibuka di pembaca PDF biasa dan lolos pemeriksaan `pdfcpu`.

Dengan `seed` yang sama, isi PDF selalu sama, sehingga hasilnya bisa diulang untuk
pengujian. Nama file tetap berbeda tiap kali generate supaya tidak saling menimpa.

## Persyaratan dan instalasi

- Go 1.27.3 atau lebih baru.
- Instalasi ke `$GOBIN`/`$GOPATH/bin` (pastikan direktori tersebut ada di `PATH`):

```bash
go install ./cmd/pdfgen
```

Alternatif tanpa instalasi: `make build` lalu jalankan `./pdfgen`.

## Cara pakai

Pilih cara yang sesuai kebutuhan:

| Kebutuhan | Cara |
|---|---|
| Coba-coba dan ingin dituntun form | jalankan `pdfgen` di terminal |
| Satu PDF dengan pengaturan tertentu | `pdfgen --no-input` + flag |
| Banyak PDF sekaligus dengan daftar pengaturan | `pdfgen --batch jobs.yaml` |

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

Bagian ini dipakai untuk membuat **banyak PDF sekaligus** tanpa menjawab form satu per
satu. Cocok saat Anda butuh puluhan atau ratusan file contoh untuk menguji aplikasi,
mengisi folder demo, atau menyiapkan bahan presentasi.

Caranya: tulis "daftar pesanan" PDF di satu file YAML bernama `jobs.yaml` (disebut
**manifest**), lalu jalankan satu perintah. Setiap entri di daftar itu menjadi satu file
PDF.

Contoh: Anda ingin tiga PDF — invoice 2 halaman, laporan 3 halaman, dan surat jalan
1 halaman.

```yaml
jobs:
  - type: invoice       # jenis dokumen
    pages: 2            # PDF ini harus tepat 2 halaman
    seed: 7             # penentu isi; seed sama = isi sama
    out: invoice-a.pdf  # nama file; relatif -> <out-dir>/invoice-a.pdf
  - type: laporan
    pages: 3
    seed: 8             # out tidak ditulis -> nama unik otomatis
  - type: surat-jalan   # jenis tanpa tabel, wajib 1 halaman
```

Jalankan dengan:

```bash
pdfgen --batch jobs.yaml --no-input
```

### Arti tiap kolom

| Kolom | Arti sederhana |
|---|---|
| `type` | Jenis dokumen yang dibuat. Boleh juga diisi lewat `--type` sebagai default. |
| `pages` | Jumlah halaman yang **harus** dihasilkan. Default 1. |
| `seed` | Angka penentu isi. Seed sama menghasilkan isi yang sama; `0` berarti acak. |
| `out` | Nama/lokasi file. Kosong berarti nama unik otomatis. |
| `title` | Judul dokumen. Kosong berarti `--title` (bila diberikan) atau judul default sesuai jenis. |

### Yang perlu diketahui

- **`pages` harus tepat.** Untuk dokumen bertabel, program mengisi baris tabel
  sebanyak mungkin sampai pas dengan jumlah halaman itu, lalu menghitung halaman
  dengan `pdfcpu`. Kalau target tidak bisa dicapai, file itu gagal dibuat.
- **Halaman lebih dari 1 hanya untuk dokumen bertabel**: `invoice`, `laporan`, dan
  `surat-jalan`. `surat-resmi`, `kuitansi`, dan `sertifikat` selalu 1 halaman.
- **`--rows` tidak dipakai di mode batch.** Jumlah isi tabel diatur oleh `pages`.
- **Semua entri diperiksa dulu sebelum menulis file.** Satu entri salah membuat seluruh
  batch berhenti tanpa menghasilkan file apa pun.
- **Satu file gagal tidak menghentikan sisanya.** Program lanjut ke entri berikutnya,
  lalu mencetak ringkasan dan exit code `1` bila ada yang gagal.
- **`type`, `seed`, dan `title` yang kosong mewarisi pengaturan global** dari
  flag/env/file konfigurasi. `out` kosong menjadi
  `<out-dir>/<type>_<timestamp UTC>_<uuid v7>.pdf`; `out` relatif diselesaikan terhadap
  `--out-dir`, sedangkan path absolut dipakai apa adanya.
- **`--force`, `--verify`, dan `-v` berlaku untuk semua entri.** Tanpa `--force`, file
  yang sudah ada tidak ditimpa.
- **Kolom yang tidak dikenal diabaikan**, dan batch tidak pernah menampilkan form.

### Cara membaca hasilnya

Selama proses berjalan, progress muncul di stderr (spinner `[2/3] nama-file.pdf ·
target 2 halaman`, atau satu baris teks bila bukan terminal). Setelah selesai, stdout
berisi ringkasan per file, misalnya:

```
ok     output/laporan_20261021T141530Z_0192f3a4-....pdf  3 halaman  178697 bytes
ringkasan: 3 sukses, 0 gagal
```

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
