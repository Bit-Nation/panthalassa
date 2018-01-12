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
 * @return {Promise} resolves with an web3 object when the node is started successfully. If the node fail to start the promise will be rejected.
 */
export default function web3Factory(node: JsonRpcNodeInterface, ee: EventEmitter, ethUtils: EthUtilsInterface): Promise<Web3> {
    return new Promise((res, rej) => {
        // Start the ethereum node
        node
            .start()
            .then((_) => {
                const provider = new PanthalassaProvider(ethUtils, node.url);

                provider.on('error', (error) => ee.emit('eth:node:error', {error: error}));

                res(new Web3(provider));

                ee.emit('eth:node:start:success');
            })
            .catch((error) => ee.emit('eth:node:start:failed', {
                error: error,
            }));
    });
}
