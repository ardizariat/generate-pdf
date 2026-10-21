# Prompt implementasi ticket — opencode

Cara pakai: buka ticket di `.scratch/generate-pdf/issues/`, jalankan satu sesi opencode baru per ticket, copas blok yang sesuai. Urutan yang memenuhi semua `Blocked by`:

```
01 → 02 → 05 → 08 → 04 → 07 → 09 → 10 → 11 → 12 → 13 → 14 → 06 → 03 → 15 → 16
```

Aturan yang berlaku untuk semua blok:

- Ticket adalah sumber kebenaran ruang lingkup. `AGENTS.md` dan `.scratch/generate-pdf/spec.md` adalah aturan yang tidak boleh dilanggar.
- Satu sesi = satu ticket. Jangan mengerjakan ticket berikutnya, jangan refactor file di luar scope.
- Jangan menambah dependency, flag, package, atau file baru di luar yang disebut ticket.
- `go.mod` hanya boleh berubah karena `go mod tidy`, bukan karena dependency baru.
- Setiap error dibungkus `fmt.Errorf("...: %w", err)`. Tanpa `panic`, tanpa `os.Exit` di luar `cmd/pdfgen/main.go`, tanpa global mutable state.
- Berhenti dan bertanya bila: butuh mengubah spec/AGENTS, butuh dependency baru, acceptance ticket tidak bisa dicapai tanpa mengubah ticket lain, atau menemukan kontradiksi antar ticket.

---

## Ticket 01 — config

```
Implement ticket `.scratch/generate-pdf/issues/01-config.md` di repo /home/ardi/Development/golang/generate-pdf.

Baca dulu: AGENTS.md, .scratch/generate-pdf/spec.md, lalu ticket tersebut. Ticket menentukan ruang lingkup.
Scope tulis: internal/config/config.go, internal/config/config_test.go. Jangan sentuh file lain.
Prasyarat: belum ada kode di repo selain go.mod/go.sum. Bila scope ini sudah berisi kode, laporkan sebelum menimpa.

Patuhi AGENTS.md §4, §6, §8, §11: tanpa global mutable, viper hanya di package ini, validasi dipisah dari pembacaan, error dibungkus %w, test table-driven.
Verifikasi wajib: go build ./... && go vet ./... && go test ./internal/config -count=1
Selesai bila seluruh acceptance di ticket lulus dan perintah di atas hijau.

Laporkan: file yang dibuat, output perintah, penyimpangan dari ticket, dan konfirmasi bahwa ticket 02 siap dijalankan. Jangan mengerjakan ticket lain.
```

---

## Ticket 02 — entry point dan logging

```
Implement ticket `.scratch/generate-pdf/issues/02-entrypoint-logging.md` di repo /home/ardi/Development/golang/generate-pdf.

Baca dulu: AGENTS.md, .scratch/generate-pdf/spec.md, ticket tersebut.
Scope tulis: cmd/pdfgen/main.go, internal/logging/logging.go, internal/logging/logging_test.go. Jangan sentuh file lain.
Prasyarat: ticket 01 selesai (package internal/config ada dan build hijau). Bila belum, berhenti dan laporkan.

Patuhi AGENTS.md §4.5, §5, §8: logger dibuat eksplisit dan diteruskan, tanpa variabel global, os.Exit hanya di main.go, error dicatat hanya di batas aplikasi.
Verifikasi wajib: go build ./... && go vet ./... && go test ./internal/logging -count=1
Selesai bila acceptance ticket lulus dan main.go hanya berisi: buat logger, panggil cli.New/Execute, tangani error + exit code.

Laporkan: file yang dibuat, output perintah, penyimpangan, dan konfirmasi ticket 05 boleh dijalankan paralel. Jangan mengerjakan ticket lain.
```

---

## Ticket 05 — data dan terbilang

```
Implement ticket `.scratch/generate-pdf/issues/05-data.md` di repo /home/ardi/Development/golang/generate-pdf.

Baca dulu: AGENTS.md, .scratch/generate-pdf/spec.md, ticket tersebut.
Scope tulis: internal/data/indonesia.go, internal/data/generator.go, internal/data/terbilang.go, dan test-nya. Jangan sentuh file lain.
Prasyarat: tidak ada. Bila internal/report sudah ada, pastikan package data tidak mengimpornya.

Patuhi AGENTS.md §3.5, §4.5, §9, §11: slice referensi read-only, satu *gofakeit.Faker per pemanggilan (tanpa faker global), tanpa filesystem/time sistem, terbilang diuji table-driven termasuk 0 dan negatif.
Verifikasi wajib: go build ./... && go vet ./... && go test ./internal/data -count=1
Selesai bila acceptance ticket lulus, termasuk test determinisme seed dan test Terbilang.

Laporkan: file yang dibuat, output perintah, penyimpangan. Jangan mengerjakan ticket lain.
```

---

## Ticket 08 — barcode

```
Implement ticket `.scratch/generate-pdf/issues/08-report-barcode.md` di repo /home/ardi/Development/golang/generate-pdf.

Baca dulu: AGENTS.md, .scratch/generate-pdf/spec.md, ticket tersebut.
Scope tulis: internal/report/barcode.go, internal/report/barcode_test.go. Jangan sentuh file lain.
Prasyarat: tidak ada. Bila package internal/report belum ada, buat file ini saja tanpa menambah file lain.

Patuhi AGENTS.md §10.3, §8.2, §11: error code128/scaling/png tidak boleh diabaikan, input kosong atau non-ASCII menghasilkan error, tidak mengimpor maroto.
Verifikasi wajib: go build ./... && go vet ./... && go test ./internal/report -count=1 -run Barcode
Selesai bila acceptance ticket lulus dan hasil byte identik untuk input sama.

Laporkan: file yang dibuat, output perintah, penyimpangan. Jangan mengerjakan ticket lain.
```

---

## Ticket 04 — prompt interaktif

```
Implement ticket `.scratch/generate-pdf/issues/04-prompt.md` di repo /home/ardi/Development/golang/generate-pdf.

Baca dulu: AGENTS.md, .scratch/generate-pdf/spec.md, ticket tersebut.
Scope tulis: internal/prompt/prompt.go, internal/prompt/prompt_test.go. Jangan sentuh file lain.
Prasyarat: ticket 01 selesai (config.Config dan Provided tersedia). Bila belum, berhenti dan laporkan.

Patuhi AGENTS.md §3.2, §4.5, §7: package prompt tidak mengimpor viper/cobra selain tipe perintah yang dioper, tanpa logika data/report, pembatalan user dikembalikan sebagai error.
Verifikasi wajib: go build ./... && go vet ./... && go test ./internal/prompt -count=1
Selesai bila acceptance ticket lulus (skipping via Provided, non-TTY, pembatalan) tanpa membuka TTY.

Laporkan: file yang dibuat, output perintah, penyimpangan. Jangan mengerjakan ticket lain.
```

---

## Ticket 07 — komponen bersama

```
Implement ticket `.scratch/generate-pdf/issues/07-report-components.md` di repo /home/ardi/Development/golang/generate-pdf.

Baca dulu: AGENTS.md, .scratch/generate-pdf/spec.md, ticket tersebut.
Scope tulis: internal/report/components.go, internal/report/components_test.go. Jangan sentuh file lain.
Prasyarat: ticket 05 (data.Report) dan ticket 08 (BarcodePNG) selesai. Bila belum, berhenti dan laporkan.

Patuhi AGENTS.md §10.2, §10.4: newDocument adalah satu-satunya tempat config maroto dibuat; helper menerima parameter eksplisit, tanpa impor cli/cobra/viper/huh; error dari barcode diteruskan.
Verifikasi wajib: go build ./... && go vet ./... && go test ./internal/report -count=1 -run Components
Selesai bila acceptance ticket lulus, termasuk validasi pdfcpu untuk dokumen kosong ber-header/footer.

Laporkan: file yang dibuat, output perintah, penyimpangan, dan konfirmasi ticket 09–14 boleh dijalankan. Jangan mengerjakan ticket lain.
```

---

## Ticket 09 — invoice

```
Implement ticket `.scratch/generate-pdf/issues/09-report-invoice.md` di repo /home/ardi/Development/golang/generate-pdf.

Baca dulu: AGENTS.md, .scratch/generate-pdf/spec.md, ticket tersebut.
Scope tulis: internal/report/invoice.go, internal/report/invoice_test.go. Jangan sentuh file lain.
Prasyarat: ticket 05, 07, 08 selesai dan go build hijau. Bila belum, berhenti dan laporkan.

Patuhi AGENTS.md §10.4, §10.5, §11: layout hanya memakai data dari payload, aritmetika uang di package data (bukan di layout), error helper diteruskan, test memakai t.TempDir().
Verifikasi wajib: go build ./... && go vet ./... && go test ./internal/report -count=1 -run Invoice
Selesai bila acceptance ticket lulus: validasi pdfcpu lolos dan halaman == 1 untuk seed 42 rows 8.

Laporkan: file yang dibuat, output perintah, penyimpangan. Jangan mengerjakan ticket lain.
```

---

## Ticket 10 — surat resmi

```
Implement ticket `.scratch/generate-pdf/issues/10-report-letter.md` di repo /home/ardi/Development/golang/generate-pdf.

Baca dulu: AGENTS.md, .scratch/generate-pdf/spec.md, ticket tersebut.
Scope tulis: internal/report/letter.go, internal/report/letter_test.go. Jangan sentuh file lain.
Prasyarat: ticket 05, 07, 08 selesai. Bila belum, berhenti dan laporkan.

Patuhi AGENTS.md §10.4, §11: tiga paragraf dari payload, tanpa tabel, tanda tangan memakai komponen signature, error diteruskan.
Verifikasi wajib: go build ./... && go vet ./... && go test ./internal/report -count=1 -run Letter
Selesai bila acceptance ticket lulus: validasi lolos dan halaman == 1.

Laporkan: file yang dibuat, output perintah, penyimpangan. Jangan mengerjakan ticket lain.
```

---

## Ticket 11 — laporan

```
Implement ticket `.scratch/generate-pdf/issues/11-report-laporan.md` di repo /home/ardi/Development/golang/generate-pdf.

Baca dulu: AGENTS.md, .scratch/generate-pdf/spec.md, ticket tersebut.
Scope tulis: internal/report/report_document.go, internal/report/report_document_test.go. Jangan sentuh file lain.
Prasyarat: ticket 05, 07, 08 selesai. Bila belum, berhenti dan laporkan.

Patuhi AGENTS.md §10.4, §11: page-break alami tanpa AddPage manual, header tabel dari addTableHeader, tanpa footer kedua, test memakai t.TempDir().
Verifikasi wajib: go build ./... && go vet ./... && go test ./internal/report -count=1 -run ReportDocument
Selesai bila acceptance ticket lulus: rows 70 → halaman ≥ 3, rows 5 → halaman == 1.

Laporkan: file yang dibuat, output perintah, penyimpangan. Jangan mengerjakan ticket lain.
```

---

## Ticket 12 — kuitansi

```
Implement ticket `.scratch/generate-pdf/issues/12-report-receipt.md` di repo /home/ardi/Development/golang/generate-pdf.

Baca dulu: AGENTS.md, .scratch/generate-pdf/spec.md, ticket tersebut.
Scope tulis: internal/report/receipt.go, internal/report/receipt_test.go. Jangan sentuh file lain.
Prasyarat: ticket 05, 07, 08 selesai. Bila belum, berhenti dan laporkan.

Patuhi AGENTS.md §10.4, §11: kotak berbingkai, nominal dari payload (jangan hitung ulang di layout), terbilang dari package data, error diteruskan.
Verifikasi wajib: go build ./... && go vet ./... && go test ./internal/report -count=1 -run Receipt
Selesai bila acceptance ticket lulus: validasi lolos, halaman == 1, terbilang tidak memuat digit.

Laporkan: file yang dibuat, output perintah, penyimpangan. Jangan mengerjakan ticket lain.
```

---

## Ticket 13 — surat jalan

```
Implement ticket `.scratch/generate-pdf/issues/13-report-delivery.md` di repo /home/ardi/Development/golang/generate-pdf.

Baca dulu: AGENTS.md, .scratch/generate-pdf/spec.md, ticket tersebut.
Scope tulis: internal/report/delivery.go, internal/report/delivery_test.go. Jangan sentuh file lain.
Prasyarat: ticket 05, 07, 08 selesai. Bila belum, berhenti dan laporkan.

Patuhi AGENTS.md §10.4, §11: tabel 4 kolom via helper, dua blok tanda tangan berdampingan, biarkan mengalir tanpa AddPage manual.
Verifikasi wajib: go build ./... && go vet ./... && go test ./internal/report -count=1 -run Delivery
Selesai bila acceptance ticket lulus: rows 25 → halaman == 1, rows 200 → halaman ≥ 2.

Laporkan: file yang dibuat, output perintah, penyimpangan. Jangan mengerjakan ticket lain.
```

---

## Ticket 14 — sertifikat

```
Implement ticket `.scratch/generate-pdf/issues/14-report-certificate.md` di repo /home/ardi/Development/golang/generate-pdf.

Baca dulu: AGENTS.md, .scratch/generate-pdf/spec.md, ticket tersebut.
Scope tulis: internal/report/certificate.go, internal/report/certificate_test.go. Jangan sentuh file lain.
Prasyarat: ticket 05, 07, 08 selesai. Bila belum, berhenti dan laporkan.

Patuhi AGENTS.md §10.4, §11: jangan mengubah orientasi (sudah landscape dari newDocument), tanpa tabel, QR verifikasi di kanan bawah, tanpa font TTF kustom.
Verifikasi wajib: go build ./... && go vet ./... && go test ./internal/report -count=1 -run Certificate
Selesai bila acceptance ticket lulus: validasi lolos, halaman == 1. Jika cek dimensi tidak tersedia di pdfcpu, tulis komentar `shortcut:` dan laporkan.

Laporkan: file yang dibuat, output perintah, penyimpangan. Jangan mengerjakan ticket lain.
```

---

## Ticket 06 — report core

```
Implement ticket `.scratch/generate-pdf/issues/06-report-core.md` di repo /home/ardi/Development/golang/generate-pdf.

Baca dulu: AGENTS.md, .scratch/generate-pdf/spec.md, ticket tersebut.
Scope tulis: internal/report/report.go, internal/report/report_test.go. Jangan sentuh file lain.
Prasyarat: ticket 09–14 selesai (keenam fungsi Add* ada dan build hijau). Bila belum, berhenti dan laporkan.

Patuhi AGENTS.md §10.1, §10.5, §11: dispatch eksplisit tanpa fallback diam-diam, config maroto hanya lewat newDocument, Verify memakai pdfcpu, error dibungkus %w.
Verifikasi wajib: go build ./... && go vet ./... && go test ./internal/report -count=1
Selesai bila acceptance ticket lulus: jenis tak dikenal error tanpa file, dua generate seed sama byte-identik, Verify menolak file rusak.

Laporkan: file yang dibuat, output perintah, penyimpangan, dan konfirmasi ticket 03 siap. Jangan mengerjakan ticket lain.
```

---

## Ticket 03 — CLI

```
Implement ticket `.scratch/generate-pdf/issues/03-cli.md` di repo /home/ardi/Development/golang/generate-pdf.

Baca dulu: AGENTS.md, .scratch/generate-pdf/spec.md, ticket tersebut.
Scope tulis: internal/cli/root.go, internal/cli/generate.go, internal/cli/cli_test.go. Jangan sentuh file lain.
Prasyarat: ticket 01, 02, 05, 06 selesai. Bila belum, berhenti dan laporkan.

Patuhi AGENTS.md §5, §8: urutan load → validate → prompt → validate ulang → cek overwrite → data → report → log; tanpa os.Exit; tanpa rendering PDF di root.go; semua error dibungkus %w dan dicatat hanya oleh main.go.
Verifikasi wajib: go build ./... && go vet ./... && go test ./internal/cli -count=1 && go run ./cmd/pdfgen types
Selesai bila acceptance ticket lulus: types exit 0, jenis salah exit 1 tanpa pesan ganda, overwrite ditolak tanpa --force.

Laporkan: file yang dibuat, output perintah, penyimpangan, dan konfirmasi ticket 15/16 siap. Jangan mengerjakan ticket lain.
```

---

## Ticket 15 — README, Makefile, lint

```
Implement ticket `.scratch/generate-pdf/issues/15-docs-tooling.md` di repo /home/ardi/Development/golang/generate-pdf.

Baca dulu: AGENTS.md (§2, §13, §14), .scratch/generate-pdf/spec.md, ticket tersebut, dan kode CLI yang sudah ada.
Scope tulis: README.md, Makefile, .golangci.yml, .gitignore. Jangan sentuh file Go.
Prasyarat: ticket 01 dan 03 selesai (flag dan perilaku CLI final). Bila belum, berhenti dan laporkan.

Patuhi AGENTS.md §14: dokumentasi hanya menyebut perilaku yang benar-benar ada, termasuk default menolak menimpa file dan flag --force, prioritas config, serta seed 0 = acak.
Verifikasi wajib: make build && make test && make fmt && (golangci-lint run || echo "golangci-lint tidak tersedia")
Selesai bila acceptance ticket lulus dan setiap flag di spec tercantum di README.

Laporkan: file yang dibuat, output perintah, penyimpangan, dan konfirmasi ticket 16 siap. Jangan mengerjakan ticket lain.
```

---

## Ticket 16 — integrasi dan DoD

```
Implement ticket `.scratch/generate-pdf/issues/16-integration-dod.md` di repo /home/ardi/Development/golang/generate-pdf.

Baca dulu: AGENTS.md (§11, §13, §15), .scratch/generate-pdf/spec.md, ticket tersebut.
Scope tulis: internal/cli (test integrasi), internal/report (test lintas jenis), go.mod/go.sum hasil tidy. Jangan menambah fitur atau mengubah perilaku flag.
Prasyarat: ticket 03, 04, 06, 09–14, 15 selesai. Bila belum, berhenti dan laporkan.

Patuhi AGENTS.md §11: table-driven, t.TempDir(), tanpa urutan test yang saling bergantung, uji jalur gagal (overwrite, jenis tak dikenal, file rusak).
Verifikasi wajib (lampirkan output apa adanya):
gofmt -w . && go vet ./... && go test ./... -count=1 && go test -race ./... -count=1 && golangci-lint run && go mod tidy
lalu generate keenam jenis ke /tmp/opencode/pdfcheck dengan --no-input --verify -v dan tampilkan ls -l.
Selesai bila semua hijau, enam file valid, halaman sesuai tabel spec, dan go mod tidy idempoten.

Laporkan: file yang diubah, output lengkap perintah, halaman tiap PDF, penyimpangan. Jangan memulai pekerjaan baru setelah ticket ini.
```
