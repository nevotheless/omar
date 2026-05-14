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

## Workflow

1. Check GitHub Issues: `gh issue list --repo nevotheless/omar`
2. Create todos: use the `todowrite` tool
3. Use `sequential-thinking` for complex multi-step planning
4. Use `pty_spawn`/`pty_write`/`pty_read` for long-running tasks in parallel
5. Implement, test, commit
6. Close issues with `Closes #N` in commit messages

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
