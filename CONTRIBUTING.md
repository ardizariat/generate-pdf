# Contributing

Thanks for taking the time to improve `pdfgen`.

By participating, you agree to follow the [Code of Conduct](CODE_OF_CONDUCT.md).

## Development setup

Requirements: Go 1.26.10 or newer, plus `golangci-lint` v2 for linting. Go 1.26.10+
is also what CI uses, and `go` will download that toolchain automatically when
`GOTOOLCHAIN=auto` (the default).

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
make vuln
```

`make lint` needs `golangci-lint` v2. If it is not installed, `go vet ./...` and
`gofmt -l .` are the minimum. `make vuln` downloads `govulncheck` on first use.

CI runs tests on Linux, macOS, and Windows; race detection, lint, formatting,
`go mod tidy -diff`, `go vet`, and `govulncheck` run on Linux.

## Project layout

| Path | Responsibility |
|---|---|
| `cmd/pdfgen` | entry point and exit codes |
| `internal/cli` | Cobra commands, flags, generation flow, progress |
| `internal/config` | Viper loading, defaults, manifest parsing, validation |
| `internal/prompt` | interactive form |
| `internal/logging` | Zap logger construction |
| `internal/data` | fictional data, document payloads, `terbilang` |
| `internal/report` | PDF layout, rendering, verification, atomic output |

Dependencies point downward: `cmd` → `cli` → (`config`, `prompt`, `logging`, `data`,
`report`). `data` and `report` must not import Cobra, Viper, or Huh.

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

Security problems must not be reported in a public issue; follow
[SECURITY.md](SECURITY.md).

## Maintaining this repository

These are the one-time GitHub settings and the release routine for the current
single-maintainer setup.

### Issue labels

Create these labels in **Issues → Labels** (or with the `gh` commands below). The
issue templates already reference the first three.

| Label | Suggested color | Purpose |
|---|---|---|
| `bug` | `d73a4a` | confirmed or suspected defect |
| `enhancement` | `a2eeef` | feature request |
| `needs-triage` | `fbca04` | not reviewed yet |
| `good first issue` | `7057ff` | small, well-scoped contribution |
| `help wanted` | `008672` | maintainer welcomes a contributor |
| `question` | `d876e3` | usage question |

```bash
gh label create bug --color d73a4a --description "Something is broken"
gh label create enhancement --color a2eeef --description "New feature or request"
gh label create needs-triage --color fbca04 --description "Not reviewed yet"
gh label create "good first issue" --color 7057ff --description "Small, well-scoped contribution"
gh label create "help wanted" --color 008672 --description "Contributor welcome"
gh label create question --color d876e3 --description "Usage question"
```

### Branch protection for `main`

Use **Settings → Branches → Add branch protection rule** (public repositories on
GitHub Free include this) or the equivalent repository ruleset. Recommended for a
solo maintainer:

1. Branch name pattern: `main`.
2. Enable **Require a pull request before merging**.
3. Set **Required approvals** to `0`. You can merge your own PR; CI is the gate.
4. Enable **Require status checks to pass before merging** and select the CI jobs:
   the three `test (...)` jobs, `lint`, and `vuln`. Enable **Require branches to be
   up to date** if you want the final merge commit tested against current `main`.
5. Enable **Do not allow force pushes** and **Do not allow deletions**.
6. Leave **Do not allow bypassing the above settings** unchecked unless you want
   even administrators to follow the PR path. Unchecked keeps an emergency escape
   hatch for a solo maintainer.

If GitHub does not offer branch protection for your plan/repository type, use a
repository ruleset with the same options, or keep the same rules by convention and
never push directly to `main`.

Suggested workflow: create a short-lived feature branch, open a PR, wait for CI,
then squash-merge and delete the branch. Dependabot PRs use the same path.

### Release checklist

1. `main` is green and the working tree is clean.
2. Run `make test`, `make race`, `make lint`, and `make vuln` locally.
3. Update documentation if flags, behavior, Go version, or dependencies changed.
4. Create and push a SemVer tag, for example `git tag -a v0.1.0 -m v0.1.0` and
   `git push origin v0.1.0`. The release workflow runs GoReleaser, builds Linux,
   macOS, and Windows archives for `amd64`/`arm64`, and publishes `checksums.txt`.
5. Check the generated GitHub release notes, archives, and checksums. Only the
   latest release is supported, so note breaking changes in the release notes.
