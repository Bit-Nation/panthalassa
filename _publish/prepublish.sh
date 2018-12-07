#!/usr/bin/env bash
# This script is intended to run from Bitrise CI - Assumes Mac OS in order to install 'jq' command

# fail if any commands fails
set -e
# debug log
set -x

NPM_PACKAGE_NAME=panthalassa

echo "Exit if the 'production' branch does not contain this tag."
git merge-base --is-ancestor $BITRISE_GIT_TAG production

echo "Exit if package already exists on npm:"
! npm show $NPM_PACKAGE_NAME@* version | grep $BITRISE_GIT_TAG

echo "Exit if package doesn't exist on the web:"
github_download_status=$(curl --head --silent https://github.com/Bit-Nation/panthalassa/releases/download/${BITRISE_GIT_TAG}/panthalassa-binaries-${BITRISE_GIT_TAG}.zip | head -n 1)
if ! echo "$github_download_status" | grep -q 200
then
  echo "Published release not found on Github. Aborting. Status code: ${github_download_status}"
  return 1
fi
