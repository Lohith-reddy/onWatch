#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$ROOT_DIR"

echo "[onWatch] Starting launcher..."

if ! command -v go >/dev/null 2>&1; then
  echo "[onWatch] Go not found."
  if command -v brew >/dev/null 2>&1; then
    echo "[onWatch] Installing Go with Homebrew..."
    brew install go
  else
    echo "[onWatch] Homebrew not found. Install Go manually from https://go.dev/dl"
    exit 1
  fi
fi

if [[ ! -f ".env" ]]; then
  echo "[onWatch] No .env found. Copying from .env.example"
  cp .env.example .env
  echo "[onWatch] Update .env with your provider tokens/accounts before first real use."
fi

if command -v open >/dev/null 2>&1; then
  (sleep 2; open "http://127.0.0.1:9211" >/dev/null 2>&1 || true) &
fi

echo "[onWatch] Running app (debug mode)..."
echo "[onWatch] URL: http://127.0.0.1:9211"
go run . --debug

