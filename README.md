# pdfgen

[![CI](https://github.com/ardizariat/generate-pdf/actions/workflows/ci.yml/badge.svg)](https://github.com/ardizariat/generate-pdf/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go 1.26.1](https://img.shields.io/badge/go-1.26.1-00ADD8?logo=go)](go.mod)

`pdfgen` is a small CLI that generates realistic **dummy PDFs** for six common
Indonesian office documents: invoices, formal letters, reports, receipts, delivery
notes, and certificates.

All data is fictional, every output is a valid PDF, and document content is
reproducible from a seed. Use it for demos, test fixtures, layout checks, or filling
a staging folder without touching real customer data.

[Versi bahasa Indonesia](README.id.md)

![Invoice example](assets/preview.png)

## Features

- Six document types with layouts that look like real office paperwork.
- Three ways to work: one interactive form, flag-driven single files, or a YAML
  manifest that generates many files in one run.
- Exact page targets for table-based documents: content is sized to fill the
  requested number of pages.
- Deterministic content from a seed, with unique file names so runs never collide.
- Page counts measured with [pdfcpu](https://github.com/pdfcpu/pdfcpu); `--verify`
  additionally validates the whole file.
- Configuration through flags, environment variables, or a YAML file.
- No network access; generated documents are entirely local.

## Requirements

- Go 1.26.1 or newer when building from source.
- Prebuilt binaries need no runtime dependencies.

## Installation

Download the archive for your OS from
[Releases](https://github.com/ardizariat/generate-pdf/releases), or install with Go:

```bash
go install github.com/ardizariat/generate-pdf/cmd/pdfgen@latest
```

From a clone:

```bash
make build    # builds ./pdfgen
```

Check the installed version:

```bash
pdfgen --version
```

## Quick start

```bash
# Interactive: one form, then generate
pdfgen

# One file, no prompt
pdfgen --no-input

# Several documents from a manifest
pdfgen --batch jobs.yaml --no-input
```

Default output goes to `./output/` with a unique name:

```
output/invoice_20261021T141530Z_0192f3a4-5b6c-7d8e-9f00-112233445566.pdf
```

## Document types

| `--type` | Default `--rows` | Default `--title` | Content |
|---|---|---|---|
| `invoice` | 8 | Invoice | Letterhead, seller/buyer, item table, subtotal/VAT 11%/total, signature |
| `surat-resmi` | 0 (no table) | Surat Resmi | Letterhead, number, subject, recipient, three paragraphs, signature |
| `laporan` | 70 | Laporan | Two paragraphs and a table that flows across pages |
| `kuitansi` | 0 (no table) | Kuitansi | Framed box, amount in digits and words, signature |
| `surat-jalan` | 25 | Surat Jalan | Recipient, item table, sender and receiver signatures |
| `sertifikat` | 0 (no table) | Sertifikat | Landscape, large name, award paragraph, verification QR |

`--rows` only affects `invoice`, `laporan`, and `surat-jalan`; it is ignored for the
other types.

## Flags

| Flag | Default | Description |
|---|---|---|
| `-t, --type` | `invoice` | Document type: `invoice`, `surat-resmi`, `laporan`, `kuitansi`, `surat-jalan`, `sertifikat` |
| `-o, --out` | `<out-dir>/<type>_<UTC>_<uuid v7>.pdf` | Exact output path |
| `--out-dir` | `output` | Base folder for default names and relative paths; created automatically |
| `-r, --rows` | `-1` (per-type default) | Table rows in single-file mode; ignored for non-table types |
| `-s, --seed` | `42` | Data seed (`uint64`); `0` means random |
| `--title` | per type | Document title |
| `--batch` | empty | YAML manifest of jobs; disables prompts |
| `--no-input` | `false` | Disable the interactive form |
| `--force` | `false` | Allow overwriting an existing output file |
| `--verify` | `false` | Validate the PDF with pdfcpu after writing |
| `-v, --verbose` | `false` | Info-level logs (default is `warn`) |
| `--config` | empty | Explicit configuration file |
| `--version` | | Print the version |

Exit codes: `0` on success, `1` for any failure (invalid configuration, existing
output without `--force`, cancelled prompt, render or verification error).

## Usage

### Interactive

Run `pdfgen` in a terminal without `--no-input`. A single form asks for:

- the number of files,
- the document types (multi-select; several types are cycled in order),
- one page count applied to every file.

Generation starts as soon as the form is submitted; the summary is printed to
stdout. There is no confirmation step. Cancelling the form exits with code `1`.

Notes:

- If `--type` (or an environment variable/config file) already sets the type, the
  multi-select is skipped and every file uses that type.
- `--rows` and `--out` are ignored here. Files go to `--out-dir` with automatic names.
- A page count above 1 applies only to table types; the other types always get one page.

### Single file (non-interactive)

```bash
pdfgen --no-input
pdfgen -t laporan -s 7 --no-input --verify
pdfgen -t sertifikat --title "Sertifikat Pelatihan" --no-input
pdfgen -t kuitansi --out output/kuitansi.pdf
pdfgen -t laporan --out-dir hasil --no-input
```

`--rows` controls the table size in this mode. There is no exact-page option; use the
batch manifest for page targets.

### Batch

```bash
pdfgen --batch jobs.yaml --no-input
```

See [Batch and page targets](#batch-and-page-targets).

### Other commands

```bash
pdfgen types    # list supported document types
```

When stdin is not a terminal (pipes, CI), `pdfgen` automatically uses single-file
non-interactive mode. Output folders are created automatically; existing files are
only replaced with `--force`.

## Configuration

Priority, highest first:

1. CLI flags
2. `GENPDF_*` environment variables
3. Configuration file
4. Defaults (general and per type)

Configuration file lookup:

- `--config <path>`, or
- `GENPDF_CONFIG=<path>`, or
- `generate-pdf.yaml` in the current directory.

The file is optional; a missing file is not an error. `pdfgen` never writes a
configuration file. Keys match the flag names:

```yaml
# generate-pdf.yaml
type: laporan
rows: 20
seed: 7
out-dir: hasil
verify: true
```

Environment variables use the `GENPDF_` prefix with `-` replaced by `_`:
`GENPDF_TYPE`, `GENPDF_OUT`, `GENPDF_OUT_DIR`, `GENPDF_ROWS`, `GENPDF_SEED`,
`GENPDF_TITLE`, `GENPDF_BATCH`, `GENPDF_NO_INPUT`, `GENPDF_FORCE`, `GENPDF_VERIFY`,
`GENPDF_VERBOSE`, `GENPDF_CONFIG`.

Configuration is validated after merging all sources and again after the interactive
form: `type` must be one of the supported types and `rows` must be `0` or greater.
When `--out` is omitted, the path is built from `--out-dir` with an automatic name.

## Batch and page targets

Batch mode creates **many PDFs at once** without answering a form for each one. It is
useful for generating dozens or hundreds of fixtures, filling a demo folder, or
preparing presentation material.

Write a "shopping list" of PDFs in a YAML file called a **manifest**, then run one
command. Every entry under `jobs` becomes one PDF.

Example: an invoice with 2 pages, a report with 3 pages, and a delivery note with
1 page.

```yaml
jobs:
  - type: invoice       # document type
    pages: 2            # this PDF must be exactly 2 pages
    seed: 7             # content seed; same seed = same content
    out: invoice-a.pdf  # relative to <out-dir>/invoice-a.pdf
  - type: laporan
    pages: 3
    seed: 8             # no out -> automatic unique name
  - type: surat-jalan   # no table, must be 1 page
```

Run it with:

```bash
pdfgen --batch jobs.yaml --no-input
```

### Manifest fields

| Field | Meaning |
|---|---|
| `type` | Document type. Can also come from `--type` as a global default. |
| `pages` | Exact page count. Defaults to 1. |
| `seed` | Content seed. Same seed produces the same content; `0` means random. |
| `out` | File name or path. Empty means an automatic unique name. |
| `title` | Document title. Empty uses `--title` when provided, otherwise the per-type default. |

### Rules

- **`pages` must be exact.** For table documents, the program fills as many table
  rows as fit the requested page count, then measures the page count with `pdfcpu`.
  If the target cannot be met, that job fails.
- **More than one page is only supported for table documents**: `invoice`, `laporan`,
  and `surat-jalan`. `surat-resmi`, `kuitansi`, and `sertifikat` always use one page.
- **`--rows` is not used in batch mode.** `pages` controls how much table content is
  generated.
- **All entries are validated before any file is written.** One invalid entry stops
  the whole batch without producing files.
- **A failed job does not stop the others.** The program continues, prints a summary,
  and exits with code `1` if anything failed.
- **Empty `type`, `seed`, and `title` inherit global settings** from flags,
  environment, or the configuration file. An empty `out` becomes
  `<out-dir>/<type>_<timestamp UTC>_<uuid v7>.pdf`; relative `out` paths are resolved
  against `--out-dir`, while absolute paths are used as-is.
- **`--force`, `--verify`, and `-v` apply to every entry.** Without `--force`,
  existing files are not overwritten.
- Unknown fields are ignored, and batch mode never shows a prompt.

### Reading the result

While jobs run, progress is written to stderr (a spinner such as
`[2/3] file.pdf · target 2 halaman` in a terminal, or one line per update otherwise).
Afterwards, stdout contains one summary line per file:

```
ok     output/laporan_20261021T141530Z_0192f3a4-....pdf  3 halaman  178697 bytes
ringkasan: 3 sukses, 0 gagal
```

`-v` disables the progress display because zap logs already show each stage.

## Output behavior

- Default name: `<out-dir>/<type>_<timestamp UTC>_<uuid v7>.pdf`, for example
  `output/invoice_20261021T141530Z_0192f3a4-....pdf`. The timestamp format is
  `YYYYMMDDTHHMMSSZ`; a new UUID v7 is generated per file. `--out-dir` defaults to
  `output`.
- `--out` takes precedence over `--out-dir`; the value is used verbatim, including
  absolute paths.
- Output folders are created automatically (`MkdirAll`), including parent folders of
  `--out` or a relative manifest `out`. A path that exists but is not a directory is
  an error.
- If the target file exists, generation is **rejected** with an error naming the path
  and how to allow it. `--force` allows overwriting. In batch mode only that job
  fails; the others continue.
- `--verify` validates the result with pdfcpu after writing; a failed verification
  fails the job. Page counts are always measured with pdfcpu as part of generation,
  regardless of `--verify`.
- Generation is only reported as successful after the file has been fully written.

## Determinism

- The same seed, type, row count, and title produce identical content and metadata.
  File names are **not** deterministic because they include a timestamp and a UUID
  v7; use `--out` or `out` when you need a predictable path.
- The document date is derived from the seed (`time.Unix(seed, 0).UTC()`) and the
  document ID is a UUID v5 of `<type>-<seed>`. In batch mode every job uses its own
  seed.
- `--seed 0` means random data and produces non-deterministic output.
- The default seed is `42`, so runs without flags are reproducible.

## Examples

Sample outputs generated with fixed seeds:

| Document | File |
|---|---|
| Invoice | [examples/invoice.pdf](examples/invoice.pdf) |
| Formal letter | [examples/surat-resmi.pdf](examples/surat-resmi.pdf) |
| Report | [examples/laporan.pdf](examples/laporan.pdf) |
| Receipt | [examples/kuitansi.pdf](examples/kuitansi.pdf) |
| Delivery note | [examples/surat-jalan.pdf](examples/surat-jalan.pdf) |
| Certificate | [examples/sertifikat.pdf](examples/sertifikat.pdf) |

## Development

```bash
make build   # go build ./...
make test    # go test ./... -count=1
make fmt     # gofmt -w .
make lint    # golangci-lint run (govet, errcheck, staticcheck, gofmt, unused)
make race    # go test -race ./... -count=1
make run     # go run ./cmd/pdfgen
```

`make lint` requires `golangci-lint` v2. Without it, run `go vet ./...` and
`gofmt -l .`.

See [CONTRIBUTING.md](CONTRIBUTING.md) before opening a pull request. Report security
issues privately as described in [SECURITY.md](SECURITY.md).

## License

[MIT](LICENSE) © 2026 Ardi Dzariat
