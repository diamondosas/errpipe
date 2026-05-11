#!/bin/sh
# install.sh — Run with:
# curl -fsSL https://diamondosas.github.io/errpipe/install/install.sh | sh

set -e

TOOL_NAME="errpipe"
REPO="diamondosas/errpipe"
INSTALL_DIR="$HOME/.local/bin"

# ── 1. Detect OS/arch ─────────────────────────────────────────────────────────
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"
case "$ARCH" in
  x86_64)  ARCH="amd64" ;;
  aarch64) ARCH="arm64" ;;
  armv7l)  ARCH="arm"   ;;
esac

# Try to find a match for the OS and ARCH
ASSET_PATTERN="${TOOL_NAME}.*${OS}.*${ARCH}"

# ── 2. Fetch latest release URL ───────────────────────────────────────────────
echo "Fetching latest release..."
RELEASES_JSON=$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest")
DOWNLOAD_URL=$(echo "$RELEASES_JSON" \
  | grep "browser_download_url" \
  | grep -i "$OS" \
  | grep -i "$ARCH" \
  | head -n 1 \
  | cut -d '"' -f 4)

# Fallback if the pattern matching is too strict
if [ -z "$DOWNLOAD_URL" ]; then
    DOWNLOAD_URL=$(echo "$RELEASES_JSON" \
      | grep "browser_download_url" \
      | grep -i "$TOOL_NAME" \
      | head -n 1 \
      | cut -d '"' -f 4)
fi

if [ -z "$DOWNLOAD_URL" ]; then
  echo "Error: no asset matching OS=$OS and ARCH=$ARCH found." && exit 1
fi

# ── 3. Download binary ────────────────────────────────────────────────────────
mkdir -p "$INSTALL_DIR"
echo "Downloading $TOOL_NAME from $DOWNLOAD_URL..."
curl -fsSL "$DOWNLOAD_URL" -o "$INSTALL_DIR/$TOOL_NAME"
chmod +x "$INSTALL_DIR/$TOOL_NAME"

# ── 4. Add to PATH in shell profile ──────────────────────────────────────────
add_to_path() {
  PROFILE="$1"
  LINE="export PATH=\"\$HOME/.local/bin:\$PATH\""
  if [ -f "$PROFILE" ] && ! grep -qF "$LINE" "$PROFILE"; then
    echo "" >> "$PROFILE"
    echo "# Added by $TOOL_NAME installer" >> "$PROFILE"
    echo "$LINE" >> "$PROFILE"
    echo "  → Added PATH to $PROFILE"
  fi
}

add_to_path "$HOME/.bashrc"
add_to_path "$HOME/.zshrc"
add_to_path "$HOME/.profile"

echo ""
echo "✅ $TOOL_NAME installed to $INSTALL_DIR"
echo "   Restart your shell or run: export PATH=\"\$HOME/.local/bin:\$PATH\""
