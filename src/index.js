// @flow

import type {SecureStorage} from "./specification/secureStorageInterface";
import type {JsonRpcNodeInterface} from "./specification/jsonRpcNode";
import iocc from './iocc';
import {asFunction, asValue} from 'awilix';

const EventEmitter = require('eventemitter3');

/**
 *
 * @param ee
 */
const preBoot = (ee:EventEmitter) => {
    "use strict";

    // Register event handler

};

const PanthalassaApi = (iocc:iocc) => {
    "use strict";

    const ee = iocc.resolve('_eventEmitter');

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
                    eth: iocc.resolve('ethereum:ethUtils'),
                    web3: iocc.resolve('ethereum:web3'),
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

    iocc.register({
        '_ethNode' : asValue(ethNode),
        '_secureStorage' : asValue(ss),
        '_eventEmitter': asFunction(() => { new EventEmitter() })
    });

    return PanthalassaApi(iocc)

}
