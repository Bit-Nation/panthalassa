// @flow
const aes = require('crypto-js/aes');
const dirty = require('dirty');
const fs = require('fs');
import {SecureStorage} from "../specification/secureStorageInterface";

/////////////////////////////////////////////////////////////////////////////////////
// This key value storage is not save. It's only intend to be used for development //
/////////////////////////////////////////////////////////////////////////////////////

/**
 * todo change return type any to a function that returns a promise
 * @param db
 * @returns {function(string, any)}
 */
export function set(db:any): any {
    "use strict";

    return (key:string, value:any) : Promise<void> => {
        return new Promise((res, rej) => {
            db.set(key, value, () => res())
        })
    }
}

/**
 *
 * @param db
 * @returns {function(string)}
 */
export function get(db:any) : (key:string) => Promise<*> {
    "use strict";

    return (key:string) : Promise<*> => {

        return new Promise((res, rej) => {

            res(db.get(key));

        });

    };

}

/**
 * Todo change return type any
 * @param db
 * @returns {function(string)}
 */
export function remove(db) : (key:string) => Promise<void> {
    "use strict";

    return (key:string) : Promise<void> => {

        return new Promise((res, rej) : void => {
            db.rm(key, () => {
                res();
            })
        })
    };

}

/**
 *
 * @param db
 * @returns {function(string)}
 */
export function has(db:any) : (key:string) => Promise<*> {
    "use strict";

    return (key:string) : Promise<boolean> => {

        return new Promise((res, rej) => {

            if('undefined' === typeof db.get(key)){
                res(false);
                return;
            }

            res(true);

        });

    };

}

/**
 * Todo change return type any
 * @param fs
 * @param path
 * @returns {Promise}
 */
export function destroyStorage(fs:any, path:string) : () => Promise<void> {
    "use strict";

    return () => {

        return new Promise((res, rej) => {

            fs.unlink(path, (err) => {

                if(err){
                    rej(err);
                    return;
                }

                res();

            })

        });

    }

}

/**
 * Fetch items based on filter function (key, value) are passed to the function.
 * @param db
 * @returns {function(*)}
 */
export function fetchItems(db) : ((filter: (key:string, value:string) => boolean) => Promise<Array<{key: string, value: string}>>) {
    "use strict";

    return (filter: (key:string, value:any) => boolean) : Promise<Array<{key: string, value: string}>> => {

        return new Promise((res, rej) => {

            const elements:Array<{key: string, value: string}> = [];

            db.forEach(function(key, val) {

                if(true === filter(key, val)){
                    elements.push({
                        key: key,
                        value: val
                    })
                }

            });

            res(elements);

        });

    };

};

/**
 *
 * @param path
 */
export default function(path:string) : SecureStorage {

    const db = dirty(path);

    const secureStorageImplementation : SecureStorage = {

        set: set(db),

        get: get(db),

        remove: remove(db),

        has: has(db),

        fetchItems: fetchItems(db),

        destroyStorage: destroyStorage(fs, path)

    };

    return secureStorageImplementation;
}
