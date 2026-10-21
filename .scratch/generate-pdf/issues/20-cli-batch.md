# 20 — CLI batch: flag, alur, ringkasan

Status: done
Blocked by: 17, 18, 19
Type: task

## Scope

`internal/cli/root.go`, `internal/cli/generate.go`, file baru alur batch, test di package yang sama.

## Requirements

- Flag baru `--batch <path>` (default kosong) didokumentasikan di help.
- Pemilihan mode: `--batch` diisi → muat dan validasi seluruh job sebelum render; tanpa `--batch` dan interaktif → prompt batch (ticket 19); selain itu mode satu file yang sekarang, tidak berubah.
- Satu job tidak valid menggagalkan batch sebelum file apa pun ditulis.
- Generate job berurutan: cek `os.Stat` + `--force` per job, cari rows lewat `report.GeneratePages` untuk target halaman, tulis file, verifikasi bila `--verify`.
- Satu job gagal tidak menghentikan job lain; setelah semua selesai cetak ringkasan per file (path, halaman, ukuran atau alasan gagal) ke stdout; exit code 1 bila ada kegagalan. Error tetap hanya dicatat di `main.go`.
- `-v` memberi satu log `info` per job (`type`, `out`, `pages`, `bytes`, `duration`) dan ringkasan akhir.
- Determinisme dipertahankan: seed per job, tanpa `time.Now()` di isi dokumen.

## Acceptance

- Test integrasi (`cmd.SetArgs`, `t.TempDir()`): manifest 3 job halaman 1/2/3 → tiga file, jumlah halaman tepat, `report.Verify` lolos.
- Test lanjut-setelah-gagal: satu job path-nya file yang sudah ada tanpa `--force` → job itu gagal, job lain sukses, exit/error non-nil, ringkasan menyebut file gagal.
- Test `--force` menimpa; test manifest tidak valid tidak menulis file apa pun; test mode satu file lama tetap lulus.
- `go test ./internal/cli -count=1` lulus; `go vet ./...` dan `golangci-lint run` bersih.

## Out of scope

Generate paralel, override `--verify`/`--force` per job, tipe dokumen baru, perubahan perilaku flag lama.

## Comments

- Selesai. Ditemukan bug lama: `--title` tidak pernah dipakai karena `data.Generate` selalu memakai judul default. Diperbaiki minimal di CLI (`payload.Title`/`job.Title`), diuji lewat batch dua judul berbeda yang menghasilkan byte berbeda.
