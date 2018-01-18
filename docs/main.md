# Pangea libs
> All code in this repository is used for our mobile and desktop client. The intention is to avoid redundant work.

## Tech stack
- [Realm](https://www.npmjs.com/package/realm) is used as the database (can be used on node js & react native)
- [Jest](https://www.npmjs.com/package/jest) for testing
- [Flow](https://flow.org/) as a static type checker for JS (it's opt in via `//@flow` at the top of your file and has almost the same syntax as typescript)
- [esdoc](https://esdoc.org/) is used for the technical documentation.

## Modules
The pangea lib's are splitted into different modules e.g. `ethereum`, `profile` and so on. See the following list:

- ethereum
    - [wallet](./ethereum/wallet.md)
    - [web3](./ethereum/web3.md)
    - [utils](./ethereum/utils.md)
    - [PangeaProvider](./ethereum/pangeaProvider.md)