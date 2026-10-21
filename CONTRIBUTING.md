# Contributing

Thanks for taking the time to improve `pdfgen`.

## Development setup

Requirements: Go 1.26.1 or newer, plus `golangci-lint` v2 for linting.

```bash
git clone https://github.com/ardizariat/generate-pdf.git
cd generate-pdf
make build
./pdfgen --help
```

## Before opening a pull request

Run the same checks as CI:

```bash
make fmt
make test
make race
make lint
```

`make lint` needs `golangci-lint` v2. If it is not installed, `go vet ./...` and
`gofmt -l .` are the minimum.

## Guidelines

- Keep changes small and focused; one purpose per pull request.
- Follow the package boundaries already in the repo: `cli`, `config`, `prompt`,
  `logging`, `data`, and `report`.
- Wrap errors with context (`fmt.Errorf("...: %w", err)`) and return them; do not
  call `os.Exit` or `panic` outside `cmd/pdfgen/main.go`.
- Avoid global mutable state and hidden configuration. Pass dependencies explicitly.
- New behavior needs a test. Bug fixes need a test that fails before the fix.
- Do not add a dependency for something the standard library or an existing
  dependency already does. Open an issue first if a new dependency is unavoidable.
- Keep user-facing CLI messages in Indonesian; documentation is English with an
  Indonesian translation.
- Commit messages follow Conventional Commits (`feat:`, `fix:`, `docs:`, `test:`,
  `refactor:`, `chore:`).

## Reporting bugs

Open an issue with the bug report template. Include the command you ran, the
expected result, the actual result, and your OS/Go version. Generated PDFs are
fake data, so attaching one is fine.
