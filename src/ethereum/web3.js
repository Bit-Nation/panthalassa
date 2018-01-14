// @flow

import {JsonRpcNodeInterface} from '../specification/jsonRpcNode';
import type {EthUtilsInterface} from './utils';
import PanthalassaProvider from './PanthalassaProvider';
const EventEmitter = require('eventemitter3');
const Web3 = require('web3');

/**
 * @name web3Factory
 * @param {JsonRpcNodeInterface} node
 * @param {EventEmitter} ee
 * @param {EthUtilsInterface} ethUtils
 * @param {boolean} networkAccess
 * @return {Promise} resolves with an web3 object when the node is started successfully. If the node fail to start the promise will be rejected.
 */
export default function web3Factory(node: JsonRpcNodeInterface, ee: EventEmitter, ethUtils: EthUtilsInterface, networkAccess: boolean): Promise<Web3> {
    return new Promise((res, rej) => {
        ethUtils
            .allKeyPairs()
            .then((keyPairsMap) => {
                // the keys are valid ethereum addresses
                const addresses = Object.keys(keyPairsMap);

                if (networkAccess === true) {
                    const web3 = new Web3();

                    if (addresses.length !== 0) {
                        web3.eth.defaultAccount = addresses[0];
                    }

                    return res(web3);
                }

                // Start the ethereum node
                node
                    .start()
                    .then((_) => {
                        const provider = new PanthalassaProvider(ethUtils, node.url);

                        const web3 = new Web3(provider);

                        if (addresses.length !== 0) {
                            web3.eth.defaultAccount = addresses[0];
                        }

                        return res(web3);
                    })
                    .catch(rej);
            })
            .catch(rej);
    });
}
