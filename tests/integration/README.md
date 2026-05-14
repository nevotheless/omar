# Integration Tests

This directory contains end-to-end tests for the `omar` CLI.

## Structure

- `integration_test.go` — Go integration tests (builds the binary, invokes it as a subprocess)
- `README.md` — This file

## Running

```bash
# Quick: Go integration tests only
make test-integration-go

# Full: Shell + Go integration tests
make test-integration

# Direct
go test -tags=integration -count=1 -v ./tests/integration/
```

## What's tested

| Command | Test | What it verifies |
|---------|------|-----------------|
| `omar version` | Text + JSON output | Version string, JSON validity, stable fields |
| `omar status` | Output format | Key labels present (Image, Version, Booted, Staged, Rollback) |
| `omar --help` | All commands listed | All 6 subcommands appear in help |
| `omar nonexistent` | Error handling | Unknown command produces error |
| `omar install --fresh` | Missing disk | Error when `--disk` not provided |
| `omar pkg add` | No args | Error when package name missing |
| `omar update` | Non-ostree | Graceful error on mutable system |
| `omar rollback` | Non-ostree | Graceful error on mutable system |

## Test architecture

- **No mocks** — the real `omar` binary is built and invoked
- **System-independent** — tests pass on both mutable and immutable systems
- **Graceful degradation** — commands that require bootc/flatpak/distrobox are tested for sensible error messages rather than crashes
