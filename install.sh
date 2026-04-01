#!/bin/sh
set -e

# gh-quoi installer script
# Usage: curl -sSL https://raw.githubusercontent.com/BenbenIO/gh-quoi/main/install.sh | sh

REPO="BenbenIO/gh-quoi"
BINARY_NAME="gh-quoi"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

# Detect OS and architecture
OS="$(uname -s)"
ARCH="$(uname -m)"

case "$OS" in
Linux*)
    OS="Linux"
    ;;
Darwin*)
    OS="Darwin"
    ;;
*)
    echo "Unsupported operating system: $OS"
    exit 1
    ;;
esac

case "$ARCH" in
x86_64)
    ARCH="x86_64"
    ;;
amd64)
    ARCH="x86_64"
    ;;
arm64 | aarch64)
    ARCH="arm64"
    ;;
armv7l | armv7*)
    ARCH="armv7"
    ;;
*)
    echo "Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

echo "Detected OS: $OS"
echo "Detected Architecture: $ARCH"
echo "Install directory: $INSTALL_DIR"

# Get latest release version
echo "Fetching latest release..."
LATEST_VERSION=$(curl -sSL "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_VERSION" ]; then
    echo "Failed to fetch latest version"
    exit 1
fi

echo "Latest version: $LATEST_VERSION"

# Construct download URL
DOWNLOAD_URL="https://github.com/$REPO/releases/download/$LATEST_VERSION/${BINARY_NAME}_${LATEST_VERSION#v}_${OS}_${ARCH}.tar.gz"

echo "Downloading from: $DOWNLOAD_URL"

# Create temporary directory
TMP_DIR=$(mktemp -d)
trap 'rm -rf "$TMP_DIR"' EXIT

# Download and extract
cd "$TMP_DIR"
if command -v curl >/dev/null 2>&1; then
    curl -sSL "$DOWNLOAD_URL" -o "${BINARY_NAME}.tar.gz"
elif command -v wget >/dev/null 2>&1; then
    wget -q "$DOWNLOAD_URL" -O "${BINARY_NAME}.tar.gz"
else
    echo "Error: curl or wget is required"
    exit 1
fi

tar -xzf "${BINARY_NAME}.tar.gz"

# Check if we need sudo
if [ -w "$INSTALL_DIR" ]; then
    SUDO=""
else
    SUDO="sudo"
    echo "Note: Sudo access is required to install to $INSTALL_DIR"
fi

# Install binary
echo "Installing $BINARY_NAME to $INSTALL_DIR..."
$SUDO mv "$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
$SUDO chmod +x "$INSTALL_DIR/$BINARY_NAME"

echo ""
echo "✓ $BINARY_NAME $LATEST_VERSION installed successfully!"
echo ""
echo "Run '$BINARY_NAME --help' to get started"
echo ""

# Verify installation
if command -v "$BINARY_NAME" >/dev/null 2>&1; then
    "$BINARY_NAME" version
else
    echo "Warning: $INSTALL_DIR may not be in your PATH"
    echo "Add it to your PATH or run: $INSTALL_DIR/$BINARY_NAME"
fi
