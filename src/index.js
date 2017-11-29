// @flow

import type {SecureStorage} from "./specification/secureStorageInterface";
import ethUtils, {EthUtilsInterface} from "./ethereum/utils";
import web3 from './ethereum/web3';
import type {JsonRpcNodeInterface} from "./specification/jsonRpcNode";

const EventEmitter = require('eventemitter3');

/**
 *
 * @param ee
 */
const preBoot = (ee:EventEmitter) => {
    "use strict";

    // Register event handler

};

/**
 *
 * @param ethNode
 * @param secureStorage
 * @param ee
 * @returns {{on: (function(string, *)), emit: (function(string)), boot: (function())}}
 * @constructor
 */
const PanthalassaApi = (ethNode:JsonRpcNodeInterface, secureStorage:SecureStorage, ee:EventEmitter) => {
    "use strict";

    const ethUtilsInstance:EthUtilsInterface = ethUtils(secureStorage, ee);

    return {

        on : (event:string, listener:(...any) => void) => {

            ee.on(event, listener);

        },

        emit: (event:string) : void => {

            ee.emit(event)

        },

        boot : () : Promise<{...any}> => {

            return new Promise((res, rej) => {

                preBoot(ee);

                res({
                    eth: ethUtilsInstance,
                    web3: web3(ethNode, ee, ethUtilsInstance)(),
                    bootNetwork : () => {
                        throw new Error("This is currently not implemented");
                    }
                })

            })

        }

    }

};

/**
 *
 * @param ethNode
 * @param ss
 * @returns {{on: (function(string, *)), emit: (function(string)), boot: (function())}}
 */
export default function(ethNode: JsonRpcNodeInterface, ss:SecureStorage) : PanthalassaApi {

    return PanthalassaApi(
        ethNode,
        ss,
        new EventEmitter()
    )

}
