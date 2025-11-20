#!/bin/bash
# Copyright (c) 2024-2025 Six After, Inc.
#
# This source code is licensed under the Apache 2.0 License found in the
# LICENSE file in the root directory of this source tree.
set -e

__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${__dir}"/os-type.sh

# Windows
if is_windows; then
    echo "[ERROR] Windows is not currently supported." >&2
    exit 1
fi

# ------------------------------------------------------------
# Detect latest release (your exact README method)
# ------------------------------------------------------------
TAG=$(curl -s https://api.github.com/repos/sixafter/aes-ctr-drbg/releases/latest | jq -r .tag_name)
VERSION=${TAG#v}

MODULE="github.com/sixafter/aes-ctr-drbg"

echo "Latest release: $TAG (version: $VERSION)"

# ------------------------------------------------------------
# Portable SHA-256 function (works on macOS + Linux)
# ------------------------------------------------------------
if command -v sha256sum >/dev/null 2>&1; then
  SHA256="sha256sum"
else
  SHA256="shasum -a 256"
fi

# ------------------------------------------------------------
# 1. GitHub Tag ZIP (this is what proxies ingest)
# ------------------------------------------------------------
echo "Downloading GitHub tag archive..."
curl -sSfL -o github.zip \
  "https://github.com/sixafter/aes-ctr-drbg/archive/refs/tags/${TAG}.zip"

GITHUB_SHA=$($SHA256 github.zip | awk '{print $1}')
echo "GitHub ZIP SHA256:  $GITHUB_SHA"

# ------------------------------------------------------------
# 2. Direct go mod ZIP (local module cache)
# ------------------------------------------------------------
echo "Downloading go mod ZIP using direct mode..."

MOD_JSON=$(GOPROXY=direct go mod download -json "${MODULE}@${TAG}")
MOD_ZIP_PATH=$(echo "$MOD_JSON" | jq -r '.Zip')

if [ ! -f "$MOD_ZIP_PATH" ]; then
  echo "ERROR: The go mod ZIP path does not exist:"
  echo "$MOD_ZIP_PATH"
  exit 1
fi

cp "$MOD_ZIP_PATH" gomod.zip
GOMOD_SHA=$($SHA256 gomod.zip | awk '{print $1}')
echo "go mod ZIP SHA256:  $GOMOD_SHA"

# ------------------------------------------------------------
# 3. Go Proxy ZIP
# ------------------------------------------------------------
echo "Downloading Go module proxy ZIP..."
curl -sSfL -o proxy.zip \
  "https://proxy.golang.org/${MODULE}/@v/${TAG}.zip"

PROXY_SHA=$($SHA256 proxy.zip | awk '{print $1}')
echo "Proxy ZIP SHA256:   $PROXY_SHA"

# ------------------------------------------------------------
# Comparison
# ------------------------------------------------------------
echo
echo "Comparing checksums..."
echo "GitHub : $GITHUB_SHA"
echo "go mod : $GOMOD_SHA"
echo "Proxy  : $PROXY_SHA"
echo

if [ "$GITHUB_SHA" != "$GOMOD_SHA" ] || [ "$GITHUB_SHA" != "$PROXY_SHA" ]; then
  echo "ERROR: CHECKSUM MISMATCH DETECTED!"
  exit 1
fi

echo "Go module archive is fully reproducible across GitHub, direct, and proxy."
