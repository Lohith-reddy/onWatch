#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
APP_PATH="$ROOT_DIR/OnWatch Launcher.app"
LAUNCHER="$ROOT_DIR/launch_onwatch.sh"

if [[ ! -f "$LAUNCHER" ]]; then
  echo "Launcher script not found: $LAUNCHER"
  exit 1
fi

if ! command -v osacompile >/dev/null 2>&1; then
  echo "osacompile not found. This script only works on macOS."
  exit 1
fi

TMP_APPLESCRIPT="$(mktemp)"
cat > "$TMP_APPLESCRIPT" <<EOF
on run
  tell application "Terminal"
    activate
    do script "bash " & quoted form of POSIX path of "$LAUNCHER"
  end tell
end run
EOF

rm -rf "$APP_PATH"
osacompile -o "$APP_PATH" "$TMP_APPLESCRIPT"
rm -f "$TMP_APPLESCRIPT"

echo "Created app launcher: $APP_PATH"
echo "You can now double-click 'OnWatch Launcher.app' to start onWatch."

