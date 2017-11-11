// @flow

import type {SecureStorage} from "./specification/secureStorageInterface";

const EventEmitter = require('eventemitter3');
import ethUtils, {EthUtilsInterface} from "./ethereum/utils";

/**
 *
 * @param ee
 */
const preBoot = (ee:EventEmitter) => {
    "use strict";

    // Register event handler

};

/**
 * Panthalassa
 * @param secureStorage
 * @param ee
 * @returns {{on: (function(string, *)), emit: (function(string)), boot: (function())}}
 * @constructor
 */
const PanthalassaApi = (secureStorage:SecureStorage, ee:EventEmitter) => {
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
                    bootNetwork : () => {
                        throw new Error("This is currently not implemented");
                    }
                })

            })

        }

    }

};

export default function(ss:SecureStorage){
    return PanthalassaApi(
        ss,
        new EventEmitter()
    )
}
