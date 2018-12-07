#!/usr/bin/env bash

# fail if any commands fails
set -e
set -x

PACKAGE_VERSION=$1

echo "Downloading Panthalassa Binary Release ${PACKAGE_VERSION}"
curl -L "https://github.com/Bit-Nation/panthalassa/releases/download/${PACKAGE_VERSION}/panthalassa-binaries-${PACKAGE_VERSION}.zip" -o panthalassa.zip
unzip panthalassa.zip
rm panthalassa.zip