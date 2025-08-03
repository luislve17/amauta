#!/bin/bash

# Re-exec from temp file if running via pipe
if [ -p /dev/stdin ]; then
  TMPFILE=$(mktemp)
  cat - > "$TMPFILE"
  chmod +x "$TMPFILE"
  exec "$TMPFILE"
fi

set -e

BINARY_NAME="amauta"
BINARY_VERSION="alpha-0.5"
DOWNLOAD_URL="https://github.com/luislve17/amauta/releases/download/$BINARY_VERSION/$BINARY_NAME-$BINARY_VERSION"

echo "Downloading Amauta:$BINARY_NAME..."
curl -L "$DOWNLOAD_URL" -o "$BINARY_NAME"
chmod +x "$BINARY_NAME"

echo "✅ Downloaded '$BINARY_NAME' to current directory."
echo "➡️  To use it globally, move it to a directory in your \$PATH (e.g. ~/.local/bin or /usr/local/bin)"

