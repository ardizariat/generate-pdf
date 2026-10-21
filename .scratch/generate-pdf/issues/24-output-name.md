# 24 — Format nama file default

Status: done
Blocked by: none
Type: task

## Scope

`internal/config` (helper nama default dan pemakaian di Load/LoadBatch), `internal/prompt`, dokumentasi dan test.

## Requirements

- Nama default menjadi `<type>_<timestamp UTC>_<uuid v7>.pdf`, mis. `invoice_20261021T141530Z_0192f3a4-....pdf`.
- Timestamp format `YYYYMMDDTHHMMSSZ`, diambil saat tiap file mulai diproses (bukan sekali per proses).
- UUID v7 baru per file.
- Berlaku untuk semua nama default: single, interaktif, dan manifest dengan `out` kosong.
- `--out` dan `out` eksplisit tetap dipakai persis; path absolut tidak diubah.
- Nama file tidak lagi deterministik; isi PDF tetap deterministik dari seed. Dokumentasikan bahwa path stabil hanya lewat `--out`/`out`.

## Acceptance

- `pdfgen --no-input` menulis satu file `output/<type>_<UTC>_<uuid>.pdf` dan UUID-nya versi 7.
- Batch dengan `out` kosong menghasilkan nama unik per file.
- `--out` dan manifest `out` tidak berubah perilakunya.
- Semua test lulus setelah ekspektasi nama default diganti pola regex.
