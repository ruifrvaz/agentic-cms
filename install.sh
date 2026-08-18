#!/bin/bash
# agentic-cms installer script
# Usage: curl -fsSL https://raw.githubusercontent.com/ruifrvaz/agentic-cms/main/install.sh | bash
#
# Options (set as environment variables):
#   AGENTIC_CMS_VERSION=latest      Install latest release (default)
#   AGENTIC_CMS_VERSION=v0.3.0      Install a specific version
#
# Examples:
#   curl -fsSL https://raw.githubusercontent.com/ruifrvaz/agentic-cms/main/install.sh | bash
#   curl -fsSL https://raw.githubusercontent.com/ruifrvaz/agentic-cms/main/install.sh | AGENTIC_CMS_VERSION=v0.3.0 bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
REPO="ruifrvaz/agentic-cms"
INSTALL_DIR="${HOME}/.local/bin"
AGENTIC_CMS_VERSION="${AGENTIC_CMS_VERSION:-latest}"

# Helper functions
info() {
    echo -e "${GREEN}✓${NC} $1"
}

warn() {
    echo -e "${YELLOW}⚠${NC} $1"
}

error() {
    echo -e "${RED}✗${NC} $1"
    exit 1
}

# Detect OS and architecture
detect_platform() {
    local os arch
    os=$(uname -s | tr '[:upper:]' '[:lower:]')
    arch=$(uname -m)

    case "$os" in
        linux)
            OS="linux"
            ;;
        darwin)
            OS="darwin"
            ;;
        mingw*|msys*|cygwin*)
            OS="windows"
            ;;
        *)
            error "Unsupported operating system: $os"
            ;;
    esac

    case "$arch" in
        x86_64|amd64)
            ARCH="amd64"
            ;;
        arm64|aarch64)
            ARCH="arm64"
            ;;
        *)
            error "Unsupported architecture: $arch"
            ;;
    esac

    info "Detected platform: ${OS}/${ARCH}"
}

# Get the release version to install from GitHub
get_version() {
    info "Fetching release version..."

    case "$AGENTIC_CMS_VERSION" in
        latest)
            VERSION=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/' || echo "")
            ;;
        v*.*.*)
            VERSION="$AGENTIC_CMS_VERSION"
            ;;
        *)
            error "Invalid AGENTIC_CMS_VERSION: $AGENTIC_CMS_VERSION (use 'latest' or 'vX.Y.Z')"
            ;;
    esac

    if [ -z "$VERSION" ]; then
        error "Failed to fetch release version"
    fi

    info "Installing version: ${VERSION}"
}

# Download binary
download_binary() {
    local binary_name="agentic-cms_${OS}_${ARCH}"

    if [ "$OS" = "windows" ]; then
        binary_name="${binary_name}.exe"
    fi

    local download_url="https://github.com/${REPO}/releases/download/${VERSION}/${binary_name}"
    TEMP_FILE="/tmp/agentic-cms_${VERSION}"

    info "Downloading from ${download_url}..."

    if ! curl -fsSL -o "$TEMP_FILE" "$download_url"; then
        error "Failed to download binary. If no asset is published yet for ${OS}/${ARCH}, install from source instead (see README.md)."
    fi

    info "Download complete"
}

# Install binary
install_binary() {
    local target="${INSTALL_DIR}/agentic-cms"

    mkdir -p "$INSTALL_DIR"
    chmod +x "$TEMP_FILE"
    mv "$TEMP_FILE" "$target"

    info "Installed to ${target}"
}

# Verify installation
verify_installation() {
    local target="${INSTALL_DIR}/agentic-cms"

    if ! "$target" --version &>/dev/null; then
        error "Installation verification failed"
    fi

    local installed_version
    installed_version=$("$target" --version 2>&1 || echo "unknown")
    info "Verified installation: ${installed_version}"
}

# Check if install directory is in PATH
check_path() {
    if [[ ":$PATH:" != *":${INSTALL_DIR}:"* ]]; then
        warn "${INSTALL_DIR} is not in your PATH"
        echo ""
        echo "Add to your shell config (~/.bashrc, ~/.zshrc, etc.):"
        echo "  export PATH=\"${INSTALL_DIR}:\$PATH\""
        echo ""
        echo "Then reload your shell:"
        echo "  source ~/.bashrc  # or ~/.zshrc"
        echo ""
    fi
}

# Main installation flow
main() {
    echo "agentic-cms installer"
    echo "======================"
    echo ""

    detect_platform
    get_version

    download_binary
    install_binary
    verify_installation
    check_path

    echo ""
    info "Installation complete!"
    echo ""
    echo "Get started:"
    echo "  cd your-project"
    echo "  agentic-cms init      # install the scaffolding"
    echo "  agentic-cms --help    # view available commands"
    echo ""
}

main
