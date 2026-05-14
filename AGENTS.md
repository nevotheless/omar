# omar Agent Instructions

This file configures the opencode agent for the omar project.

## Project

- **Repo**: github.com/nevotheless/omar
- **Description**: Immutable OCI-image-based omarchy OS
- **Version**: 2026.5.0 (CalVer)
- **Language**: Go (CLI), mkosi (OCI images)

## Task Management

Use **GitHub Issues** as the project management backend.

### Workflow

1. **Issue erstellen** für jede Aufgabe (enhancement, bug, feature)
2. **Labels** verwenden: `enhancement`, `bug`, `priority-high`, `good-first-issue`
3. **Issue schließen** wenn der entsprechende PR gemerged wurde
4. **Issues referenzieren** in Commit-Nachrichten (`Closes #N`)

### Issue-Format

```
Titel: Kurze prägnante Beschreibung
Body:
- Was: Beschreibung der Aufgabe
- Warum: Motivation / Kontext
- Definition of Done: Checkliste
- Refs: Links zu verwandten Issues/PRs
```

## Repository-Struktur

```
cmd/omar/       # CLI Entrypoint (Go + Cobra)
internal/       # Pakete (bootc, convert, pkg, image, update, version)
images/         # Build-Pipeline (mkosi)
  packages/     # Paketlisten (base, hyprland, omarchy, immutable)
  scripts/      # Postinstall- und Konfigurations-Skripte
.devcontainer/  # Entwicklungsumgebung
.github/workflows/  # CI/CD
```

## Skill

Für detaillierte Projekt-Kontexte den omar-Skill laden:
```
skill load omar
```
(Der Skill ist automatisch in `.opencode/skills/omar/` verfügbar.)

## Build

```bash
make build VERSION=2026.5.0-dev    # CLI bauen
make test                           # Tests ausführen
make image                          # OCI-Image bauen (braucht root)
```

## Versioning

- CalVer: `2026.5.0` (Jahr.Monat.Patch)
- Dev: `2026.5.0-dev`
- Image Tags: `rolling-YYYYMMDD`, `rolling`, `vYYYY.M.P`

## Token

Der GitHub-Token liegt in `~/.gh_access_token_omar`.
Setzen für gh: `export GH_TOKEN=$(cat ~/.gh_access_token_omar)`
