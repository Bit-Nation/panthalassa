#!/usr/bin/env bash
# This script is intended to run from Bitrise CI - Assumes Mac OS in order to install 'jq' command

# fail if any commands fails
set -e
# debug log
set -x

# Set version in package.json
echo "Setting package.json version: ${BITRISE_GIT_TAG}" 
brew install jq
cat package.json | jq ".version = \"${BITRISE_GIT_TAG}\"" > tmp.json
cat tmp.json > package.json
rm tmp.json
