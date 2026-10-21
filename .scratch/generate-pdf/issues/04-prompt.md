# 04 — Prompt interaktif huh v2

Status: ready-for-agent
Blocked by: 01
Type: task

## Scope

`internal/prompt/prompt.go`, `internal/prompt/prompt_test.go`.

## Requirements

- `func Fill(cmd *cobra.Command, cfg config.Config) (config.Config, error)`; hanya berisi interaksi, tanpa viper, tanpa data, tanpa report.
- Aktif bila `os.Stdin.Stat()` menunjukkan `os.ModeCharDevice` dan `!cfg.NoInput`; selain itu mengembalikan `cfg` apa adanya.
- Form 1: `huh.NewSelect[string]` jenis dokumen (label Indonesia, value slug).
- Form 2: `huh.NewInput` rows (prefill default jenis terpilih, validasi angka ≥ 0), seed (validasi angka), path output (prefill `<type>-<seed>.pdf`). Field dengan `cfg.Provided(key)` true dilewati.
- Jenis tanpa tabel (rows default 0) tidak menanyakan rows.
- Batalkan user: `errors.Is(err, huh.ErrUserAborted)` → `fmt.Errorf("input dibatalkan")`, bukan panic.
- Nilai hasil prompt diisi ke `cfg` dan dikembalikan; validasi akhir tetap tugas `internal/cli`.

## Acceptance

- Test `Provided`-aware skipping tanpa membuka TTY (extract fungsi murni `fieldsToAsk(cmd, cfg) []string`, uji langsung).
- Test non-TTY (stdin pipe) mengembalikan cfg tanpa perubahan.
- Test pembatalan dipetakan ke error, bukan panic.
- `go test ./internal/prompt -count=1` lulus.

## Out of scope

Flag baru, validasi konfigurasi (milik 01), rendering.
