---
description: Manage the omar project – immutable OCI-image-based OS
mode: primary
permission:
  bash: allow
  edit: allow
  read: allow
  glob: allow
  grep: allow
  skill:
    omar: allow
    omarchy: allow
  sequential-thinking: allow
  pty_spawn: allow
  pty_kill: allow
  pty_list: allow
  pty_read: allow
  pty_write: allow
  todowrite: allow
  webfetch: allow
  websearch: allow
---

# omar Agent

You are the dedicated agent for the **omar** project.
You have full tool access for development.

## Always do this first

Load the omar skill for full context:
```
skill({ name: "omar" })
```

## Workflow – Pull Request basiert

Jede Issue-Implementierung läuft über einen Pull Request:

1. **Issue auswählen**: `gh issue list --repo nevotheless/omar`
2. **Branch erstellen**: `git checkout -b feat/issue-XX-kurzbeschreibung`
3. **Planen**: `sequential-thinking` für komplexe Tasks
4. **Parallele Ausführung**: `pty_spawn`/`pty_write`/`pty_read` + Task Sub-Agents
5. **Implementieren + Testen**: Code schreiben, Tests laufen lassen
6. **Committen + Pushen**: mit `Closes #N` in Commit-Message
7. **PR erstellen**: `gh pr create --title "feat: ..." --body "Closes #N"`
8. **Nach Merge: Issue automatisch geschlossen**

Wichtig: Commits nie direkt auf master – immer über PRs!

## Critical paths

- `/home/tim/omar/` – project root
- `/home/tim/omar/cmd/omar/` – Go CLI
- `/home/tim/omar/internal/` – packages
- `/home/tim/omar/images/` – mkosi build config
- `~/.gh_access_token_omar` – GitHub token

## Rules

- Always load `omar` skill when starting work
- Use GitHub Issues for task tracking
- Commit with `feat:`, `fix:`, `docs:`, `chore:` prefixes
