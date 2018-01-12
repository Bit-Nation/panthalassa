//@flow

/**
 * @file Manages the configuration settings for the widget.
 * @author Rowina Sanela
 */

import {JsonRpcNodeInterface} from "../specification/jsonRpcNode";
import type {EthUtilsInterface} from "./utils";
import PanthalassaProvider from './PanthalassaProvider';
const EventEmitter = require('eventemitter3');
const Web3 = require('web3');

/**
 * @name ethereum/web3.js
 * @param node {JsonRpcNodeInterface} an instance of an object that satisfies the JsonRpcNodeInterface
 * @param ee {EventEmitter} an instance of the event emitter
 * @param ethUtils {EthUtilsInterface} an object that satisfies the EthUtilsInterface
 * @returns {Promise} resolves with an web3 object when the node is started successfully. If the node fail to start the promise will be rejected.
 */
export default function (node:JsonRpcNodeInterface, ee:EventEmitter, ethUtils:EthUtilsInterface) : () => Promise<Web3> {

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
