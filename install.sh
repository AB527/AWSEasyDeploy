#!/bin/bash
set -e

REPO="AB527/AWSEasyDeploy"
PROJECT="AWSEasyDeploy"
BINARY="easy-deploy"

OS=$(uname | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

[ "$ARCH" = "x86_64" ] && ARCH="amd64"

echo "Fetching latest version..."

VERSION=$(curl -s https://api.github.com/repos/$REPO/releases/latest \
  | grep tag_name \
  | cut -d '"' -f4)

FILE="${PROJECT}_${VERSION}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/$REPO/releases/download/${VERSION}/${FILE}"

echo "Installing $PROJECT $VERSION..."
echo "Downloading: $FILE"

TMP=$(mktemp)

if ! curl -fL "$URL" -o "$TMP"; then
  echo "❌ Failed to download binary"
  echo "Expected: $FILE"
  exit 1
fi

tar -xzf "$TMP"
rm "$TMP"

chmod +x $BINARY
sudo mv $BINARY /usr/local/bin/

echo "Installed successfully!"
echo "Run: easy-deploy"