import {JsonRpcNodeInterface} from "../specification/jsonRpcNode";
import type {EthUtilsInterface} from "./utils";
import PanthalassaProvider from './PanthalassaProvider';
const EventEmitter = require('eventemitter3');
const Web3 = require('web3');

/**
 *
 * @param node
 * @param ee
 * @param ethUtils
 * @returns {function()}
 */
export default function (node:JsonRpcNodeInterface, ee:EventEmitter, ethUtils:EthUtilsInterface) : () => Promise<Web3> {

    return () : Promise<Web3> => {

        return new Promise((res, rej) => {

            //Start the ethereum node
            node
                .start()
                .then(_ => {

                    const provider = new PanthalassaProvider(ethUtils, node.url);

                    provider.on('error', (error) => ee.emit('eth:node:error', {error : error}));

                    res(new Web3(provider));

                    ee.emit('eth:node:start:success')

                })
                .catch(error => ee.emit('eth:node:start:failed', {
                    error: error
                }));

        });

    }

}
