#!/bin/sh
set -eu

project_root=$(CDPATH= cd -- "$(dirname "$0")/.." && pwd)
env_file="$project_root/.env.local"

if [ -f "$env_file" ]; then
  set -a
  . "$env_file"
  set +a
fi

go_bin=${GO_BIN:-}

if [ -z "$go_bin" ]; then
  if command -v go >/dev/null 2>&1; then
    go_bin=$(command -v go)
  elif [ -x "$HOME/.local/bin/go" ]; then
    go_bin="$HOME/.local/bin/go"
  elif [ -x "/usr/local/go/bin/go" ]; then
    go_bin="/usr/local/go/bin/go"
  else
    echo "Go executable not found. Set GO_BIN or install Go." >&2
    exit 1
  fi
fi

cd "$project_root"
exec "$go_bin" run ./apps/api/cmd/server
