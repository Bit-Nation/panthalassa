#!/usr/bin/env bash
node_modules/.bin/flow-remove-types src --out-dir docs_build
node_modules/.bin/documentation build docs_build/database/db.js docs_build/ethereum/PanthalassaProvider.js docs_build/ethereum/utils.js docs_build/ethereum/wallet.js docs_build/ethereum/web3.js -f md > docs.md