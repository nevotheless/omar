---
name: omar
description: >
  Manage the omar project – immutable OCI-image-based omarchy OS.
  Go CLI, mkosi OCI builds, bootc deployment, GitHub Issues PM.
  Token in ~/.gh_access_token_omar.
---

# omar Skill

Manage the **omar** project – an immutable OCI-image-based omarchy OS.

## Project Snapshot

- **Repo**: `github.com/nevotheless/omar`
- **Version**: `2026.5.0` (CalVer)
- **Language**: Go (CLI), mkosi (OCI images)
- **Default image**: `ghcr.io/nevotheless/omar:rolling`

## Token

```bash
export GH_TOKEN=$(cat ~/.gh_access_token_omar)
```

## Repository

```
omar/
├── cmd/omar/           # CLI (Go + Cobra)
├── internal/
│   ├── bootc/          # bootc-Wrapper
│   ├── convert/        # mutable → immutable
│   ├── pkg/            # Flatpak + Distrobox
│   ├── image/          # bootc status
│   ├── update/         # Updates + Rollback
│   └── version/        # CalVer ldflags
├── images/
│   ├── mkosi.conf      # OCI build config
│   ├── packages/       # Paketlisten
│   └── scripts/        # Postinstall
├── .devcontainer/       # Dev env
├── scripts/
│   └── issue.sh        # GitHub Issues CLI
└── .github/
    ├── AGENTS.md
    └── workflows/
```

## Status

Done: CLI skeleton, internal packages, mkosi config, CI/CD, devcontainer, issue helper, CalVer,
install --fresh (bootc install), update registry check, version --json.
Open: OCI image build (mkosi build test), image variants, docs.

## Commands

```bash
make build VERSION=2026.5.0-dev
make test
make image                            # needs root + mkosi
./scripts/issue.sh list --label enhancement
export GH_TOKEN=$(cat ~/.gh_access_token_omar)
```

## CI

| Workflow | Trigger | Output |
|----------|---------|--------|
| `build-rolling` | daily 06:00 UTC | `:rolling-YYYYMMDD` + `:rolling` |
| `release` | tag `v*` | `:latest` + `:v2026.M.P` |
| `pr-checks` | PR | vet + test + mkosi build |
