#!/bin/sh

# Self-reexec logic for curl | sh
if [ -z "$0" ] || [ "$0" = "sh" ]; then
  TMPFILE=$(mktemp)
  cat > "$TMPFILE"
  chmod +x "$TMPFILE"
  exec "$TMPFILE"
fi

set -e

BINARY_NAME="amauta"
BINARY_VERSION="alpha-0.4"
DOWNLOAD_URL="https://github.com/luislve17/amauta/releases/download/$BINARY_VERSION/$BINARY_NAME-$BINARY_VERSION"

echo "Downloading Amauta:$BINARY_NAME..."
curl -L "$DOWNLOAD_URL" -o "$BINARY_NAME"
chmod +x "$BINARY_NAME"

echo "✅ Downloaded '$BINARY_NAME' to current directory."
echo
printf "Do you want to move it to /usr/local/bin? [y/N]: "
read -r answer
if [ "$answer" = "y" ] || [ "$answer" = "Y" ]; then
  sudo mv "$BINARY_NAME" /usr/local/bin/
  echo "✅ Installed $BINARY_NAME to /usr/local/bin"
else
  echo "ℹ️  $BINARY_NAME remains in the current directory"
fi

# Clean up if running from a temp file
case "$0" in
  /tmp/*|/var/tmp/*) rm -f "$0" ;;
esac

