# panthalassa

[![standard-readme compliant](https://img.shields.io/badge/standard--readme-OK-green.svg?style=flat-square)](https://github.com/RichardLitt/standard-readme)
[![Build Status](https://semaphoreci.com/api/v1/florianlenz/panthalassa/branches/develop/badge.svg)](https://semaphoreci.com/florianlenz/panthalassa) (Develop)
[![Build Status](https://semaphoreci.com/api/v1/florianlenz/panthalassa/branches/master/badge.svg)](https://semaphoreci.com/florianlenz/panthalassa) (Master)

> Backend for Pangea

TODO: Fill out this long description.

## Table of Contents

- [Security](#security)
- [Background](#background)
- [Install](#install)
- [Usage](#usage)
- [API](#api)
- [Maintainers](#maintainers)
- [Contribute](#contribute)
- [License](#license)

## Security
If you find a bug / vulnerability please DON'T open an issue. Write to `security@bitnation.co` PLEASE use [this](security-bitnation.co.key.pub) PGP key to encrypt your report / email.

## Background
[Pangea](https://github.com/Bit-Nation/BITNATION-Pangea-mobile) is the mobile interface to our blockchain jurisdiction. While smart contract's are "onchain" (on a blockchain like a thereum) communication happen's offchain.
Since current chat system's like whatsapp and telegram are hevaly centralized we are using a p2p system to send messages between peer's so that bitnaiton doesn't become a central point of failure.
We are using [libp2p](https://github.com/libp2p) for the p2p network which is truly a great project.

## Install

First clone the project. You can run the command's from the `Usage` section.

## Usage
We are using GX as the dependency manager since libp2p (and almost all go project from [Protocol Labs](https://protocol.ai/)) use it as the dependecny manager.
However, you don't really need to care about since you just have to use the make file. The following command's are available:
- `make list` (or just `make`) will list all command's from the Makefile.
- `make deps` will fetch tool's that you need in order to work with the project.
- `make install` will install all dependencies needed in order to work with the project.
- `make deps_hack` will "hack" your dependencies. GX rewrites your import path's `github.com/libp2p/go-libp2p` e.g. becomes `gx/ipfs/QmNh1kGFFdsPu79KNSaL4NUKUPb4Eiz4KHdMtFY6664RDp/go-libp2p`. You need this in order to work with the package version's specified in the package.json.
- `make deps_hack_revert` will undo `make deps_hack`. We never want to commit the GX import path's.
- `make deps_mobile` will install some tool's needed to build panthalassa for mobile. You need to run this before you can build for ios and android.
- `make ios` will build panthalassa for ios and place it in the `build` folder.
- `make android` will build panthalassa for android and place it in the build folder.
- `make test` will format the code and run all test's.
- `make test_coverage` will test the code and open the coverage report.

## API
> TODO - add link to godoc.org

## Maintainers

[@florianlenz](https://github.com/florianlenz)

## Contribute

PRs accepted.

Small note: If editing the README, please conform to the [standard-readme](https://github.com/RichardLitt/standard-readme) specification.

## License

MIT © 2018 florian@bitnation.co
