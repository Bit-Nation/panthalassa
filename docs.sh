#!/usr/bin/env bash
node_modules/.bin/flow-remove-types ./src --out-dir docs_build
node_modules/.bin/documentation build docs_build/ethereum/wallet.js -f html -o docs
#rm -rf docs_build