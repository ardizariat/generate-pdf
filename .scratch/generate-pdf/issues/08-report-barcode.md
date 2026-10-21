# 08 — Barcode Code128

Status: ready-for-agent
Blocked by: none
Type: task

## Scope

`internal/report/barcode.go`, `internal/report/barcode_test.go`.

## Requirements

- `func BarcodePNG(value string) ([]byte, error)`: `code128.Encode(value)` → `scale.Scale(code, 300, 60)` → `png.Encode` ke `bytes.Buffer`.
- Validasi input: string kosong → error; karakter non-ASCII → error (Code128 hanya ASCII).
- Semua error dibungkus `%w`.
- Tidak mengimpor maroto maupun layout mana pun.

## Acceptance

- Test: `BarcodePNG("INV-123")` → bytes diawali signature PNG dan `image.Decode` berhasil.
- Test: string kosong dan `"é"` → error.
- Test: nilai sama menghasilkan byte identik.
- `go test ./internal/report -count=1 -run Barcode` lulus.

## Out of scope

QR, layout, PDF.
