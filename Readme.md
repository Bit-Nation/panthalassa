# Panthalassa
> Bitnation's backend - contains the mesh and some utils

[![Coverage Status](https://coveralls.io/repos/github/Bit-Nation/panthalassa/badge.svg?branch=feature%2Faes)](https://coveralls.io/github/Bit-Nation/panthalassa?branch=develop)
[![Build Status](https://semaphoreci.com/api/v1/florianlenz/panthalassa/branches/feature-aes/badge.svg)](https://semaphoreci.com/florianlenz/panthalassa)


## Development

1. Clone the project into `$GOPATH/src/github.com/Bit-Nation/panthalassa`
2. Run `make` to see all available commands

### Install
1. Run `make deps` to get the needed dependencies
2. Run `make install` to install the gx dependencies
3. Run `make deps_hack` to rewrite the import paths

### Build for ios
- Follow the install section first
- Run `make ios` to build for ios

### Build for android
- Follow the install section first
- Run `make android` to build for android