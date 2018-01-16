// @flow

import {JsonRpcNodeInterface} from '../specification/jsonRpcNode';
import type {EthUtilsInterface} from './utils';
import PangeaProvider from './PangeaProvider';
const Web3 = require('web3');

/**
 * @name web3Factory
 * @param {JsonRpcNodeInterface} node
 * @param {EthUtilsInterface} ethUtils
 * @param {boolean} networkAccess
 * @return {Promise} resolves with an web3 object when the node is started successfully. If the node fail to start the promise will be rejected.
 */
export default function web3Factory(node: JsonRpcNodeInterface, ethUtils: EthUtilsInterface, networkAccess: boolean): Promise<Web3> {
    return new Promise((res, rej) => {
        ethUtils
            .allKeyPairs()
            .then((keyPairsMap) => {
                if (networkAccess === false) {
                    const web3 = new Web3();

                    // Even when creating an new instance of Web3
                    // the default account stay's so we need to reset it.
                    // In case it was previously set.
                    web3.eth.defaultAccount = undefined;

                    if (keyPairsMap.size !== 0) {
                        web3.eth.defaultAccount = keyPairsMap.keys().next().value;
                    }

                    return res(web3);
                }

                // Start the ethereum node
                node
                    .start()
                    .then((_) => {
                        const provider = new PangeaProvider(ethUtils, node.url);

                        const web3 = new Web3(provider);

                        // Even when creating an new instance of Web3
                        // the default account stay's so we need to reset it.
                        // In case it was previously set.
                        web3.eth.defaultAccount = undefined;

                        if (keyPairsMap.size !== 0) {
                            web3.eth.defaultAccount = keyPairsMap.keys().next().value;
                        }

                        return res(web3);
                    })
                    .catch(rej);
            })
            .catch(rej);
    });
}
