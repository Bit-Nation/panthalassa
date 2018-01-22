// @flow

import web3Factory from './ethereum/web3';
import utilsFactory from './ethereum/utils';
import profileFactory from './profile/profile';
import dbFactory from './database/db';
import {SecureStorageInterface} from './specification/secureStorageInterface';
import type {JsonRpcNodeInterface} from './specification/jsonRpcNode';
import walletFactory from './ethereum/wallet';
import type {OsDependenciesInterface} from './specification/osDependencies';
import nationFactory from './ethereum/nation';
import {APP_OFFLINE, AMOUNT_OF_ADDRESSES_CHANGED, APP_ONLINE} from './events';
import txQueueFactory from './queues/transaction';
import messagingFactory from './queues/messaging';
import {NATION_CONTRACT_ABI, NATION_CONTRACT_ADDRESS_DEV, NATION_CONTRACT_ADDRESS_PROD} from './constants';
const EventEmitter = require('eventemitter3');

/**
 * @param {SecureStorageInterface} ss
 * @param {string} dbPath
 * @param {JsonRpcNodeInterface} rpcNode
 * @param {OsDependenciesInterface} osDeps
 * @param {EventEmitter} ee
 * @param {boolean} networkAccess
 * @param {boolean} production
 * @return {Promise<{...mixed}>}
 */
export default function pangeaLibsFactory(ss: SecureStorageInterface, dbPath: string, rpcNode: JsonRpcNodeInterface, osDeps: OsDependenciesInterface, ee: EventEmitter, networkAccess: boolean, production: boolean): Promise<{...mixed}> {
    const db = dbFactory(dbPath);
    const ethUtils = utilsFactory(ss, ee, osDeps);
    const profile = profileFactory(db, ethUtils);
    const msgQueue = messagingFactory(ee, db);
    const txQueue = txQueueFactory(db, ee);

    const nationContractAddress = (production ? NATION_CONTRACT_ADDRESS_PROD : NATION_CONTRACT_ADDRESS_DEV);

    const pangeaLibs = {
        eventEmitter: ee,
        eth: {
            utils: ethUtils,
        },
        queue: {
            txQueue: txQueue,
        },
        profile: {
            profile,
        },
    };

    // /////////////////////////////////////////////////////
    // update web3 (and wallet since it depends on web3) //
    // /////////////////////////////////////////////////////

    function refreshWeb3(networkAccess) {
        web3Factory(rpcNode, ethUtils, networkAccess)
            .then((web3) => {
                // $FlowFixMe
                pangeaLibs.eth.wallet = walletFactory(ethUtils, web3, db);

                const nationContract = web3.eth.contract(NATION_CONTRACT_ABI).at(nationContractAddress);

                // $FlowFixMe
                pangeaLibs.eth.nation = nationFactory(db, txQueue, web3, ee, nationContract);
            })
            .catch((e) => {
                throw e;
            });
    }

    // listen for network change
    ee.on(APP_OFFLINE, () => {
        networkAccess = false;
        refreshWeb3(networkAccess);
    });

    // listen for network change
    ee.on(APP_ONLINE, () => {
        networkAccess = true;
        refreshWeb3(networkAccess);
    });

    /**
     * @desc When the amount of addresses changes we need to create an new instance of web3 (and wallet since it consumes web3)
     * to override the default address (in case of deleting private key, etc)
     */
    ee.on(AMOUNT_OF_ADDRESSES_CHANGED, () => {
        refreshWeb3(networkAccess);
    });

    return new Promise((res, rej) => {
        web3Factory(rpcNode, ethUtils, networkAccess)
            .then((web3) => {
                // $FlowFixMe
                pangeaLibs.eth.wallet = walletFactory(ethUtils, web3, db);

                const nationContract = web3.eth.contract(NATION_CONTRACT_ABI).at(nationContractAddress);

                // $FlowFixMe
                pangeaLibs.eth.nation = nationFactory(db, txQueue, web3, ee, nationContract);

                res(pangeaLibs);
            })
            .catch(rej);
    });
}
