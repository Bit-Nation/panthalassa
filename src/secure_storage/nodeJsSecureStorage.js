// @flow
const aes = require('crypto-js/aes');
const dirty = require('dirty');
const fs = require('fs');

/////////////////////////////////////////////////////////////////////////////////////
// This key value storage is not save. It's only intend to be used for development //
/////////////////////////////////////////////////////////////////////////////////////

/**
 * todo change return type any to a function that returns a promise
 * @param db
 * @returns {function(string, any)}
 */
const set = (db:any): any => {
    "use strict";

    return (key:string, value:any) : Promise<void> => {
        return new Promise((res, rej) => {
            db.set(key, value, () => res())
        })
    }
};

const get = (db:any) : any => {
    "use strict";

    return (key:string) : Promise<*> => {

        return new Promise((res, rej) => {

            res(db.get(key));

        });

    };

};

/**
 * Todo change return type any
 * @param db
 * @returns {function(string)}
 */
const remove = (db) : any => {
    "use strict";

    return (key:string) : Promise<void> => {

        return new Promise((res, rej) : void => {
            db.rm(key, () => {
                res();
            })
        })
    };

};

/**
 * Todo change return type any
 * @param db
 * @param get function to fetch a value from the database
 * @returns {function(string)}
 */
const has = (db:any, get: (key:string) => any) : any => {
    "use strict";

    return (key:string) : Promise<boolean> => {

        return new Promise((res, rej) => {

            get(key)
                // The promise will only be resolved when a key exist
                // so we can just resolve the parent promise
                .then(result => {

                    if('undefined' === typeof result){
                        res(false);
                    }

                    res(true);

                })
                .catch(err => rej(err));

        });

    };

};

/**
 * Todo change return type any
 * @param fs
 * @param path
 * @returns {Promise}
 */
const destroyStorage = (fs:any, path:string) : any => {
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

};

module.exports = (path:string) => {
    "use strict";

    const db = dirty(path);

    return {
        set: set(db),
        get: get(db),
        remove: remove(db),
        has: has(db, get(db)),
        destroyStorage: destroyStorage(fs, path),
        raw: {
            set,
            get,
            remove,
            has,
            destroyStorage
        }
    }

};
