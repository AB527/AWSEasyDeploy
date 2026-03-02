#!/bin/bash

set -e

REPO="AB527/AWSEasyDeploy"
PROJECT="AWSEasyDeploy"
BINARY="easy-deploy"

OS=$(uname | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

if [ "$ARCH" = "x86_64" ]; then
  ARCH="amd64"
fi

URL="https://github.com/$REPO/releases/latest/download/${PROJECT}_${OS}_${ARCH}.tar.gz"

echo "Installing $PROJECT..."

curl -sL $URL | tar xz

chmod +x $BINARY
sudo mv $BINARY /usr/local/bin/

echo "Installed successfully!"
echo "Run: easy-deploy"