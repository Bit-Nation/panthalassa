// @flow

import web3Factory from './ethereum/web3';
import utilsFactory from './ethereum/utils';
import profileFactory from './profile/profile';
import dbFactory from './database/db';
import {SecureStorageInterface} from './specification/secureStorageInterface';
import type {JsonRpcNodeInterface} from './specification/jsonRpcNode';
import walletFactory from './ethereum/wallet';
import type {OsDependenciesInterface} from './specification/osDependencies';
import {APP_OFFLINE, AMOUNT_OF_ADDRESSES_CHANGED, APP_ONLINE} from './events';
const EventEmitter = require('eventemitter3');

/**
 *
 * @param {SecureStorageInterface} ss
 * @param {string} dbPath
 * @param {JsonRpcNodeInterface} rpcNode
 * @param {OsDependenciesInterface} osDeps
 * @param {EventEmitter} ee
 * @param {boolean} networkAccess
 * @return {Promise<any>}
 */
export default function panthalassaFactory(ss: SecureStorageInterface, dbPath: string, rpcNode: JsonRpcNodeInterface, osDeps: OsDependenciesInterface, ee: EventEmitter, networkAccess: boolean): Promise<{...mixed}> {
    const db = dbFactory(dbPath);
    const ethUtils = utilsFactory(ss, ee, osDeps);
    const profile = profileFactory(db, ethUtils);

    const panthalassa = {
        eth: {
            utils: ethUtils,
        },

        profile: {
            profile,
        },
    };

    // /////////////////////////////////////////////////////
    // update web3 (and wallet since it depends on web3) //
    // /////////////////////////////////////////////////////

    // listen for network change
    ee.on(APP_OFFLINE, () => {
        networkAccess = false;

        web3Factory(rpcNode, ethUtils, false)
            .then((web3) => {
                panthalassa.eth.wallet = walletFactory(ethUtils, web3, db);
            })
            .catch((e) => {
                throw e;
            });
    });

    // listen for network change
    ee.on(APP_ONLINE, () => {
        networkAccess = true;

        web3Factory(rpcNode, ethUtils, true)
            .then((web3) => {
                panthalassa.eth.wallet = walletFactory(ethUtils, web3, db);
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
        web3Factory(rpcNode, ethUtils, networkAccess)
            .then((web3) => {
                panthalassa.eth.wallet = walletFactory(ethUtils, web3, db);
            })
            .catch((e) => {
                throw e;
            });
    });

    return new Promise((res, rej) => {
        web3Factory(rpcNode, ethUtils, networkAccess)
            .then((web3) => {
                panthalassa.eth.wallet = walletFactory(ethUtils, web3, db);
                res(panthalassa);
            })
            .catch(rej);
    });
}
