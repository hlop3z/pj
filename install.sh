#!/bin/sh

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
REPO="hlop3z/pj"
BINARY_NAME="pj"
INSTALL_DIR="/usr/local/bin"

# Print colored message
print_message() {
    color=$1
    shift
    printf "${color}%s${NC}\n" "$*"
}

# Detect OS
detect_os() {
    case "$(uname -s)" in
        Linux*)     echo "linux";;
        Darwin*)    echo "darwin";;
        CYGWIN*|MINGW*|MSYS*|MINGW32*|MINGW64*)    echo "windows";;
        *)          echo "unknown";;
    esac
}

# Detect architecture
detect_arch() {
    arch=$(uname -m)
    case "$arch" in
        x86_64|amd64)   echo "amd64";;
        aarch64|arm64)  echo "arm64";;
        *)              echo "unknown";;
    esac
}

# Get latest release version
get_latest_version() {
    curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed 's/.*"tag_name": *"\([^"]*\)".*/\1/'
}

# Main installation
main() {
    print_message "$GREEN" "================================="
    print_message "$GREEN" "  PJ Installer"
    print_message "$GREEN" "================================="
    echo

    # Detect system
    OS=$(detect_os)
    ARCH=$(detect_arch)

    if [ "$OS" = "unknown" ]; then
        print_message "$RED" "Error: Unsupported operating system"
        exit 1
    fi

    if [ "$ARCH" = "unknown" ]; then
        print_message "$RED" "Error: Unsupported architecture: $(uname -m)"
        exit 1
    fi

    print_message "$YELLOW" "Detected OS: $OS"
    print_message "$YELLOW" "Detected Architecture: $ARCH"
    echo

    # Get latest version
    print_message "$YELLOW" "Fetching latest release..."
    VERSION=$(get_latest_version)

    if [ -z "$VERSION" ]; then
        print_message "$RED" "Error: Could not fetch latest version"
        exit 1
    fi

    print_message "$GREEN" "Latest version: $VERSION"
    echo

    # Determine file extension and archive name
    if [ "$OS" = "windows" ]; then
        ARCHIVE_EXT="zip"
        BINARY_EXT=".exe"
    else
        ARCHIVE_EXT="tar.gz"
        BINARY_EXT=""
    fi

    ARCHIVE_NAME="${BINARY_NAME}-${OS}-${ARCH}.${ARCHIVE_EXT}"
    DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${ARCHIVE_NAME}"

    print_message "$YELLOW" "Downloading: $ARCHIVE_NAME"
    print_message "$YELLOW" "From: $DOWNLOAD_URL"
    echo

    # Create temporary directory
    TMP_DIR=$(mktemp -d)
    cd "$TMP_DIR"

    # Download archive
    if ! curl -L -o "$ARCHIVE_NAME" "$DOWNLOAD_URL"; then
        print_message "$RED" "Error: Failed to download $ARCHIVE_NAME"
        rm -rf "$TMP_DIR"
        exit 1
    fi

    print_message "$GREEN" "Download complete!"
    echo

    # Extract archive
    print_message "$YELLOW" "Extracting archive..."
    if [ "$OS" = "windows" ]; then
        if command -v unzip > /dev/null 2>&1; then
            unzip -q "$ARCHIVE_NAME"
        else
            print_message "$RED" "Error: unzip command not found. Please install unzip."
            rm -rf "$TMP_DIR"
            exit 1
        fi
    else
        tar -xzf "$ARCHIVE_NAME"
    fi

    # Install binary
    if [ "$OS" = "windows" ]; then
        # For Windows, install to user's local bin or provide instructions
        WINDOWS_INSTALL_DIR="$HOME/bin"
        mkdir -p "$WINDOWS_INSTALL_DIR"

        if [ -f "${BINARY_NAME}${BINARY_EXT}" ]; then
            mv "${BINARY_NAME}${BINARY_EXT}" "$WINDOWS_INSTALL_DIR/"
            chmod +x "$WINDOWS_INSTALL_DIR/${BINARY_NAME}${BINARY_EXT}"
            print_message "$GREEN" "Installed to: $WINDOWS_INSTALL_DIR/${BINARY_NAME}${BINARY_EXT}"
            echo
            print_message "$YELLOW" "Please add $WINDOWS_INSTALL_DIR to your PATH if not already added."
            print_message "$YELLOW" "You can do this by adding this line to your ~/.bashrc or ~/.bash_profile:"
            print_message "$YELLOW" "  export PATH=\"\$HOME/bin:\$PATH\""
        else
            print_message "$RED" "Error: Binary not found in archive"
            rm -rf "$TMP_DIR"
            exit 1
        fi
    else
        # For Linux/macOS
        if [ -f "${BINARY_NAME}${BINARY_EXT}" ]; then
            # Check if we need sudo
            if [ -w "$INSTALL_DIR" ]; then
                mv "${BINARY_NAME}${BINARY_EXT}" "$INSTALL_DIR/"
                chmod +x "$INSTALL_DIR/${BINARY_NAME}${BINARY_EXT}"
            else
                print_message "$YELLOW" "Installing to $INSTALL_DIR (requires sudo)..."
                sudo mv "${BINARY_NAME}${BINARY_EXT}" "$INSTALL_DIR/"
                sudo chmod +x "$INSTALL_DIR/${BINARY_NAME}${BINARY_EXT}"
            fi
            print_message "$GREEN" "Installed to: $INSTALL_DIR/${BINARY_NAME}${BINARY_EXT}"
        else
            print_message "$RED" "Error: Binary not found in archive"
            rm -rf "$TMP_DIR"
            exit 1
        fi
    fi

    # Cleanup
    cd - > /dev/null
    rm -rf "$TMP_DIR"

    echo
    print_message "$GREEN" "================================="
    print_message "$GREEN" "  Installation Complete!"
    print_message "$GREEN" "================================="
    echo
    print_message "$GREEN" "Run '${BINARY_NAME}' to get started!"
}

# Run main installation
main
