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
BINARY_VERSION="alpha-0.4"
DOWNLOAD_URL="https://github.com/luislve17/amauta/releases/download/$BINARY_VERSION/$BINARY_NAME-$BINARY_VERSION"
INSTALL_DIR="$HOME/.local/bin"

echo "Downloading Amauta:$BINARY_NAME..."
curl -L "$DOWNLOAD_URL" -o "$BINARY_NAME"
chmod +x "$BINARY_NAME"

echo "✅ Downloaded '$BINARY_NAME' to current directory."
echo
read -rp "Do you want to move it to $INSTALL_DIR? [y/N]: " answer
if [ "$answer" = "y" ] || [ "$answer" = "Y" ]; then
  mkdir -p "$INSTALL_DIR"
  mv "$BINARY_NAME" "$INSTALL_DIR/"
  echo "✅ Installed $BINARY_NAME to $INSTALL_DIR"
  echo "➕ Make sure $INSTALL_DIR is in your PATH"
else
  echo "ℹ️  $BINARY_NAME remains in the current directory"
fi

case "$0" in
  /tmp/*|/var/tmp/*) rm -f "$0" ;;
esac

