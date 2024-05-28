#!/bin/sh

set -e

# Define variables
REPO="cristianrz/SimpleGoWebServer"
LATEST_URL="https://api.github.com/repos/$REPO/releases/latest"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="simplegowebserver"

# Function to check command existence
command_exists() {
	command -v "$1" >/dev/null 2>&1
}

# Check if curl is installed
if ! command_exists curl; then
	echo "Error: curl is not installed."
	exit 1
fi

# Check if jq is installed
if ! command_exists jq; then
	echo "Error: jq is not installed."
	exit 1
fi

# Fetch the latest release tag from GitHub
echo "Fetching the latest release..."
LATEST_RELEASE=$(curl -s $LATEST_URL | jq -r ".tag_name")

if [ -z "$LATEST_RELEASE" ]; then
	echo "Error: Unable to fetch the latest release."
	exit 1
fi

echo "Latest release: $LATEST_RELEASE"

# Construct the download URL
ARCH=$(uname -m)
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
DOWNLOAD_URL="https://github.com/$REPO/releases/download/$LATEST_RELEASE/$BINARY_NAME-$OS-$ARCH"

# Download the binary
echo "Downloading $BINARY_NAME from $DOWNLOAD_URL..."
curl -L -o $BINARY_NAME $DOWNLOAD_URL

# Make the binary executable
chmod +x $BINARY_NAME

# Move the binary to the install directory
echo "Installing $BINARY_NAME to $INSTALL_DIR..."
sudo mv $BINARY_NAME $INSTALL_DIR/$BINARY_NAME

# Verify installation
if command_exists $BINARY_NAME; then
	echo "$BINARY_NAME installed successfully."
else
	echo "Error: $BINARY_NAME installation failed."
	exit 1
fi
