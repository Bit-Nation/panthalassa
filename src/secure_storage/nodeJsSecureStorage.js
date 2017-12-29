// @flow
const dirty = require('dirty');
const fs = require('fs');
import {SecureStorage} from "../specification/secureStorageInterface";

/////////////////////////////////////////////////////////////////////////////////////
// This key value storage is not save. It's only intend to be used for development //
/////////////////////////////////////////////////////////////////////////////////////

/**
 *
 * @param path
 */
export default function(path:string) : SecureStorage {

    const db = dirty(path);

    const secureStorageImplementation : SecureStorage = {

        set: (key:string, value:any) : Promise<void> => new Promise((res, rej) => {
            db.set(key, value, () => res())
        }),

        get: (key:string) : Promise<mixed> => new Promise((res, rej) => res(db.get(key))),

        remove: (key:string) : Promise<void> => new Promise((res, rej) : void => db.rm(key, () => res())),

        has: (key:string) : Promise<boolean> => new Promise((res, rej) => {

            if('undefined' === typeof db.get(key)){
                res(false);
                return;
            }

            res(true);

        }),

        fetchItems: (filter: (key:string, value:mixed) => boolean) : Promise<{}> => new Promise((res, rej) => {

            const filterdItems:{} = {};

            db.forEach(function(key, val) {

                if(true === filter(key, val)){

                    filterdItems[key] = val;

                }

            });

            res(filterdItems);

        }),

        destroyStorage: () => new Promise((res, rej) => {

            fs.unlink(path, (err) => {

                if(err){
                    rej(err);
                    return;
                }

                res();

            })

        })

    };

    return secureStorageImplementation;
}
