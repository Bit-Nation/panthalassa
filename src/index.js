// @flow

const ee = require('eventemitter3');
const ethUtils = require('./ethereum/utils');

/**
 *
 * @param ee
 */
const preBoot = (ee) => {
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
const Panthalassa = (secureStorage:any, ee:any) => {
    "use strict";

    //Ethereum utils
    const eth = ethUtils(secureStorage, ee);

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
                    eth : {
                        createPrivateKey: eth.createPrivateKey,
                        savePrivateKey: eth.savePrivateKey,
                        allKeyPairs: eth.allKeyPairs,
                        getPrivateKey: eth.getPrivateKey,
                        deletePrivateKey: eth.deletePrivateKey,
                        decryptPrivateKey: eth.decryptPrivateKey
                    },
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
 * @param secureStorage
 * @returns {{on, boot}}
 */
const factory = (secureStorage:any) => {
    "use strict";
    return Panthalassa(
        secureStorage,
        new ee()
    )
};

module.exports = {
    Panthalassa,
    factory
};
