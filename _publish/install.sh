#!/usr/bin/env bash

# fail if any commands fails
set -e

PACKAGE_VERSION=$(node -p "require('./package.json').version")

echo "Downloading Panthalassa Binary Release ${PACKAGE_VERSION}"
curl https://github.com/Bit-Nation/panthalassa/releases/download/$(PACKAGE_VERSION)/panthalassa-binaries-$(PACKAGE_VERSION).zip -o panthalassa.zip
unzip build.zip