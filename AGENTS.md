# generate-pdf

Go module `github.com/ardizariat/generate-pdf`.

# Go Coding Guidelines — PDF Generator CLI

## Agent skills

## 1. Tujuan

Guideline ini menjadi standar pengembangan CLI PDF Generator menggunakan:

- Go
- Cobra untuk CLI
- Viper untuk konfigurasi
- Validator untuk validasi konfigurasi
- Huh v2 untuk form interaktif
- Zap untuk structured logging
- Maroto untuk pembuatan PDF

Tujuan utama adalah menghasilkan kode yang sederhana, mudah dipahami, mudah diuji, dan mudah dikembangkan tanpa abstraksi yang tidak diperlukan.

## 2. Struktur Direktori

Gunakan struktur berikut sebagai acuan:

```text
pdf-generator/
├── cmd/
│   └── pdfgen/
│       └── main.go
│
├── internal/
│   ├── cli/
│   │   ├── root.go
│   │   └── generate.go
│   │
│   ├── config/
│   │   └── config.go
│   │
│   ├── prompt/
│   │   └── prompt.go
│   │
│   ├── logging/
│   │   └── logging.go
│   │
│   ├── data/
│   │   ├── indonesia.go
│   │   ├── generator.go
│   │   └── terbilang.go
│   │
│   └── report/
│       ├── report.go
│       ├── components.go
│       ├── barcode.go
│       ├── invoice.go
│       ├── letter.go
│       ├── report_document.go
│       ├── receipt.go
│       ├── delivery.go
│       ├── certificate.go
│       └── report_test.go
│
├── go.mod
├── go.sum
├── Makefile
├── README.md
├── .gitignore
└── .golangci.yml
```

Prinsip struktur:

- `cmd/` berisi entry point aplikasi.
- `internal/cli/` menangani perintah dan alur CLI.
- `internal/config/` menangani konfigurasi dan validasi.
- `internal/prompt/` menangani interaksi pengguna.
- `internal/logging/` menangani inisialisasi logger.
- `internal/data/` menyediakan data dummy dan fungsi pembentuk data.
- `internal/report/` menangani pembuatan PDF dan layout dokumen.

Jangan membuat package baru hanya untuk memisahkan satu fungsi kecil. Pemisahan harus memberikan manfaat nyata dalam hal tanggung jawab, penggunaan ulang, atau pengujian.

## 3. Prinsip Desain

### 3.1 Simplicity

Pilih implementasi paling sederhana yang memenuhi kebutuhan.

- Hindari abstraksi yang belum dibutuhkan.
- Hindari terlalu banyak lapisan pemanggilan fungsi.
- Jangan membuat framework internal untuk kebutuhan yang dapat diselesaikan dengan fungsi biasa.
- Jangan menambahkan dependency tanpa alasan yang jelas.

### 3.2 Single Responsibility

Setiap package dan fungsi harus memiliki tanggung jawab yang jelas.

- CLI mengatur interaksi dan alur perintah.
- Config membaca dan memvalidasi konfigurasi.
- Prompt menangani input interaktif.
- Logging menyediakan logger.
- Data menghasilkan data dummy.
- Report menyusun dan menghasilkan dokumen PDF.

Jangan mencampur logika CLI, pembentukan data, dan layout PDF dalam satu fungsi.

### 3.3 Explicitness

Buat alur program dan efek samping mudah ditelusuri.

- Teruskan dependency secara eksplisit.
- Kembalikan error kepada pemanggil.
- Hindari global mutable state.
- Hindari konfigurasi tersembunyi yang diakses dari berbagai bagian program.
- Pastikan operasi penulisan file dan penggantian output memiliki perilaku yang jelas.

### 3.4 DRY dan YAGNI

Hindari duplikasi yang benar-benar perlu dihilangkan, tetapi jangan membuat abstraksi terlalu dini.

- Gunakan helper bersama untuk komponen PDF yang digunakan oleh beberapa layout.
- Jangan membuat interface hanya untuk memenuhi pola desain.
- Jangan membangun fitur yang belum diperlukan.
- Jangan menyatukan dua fungsi hanya karena beberapa baris kodenya mirip jika tanggung jawabnya berbeda.

### 3.5 Testability

Logika bisnis harus dapat diuji tanpa harus menjalankan seluruh aplikasi.

- Uji generator data secara terpisah.
- Uji konversi terbilang secara terpisah.
- Uji validasi konfigurasi tanpa membuat PDF.
- Uji layout dan output PDF pada package report.
- Hindari ketergantungan yang tidak perlu pada filesystem, waktu, dan randomness.

## 4. Standar Penulisan Go

### 4.1 Formatting dan penamaan

- Gunakan `gofmt`.
- Ikuti konvensi penamaan Go.
- Gunakan `MixedCaps` atau `mixedCaps`, bukan `snake_case`, untuk identifier Go.
- Gunakan nama yang jelas dan spesifik.
- Hindari singkatan yang tidak umum.
- Berikan komentar pada identifier yang diekspor jika membantu menjelaskan kontrak atau perilakunya.

Contoh:

```go
func GenerateInvoice() error {
    return nil
}
```

Pilih nama fungsi berdasarkan tujuan, bukan detail implementasinya.

### 4.2 Fungsi

- Satu fungsi sebaiknya memiliki satu tujuan utama.
- Gunakan early return untuk mengurangi nested conditionals.
- Hindari fungsi yang terlalu panjang dan menangani banyak tahap proses.
- Pisahkan fungsi ketika ada batas tanggung jawab yang jelas.
- Hindari pemecahan fungsi yang justru membuat alur program sulit diikuti.

### 4.3 Struct dan interface

- Gunakan struct untuk merepresentasikan konfigurasi atau data yang memiliki hubungan logis.
- Gunakan interface ketika diperlukan untuk beberapa implementasi atau untuk memisahkan dependency yang perlu diuji.
- Jangan membuat interface untuk setiap struct secara otomatis.
- Jangan membuat factory atau dependency injection framework tanpa kebutuhan nyata.

### 4.4 Constants dan magic values

- Gunakan konstanta untuk nilai tetap yang memiliki makna bisnis.
- Hindari magic values yang tersebar di banyak fungsi.
- Tempatkan konstanta dekat dengan fungsi atau package yang menggunakannya.
- Jangan membuat konstanta untuk setiap literal sederhana tanpa alasan.

### 4.5 Global state

- Hindari global mutable state.
- Jangan memodifikasi slice data referensi secara langsung.
- Jangan menyimpan konfigurasi runtime dalam variabel global yang dapat berubah.
- Gunakan parameter, struct konfigurasi, atau dependency eksplisit jika nilainya perlu dibagikan.

## 5. CLI dan Cobra

File terkait:

```text
internal/cli/
├── root.go
└── generate.go
```

### `root.go`

Bertanggung jawab untuk:

- Mendefinisikan root command.
- Mendaftarkan flag umum.
- Menentukan konfigurasi command utama.
- Mendaftarkan subcommand.
- Menghubungkan aplikasi dengan konfigurasi dan logging.

### `generate.go`

Bertanggung jawab untuk:

- Menjalankan proses generate PDF.
- Mengambil konfigurasi yang diperlukan.
- Meminta input melalui prompt jika diperlukan.
- Memanggil generator data dan report.
- Mengembalikan error kepada pemanggil.

### Aturan CLI

- Gunakan `RunE` jika command perlu mengembalikan error.
- Validasi input sebelum memulai proses pembuatan dokumen.
- Pisahkan parsing flag dari logika bisnis.
- Gunakan exit code secara konsisten.
- Tampilkan pesan error yang dapat ditindaklanjuti.
- Hindari memanggil `os.Exit` dari fungsi bisnis atau helper.
- Jangan melakukan rendering PDF langsung di `root.go`.

Flag yang didukung harus memiliki nama, default, dan perilaku yang terdokumentasi.

## 6. Konfigurasi dan Validasi

File:

```text
internal/config/config.go
```

### Tanggung jawab

- Mendefinisikan struct `Config`.
- Membaca konfigurasi menggunakan Viper.
- Menentukan default value.
- Mendukung sumber konfigurasi yang memang diperlukan.
- Memvalidasi konfigurasi menggunakan Validator.

### Aturan

- Gunakan satu representasi konfigurasi setelah semua sumber input digabungkan.
- Tentukan prioritas sumber konfigurasi secara eksplisit.
- Hindari mengakses Viper dari seluruh package.
- Baca konfigurasi di satu tempat dan teruskan hasilnya kepada bagian yang membutuhkan.
- Pisahkan pembacaan konfigurasi dari validasi agar keduanya mudah diuji.
- Berikan pesan validasi yang menyebutkan field dan alasan kegagalan.
- Jangan melanjutkan proses generate jika konfigurasi wajib tidak valid.

Contoh alur:

```text
CLI Flags
   ↓
Viper
   ↓
Config
   ↓
Validation
   ↓
Prompt (jika diperlukan)
   ↓
Final Validation
   ↓
Generate PDF
```

Jika prompt melengkapi konfigurasi, lakukan validasi kembali setelah input interaktif diterima.

## 7. Prompt Interaktif

File:

```text
internal/prompt/prompt.go
```

### Aturan

- Gunakan Huh v2 hanya untuk kebutuhan interaksi pengguna.
- Pisahkan definisi form dari logika pembuatan dokumen.
- Gunakan nilai default yang masuk akal.
- Berikan label dan pesan validasi yang jelas.
- Tangani pembatalan pengguna secara eksplisit.
- Jangan menganggap pembatalan sebagai panic atau kegagalan internal.
- Jangan menaruh logika generator data atau layout PDF di dalam prompt.

Jika aplikasi mendukung mode noninteraktif, pastikan seluruh input yang dibutuhkan dapat disediakan melalui flag atau konfigurasi.

## 8. Logging dan Error Handling

File:

```text
internal/logging/logging.go
```

### 8.1 Logging dengan Zap

- Gunakan structured logging.
- Pilih level log sesuai tingkat kepentingan.
- Gunakan field yang relevan, seperti jenis dokumen dan lokasi output.
- Hindari log berulang untuk error yang sama di banyak lapisan.
- Jangan mencatat data sensitif atau isi dokumen yang tidak diperlukan.
- Pastikan logger di-flush sesuai kebutuhan sebelum aplikasi selesai.

Level log:

- `Debug`: detail untuk diagnosis.
- `Info`: tahapan penting proses.
- `Warn`: kondisi yang masih dapat ditangani.
- `Error`: kegagalan operasi.

### 8.2 Error handling

- Jangan mengabaikan error.
- Jangan menggunakan panic untuk kegagalan operasional biasa.
- Bungkus error dengan konteks menggunakan `%w`.
- Gunakan `errors.Is` atau `errors.As` jika perlu mengenali error tertentu.
- Tambahkan konteks yang berguna ketika mengembalikan error.
- Log error pada batas aplikasi jika memang perlu dicatat.

Contoh:

```go
if err != nil {
    return fmt.Errorf("generate invoice: %w", err)
}
```

Hindari mencatat error yang sama pada setiap lapisan. Lapisan internal sebaiknya mengembalikan error, sedangkan batas aplikasi menentukan bagaimana error tersebut ditampilkan dan dicatat.

## 9. Data Dummy dan Fungsi Terbilang

File:

```text
internal/data/
├── indonesia.go
├── generator.go
└── terbilang.go
```

### 9.1 `indonesia.go`

Berisi data referensi yang digunakan untuk menghasilkan dokumen dummy, seperti:

- Nama orang.
- Nama perusahaan.
- Kota dan alamat.
- Satuan barang.
- Data referensi lokal lainnya.

Aturan:

- Jangan mengubah data referensi global ketika menghasilkan data.
- Hindari menduplikasi daftar yang sama.
- Kelompokkan data berdasarkan kegunaan.
- Gunakan tipe yang sederhana dan sesuai kebutuhan.

### 9.2 `generator.go`

Berisi logika pembentukan data dummy, seperti:

- Data pelanggan.
- Data transaksi.
- Baris item invoice.
- Data pengiriman.
- Data laporan.

Aturan:

- Pisahkan pembentukan data dari rendering PDF.
- Gunakan parameter eksplisit untuk nilai yang memengaruhi hasil.
- Pastikan struktur data konsisten.
- Gunakan seed jika hasil acak perlu direproduksi.
- Hindari ketergantungan pada konfigurasi global.

### 9.3 `terbilang.go`

Berisi fungsi konversi angka menjadi teks bahasa Indonesia.

Aturan:

- Tentukan rentang angka yang didukung.
- Tangani nol secara benar.
- Tentukan perilaku untuk angka negatif.
- Pastikan penanganan nilai besar tidak menghasilkan teks yang salah.
- Jangan mencampur fungsi terbilang dengan layout PDF.
- Uji hasil konversi secara terpisah menggunakan table-driven tests.

## 10. Pembuatan PDF dan Maroto

File:

```text
internal/report/
├── report.go
├── components.go
├── barcode.go
├── invoice.go
├── letter.go
├── report_document.go
├── receipt.go
├── delivery.go
└── certificate.go
```

### 10.1 `report.go`

Bertanggung jawab untuk:

- Mengatur alur pembuatan PDF.
- Menentukan layout berdasarkan jenis dokumen.
- Menghubungkan data dengan layout.
- Mengatur proses penyimpanan hasil.
- Mengembalikan error pembuatan dokumen.

Hindari menaruh seluruh detail layout dalam file ini.

Gunakan dispatch yang jelas dan tangani jenis dokumen yang tidak didukung secara eksplisit.

### 10.2 `components.go`

Berisi komponen umum yang digunakan oleh beberapa layout, seperti:

- Header.
- Footer.
- Informasi halaman.
- Komponen teks atau tabel bersama.
- Helper layout yang benar-benar digunakan berulang kali.

Aturan:

- Hindari duplikasi komponen yang identik.
- Jangan memaksakan seluruh layout menggunakan helper yang sama jika kebutuhan tampilannya berbeda.
- Hindari helper dengan terlalu banyak parameter.
- Pastikan komponen bersama tetap mudah dipahami dan digunakan.

### 10.3 `barcode.go`

Bertanggung jawab untuk:

- Menghasilkan barcode.
- Mengembalikan hasil dalam bentuk bytes sesuai kebutuhan.
- Menangani error pembuatan barcode.
- Menjaga agar proses pembuatan barcode tidak terikat pada satu layout dokumen.

Jangan mengabaikan error atau menganggap hasil barcode selalu valid.

### 10.4 File layout dokumen

Setiap file memiliki tanggung jawab layout tersendiri:

| File                 | Jenis dokumen |
| -------------------- | ------------- |
| `invoice.go`         | Invoice       |
| `letter.go`          | Surat resmi   |
| `report_document.go` | Laporan       |
| `receipt.go`         | Kuitansi      |
| `delivery.go`        | Surat jalan   |
| `certificate.go`     | Sertifikat    |

Aturan:

- Pisahkan layout setiap jenis dokumen.
- Gunakan helper bersama jika sesuai.
- Hindari ketergantungan langsung pada Cobra, Viper, dan Huh.
- Jangan membaca konfigurasi global dari dalam layout.
- Terima data yang dibutuhkan secara eksplisit.
- Jangan menggabungkan seluruh layout dalam satu fungsi besar.
- Tangani error yang berasal dari proses rendering.

Setiap layout harus memiliki tanggung jawab yang jelas, tetapi tidak harus dipaksa mengikuti abstraksi yang sama jika kebutuhan dokumennya berbeda.

### 10.5 Output PDF

- Validasi lokasi output.
- Tentukan perilaku jika file tujuan sudah ada.
- Jangan menimpa file secara diam-diam jika perilaku tersebut tidak dinyatakan dalam konfigurasi atau dokumentasi.
- Tangani error dari proses rendering dan penulisan file.
- Gunakan temporary directory untuk pengujian output.
- Pastikan proses penyimpanan selesai sebelum menyatakan generate berhasil.

## 11. Testing

Pengujian ditempatkan sedekat mungkin dengan kode yang diuji.

Contoh:

```text
internal/
├── config/
│   ├── config.go
│   └── config_test.go
├── data/
│   ├── generator.go
│   ├── generator_test.go
│   ├── terbilang.go
│   └── terbilang_test.go
└── report/
    ├── report.go
    ├── report_test.go
    ├── invoice.go
    └── invoice_test.go
```

File pengujian tambahan dapat dibuat sesuai kebutuhan; tidak perlu membuat semuanya sejak awal.

### Aturan testing

- Gunakan table-driven tests untuk variasi input.
- Gunakan `t.TempDir()` untuk pengujian file.
- Jangan bergantung pada urutan eksekusi test.
- Jangan menggunakan randomness yang tidak terkendali.
- Uji kondisi sukses dan gagal.
- Uji error yang dikembalikan.
- Hindari test yang terlalu bergantung pada detail implementasi internal.

Prioritas pengujian:

1. Validasi konfigurasi.
2. Generator data.
3. Konversi terbilang.
4. Dispatch jenis dokumen.
5. Pembuatan barcode.
6. Rendering dan penyimpanan PDF.
7. Penanganan input atau output yang tidak valid.

Untuk PDF, jangan hanya memeriksa apakah fungsi mengembalikan `nil`. Pastikan file berhasil dibuat dan hasilnya masuk akal. Jika tersedia, gunakan parser PDF untuk memverifikasi bahwa file dapat dibaca.

## 12. Dependency dan Package

- Gunakan dependency yang benar-benar diperlukan.
- Hindari menambahkan library untuk fungsi sederhana yang dapat ditangani oleh standard library.
- Jangan membuat package baru hanya untuk satu helper sederhana.
- Jangan membuat interface jika hanya ada satu implementasi dan tidak ada kebutuhan pengujian yang membenarkannya.
- Pastikan arah dependency tetap sederhana.

Arah dependency yang disarankan:

```text
cmd
 ↓
cli ───────► config
 │          prompt
 │          logging
 │
 ├────────► data
 │
 └────────► report
               │
               └── layout dokumen
```

Diagram ini menunjukkan gambaran tanggung jawab, bukan aturan bahwa seluruh package harus saling bergantung secara langsung.

Package `data` sebaiknya tidak bergantung pada `report`, dan layout PDF sebaiknya tidak bergantung pada package CLI.

## 13. Quality Check

Jalankan pemeriksaan berikut sebelum menyelesaikan perubahan:

```bash
gofmt -w .
go vet ./...
go test ./...
go test -race ./...
```

Jika project menggunakan `golangci-lint`:

```bash
golangci-lint run
```

Gunakan pemeriksaan formatting di CI untuk memastikan seluruh file sudah diformat. CI sebaiknya memeriksa hasil formatting tanpa mengubah file secara otomatis.

## 14. Dokumentasi

README minimal menjelaskan:

- Tujuan aplikasi.
- Persyaratan Go dan cara instalasi.
- Cara menjalankan aplikasi.
- Jenis dokumen yang didukung.
- Daftar flag dan nilai default.
- Sumber konfigurasi dan prioritasnya.
- Cara menggunakan mode interaktif dan noninteraktif jika keduanya tersedia.
- Lokasi output dan perilaku terhadap file yang sudah ada.
- Cara menjalankan test.

Jika perilaku CLI, konfigurasi, atau format output berubah, perbarui dokumentasi terkait.

## 15. Definition of Done

Perubahan dianggap selesai jika:

- [ ] Tanggung jawab package dan file tetap jelas.
- [ ] Kode mengikuti konvensi Go dan sudah diformat.
- [ ] Tidak ada abstraksi atau dependency yang tidak diperlukan.
- [ ] Konfigurasi divalidasi sebelum digunakan.
- [ ] Error ditangani dan diteruskan secara konsisten.
- [ ] Logging menggunakan Zap secara konsisten.
- [ ] Data dummy tidak terikat pada rendering PDF.
- [ ] Layout dokumen terpisah berdasarkan jenis dokumen.
- [ ] Komponen bersama digunakan tanpa over-abstraction.
- [ ] Output PDF berhasil dibuat sesuai kebutuhan.
- [ ] Pengujian relevan tersedia dan lulus.
- [ ] `go vet` dan test lulus.
- [ ] Dokumentasi diperbarui jika diperlukan.

## Kesimpulan

Prioritaskan kode yang sederhana, eksplisit, dan mudah diuji.

Gunakan Cobra untuk mengatur alur CLI, Viper dan Validator untuk konfigurasi, Huh untuk interaksi pengguna, Zap untuk logging, package data untuk pembentukan data dummy, serta package report untuk rendering PDF.

Pisahkan setiap layout dokumen, tetapi jangan membuat abstraksi berlebihan. Tambahkan package atau helper baru hanya ketika kompleksitas nyata membuat pemisahan tersebut lebih mudah dipahami dan dirawat.

### Issue tracker

Issues and specs are local markdown files under `.scratch/<feature-slug>/`. See `docs/agents/issue-tracker.md`.

### Triage labels

Default five-role triage vocabulary (`needs-triage`, `needs-info`, `ready-for-agent`, `ready-for-human`, `wontfix`). See `docs/agents/triage-labels.md`.

### Domain docs

Single-context: `CONTEXT.md` plus `docs/adr/` at the repo root. See `docs/agents/domain.md`.
