// @flow

import web3Factory from './ethereum/web3';
import utilsFactory from './ethereum/utils';
import profileFactory from './profile/profile';
import dbFactory from './database/db';
import {SecureStorageInterface} from './specification/secureStorageInterface';
import type {JsonRpcNodeInterface} from './specification/jsonRpcNode';
import type {ProfileInterface} from './profile/profile';
import type {EthUtilsInterface} from './ethereum/utils';
import type {WalletInterface} from './ethereum/wallet';
import walletFactory from './ethereum/wallet';
import type {OsDependenciesInterface} from './specification/osDependencies';
import {APP_OFFLINE, AMOUNT_OF_ADDRESSES_CHANGED, APP_ONLINE} from './events';
const EventEmitter = require('eventemitter3');
const Web3 = require('web3');

export interface PanthalassaInterface {
    eventEmitter: EventEmitter,
    profile: ProfileInterface,
    eth: {
        utils: EthUtilsInterface,
        wallet: WalletInterface | null,
        web3: Web3 | null
    }
}

/**
 * @name panthalassaFactory
 * @desc Factory to create an instance of panthalassa ready to use.
 * @param {SecureStorageInterface} ss
 * @param {string} dbPath
 * @param {JsonRpcNodeInterface} rpcNode
 * @param {OsDependenciesInterface} osDeps
 * @param {EventEmitter} ee
 * @param {boolean} networkAccess
 * @return {Promise<PanthalassaInterface>}
 */
export default function panthalassaFactory(ss: SecureStorageInterface, dbPath: string, rpcNode: JsonRpcNodeInterface, osDeps: OsDependenciesInterface, ee: EventEmitter, networkAccess: boolean): Promise<PanthalassaInterface> {
    const db = dbFactory(dbPath);
    const ethUtils = utilsFactory(ss, ee, osDeps);

    const index = {
        eventEmitter: ee,
        eth: {
            utils: ethUtils,
        },
        profile: profileFactory(db, ethUtils),
    };

    // /////////////////////////////////////////////////////
    // update web3 (and wallet since it depends on web3) //
    // /////////////////////////////////////////////////////

    // listen for network change
    ee.on(APP_OFFLINE, () => {
        networkAccess = false;

        web3Factory(rpcNode, ee, ethUtils, false)
            .then((web3) => {
                index.eth.web3 = web3;
                index.eth.wallet = walletFactory(ethUtils, web3, db);
            })
            .catch((e) => {
                throw e;
            });
    });

    // listen for network change
    ee.on(APP_ONLINE, () => {
        networkAccess = true;

        web3Factory(rpcNode, ee, ethUtils, true)
            .then((web3) => {
                index.eth.web3 = web3;
                index.eth.wallet = walletFactory(ethUtils, web3, db);
            })
            .catch((e) => {
                throw e;
            });
    });

    /**
     * @desc When the amount of addresses changes we need to create an new instance of web3 (and wallet since it consumes web3)
     * to override the default address (in case of deleting private key, etc)
     */
    ee.on(AMOUNT_OF_ADDRESSES_CHANGED, () => {
        web3Factory(rpcNode, ee, ethUtils, networkAccess)
            .then((web3) => {
                index.eth.web3 = web3;
                index.eth.wallet = walletFactory(ethUtils, web3, db);
            })
            .catch((e) => {
                throw e;
            });
    });

    return new Promise((res, rej) => {
        web3Factory(rpcNode, ee, ethUtils, networkAccess)
            .then((web3) => {
                index.eth.web3 = web3;
                index.eth.wallet = walletFactory(ethUtils, web3, db);

                const impl:PanthalassaInterface = index;

                res(impl);
            });
    });
}
