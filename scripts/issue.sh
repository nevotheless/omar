#!/bin/bash
set -euo pipefail

TOKEN_FILE="$HOME/.gh_access_token_omar"
TOKEN=""

if [ -f "$TOKEN_FILE" ]; then
  TOKEN=$(cat "$TOKEN_FILE")
fi

if [ -z "$TOKEN" ]; then
  echo "Error: No token found in $TOKEN_FILE" >&2
  exit 1
fi

export GH_TOKEN="$TOKEN"

cmd="${1:-}"
shift 2>/dev/null || true

case "$cmd" in
  create)
    title="$1"
    body="$2"
    gh issue create \
      --repo nevotheless/omar \
      --title "$title" \
      --body "$body"
    ;;
  list)
    gh issue list --repo nevotheless/omar "$@"
    ;;
  view)
    gh issue view --repo nevotheless/omar "$@"
    ;;
  close)
    gh issue close --repo nevotheless/omar "$@"
    ;;
  reopen)
    gh issue reopen --repo nevotheless/omar "$@"
    ;;
  label)
    gh issue edit --repo nevotheless/omar "$@"
    ;;
  *)
    echo "Usage: $0 {create|list|view|close|reopen} [args]"
    echo ""
    echo "  create <title> <body>    Create a new issue"
    echo "  list [--label L] [--state S]  List issues"
    echo "  view <number>            View issue details"
    echo "  close <number>           Close an issue"
    echo ""
    exit 1
    ;;
esac
