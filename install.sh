#!/bin/bash
set -e

REPO="naqerl/todo-hours"
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"
BINARY_NAME="todo-hours"

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64|amd64)
        ARCH="amd64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    *)
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

case "$OS" in
    linux|darwin)
        ;;
    *)
        echo "Unsupported OS: $OS"
        exit 1
        ;;
esac

# Get latest release
LATEST_RELEASE=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_RELEASE" ]; then
    echo "Failed to get latest release"
    exit 1
fi

echo "Installing $BINARY_NAME $LATEST_RELEASE for $OS/$ARCH..."

# Create install directory
mkdir -p "$INSTALL_DIR"

# Download binary
DOWNLOAD_URL="https://github.com/$REPO/releases/download/$LATEST_RELEASE/${BINARY_NAME}_${LATEST_RELEASE}_${OS}_${ARCH}.tar.gz"
TEMP_DIR=$(mktemp -d)

curl -fsSL "$DOWNLOAD_URL" -o "$TEMP_DIR/$BINARY_NAME.tar.gz"
tar -xzf "$TEMP_DIR/$BINARY_NAME.tar.gz" -C "$TEMP_DIR"

# Install binary
cp "$TEMP_DIR/$BINARY_NAME" "$INSTALL_DIR/"
chmod +x "$INSTALL_DIR/$BINARY_NAME"

# Cleanup
rm -rf "$TEMP_DIR"

echo "Successfully installed $BINARY_NAME to $INSTALL_DIR"

# Check if install directory is in PATH
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo ""
    echo "WARNING: $INSTALL_DIR is not in your PATH"
    echo "Add the following to your shell profile:"
    echo "  export PATH=\"\$PATH:$INSTALL_DIR\""
fi
