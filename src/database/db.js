//@flow
const Realm = require('realm');
const schemata = require('./schemata');

/**
 *
 * @param realm
 * @returns {function(*)}
 */
const query = (realm:any) : ((queryAction: (realm:any) => any) => Promise<any>) => {
    "use strict";

    return (queryAction: (realm) => any) : Promise<*> => {

        return new Promise((res, rej) => {

            realm
                .then(r => {
                    res(queryAction(r))
                })
                .catch(e => rej(e));

        });

    }

};

/**
 * Executes a write function
 * @param realm
 * @returns {function(*)}
 */
const write = (realm: {...any}) : ((writeAction: (realm:any) => void) => Promise<void>) => {

    "use strict";
    return (writeAction: (realm:any) => void) : Promise<void> => {

        return new Promise((res, rej) => {
            "use strict";

            realm
                .then(r => {

                    r.write(() => {
                        writeAction(r);
                        res();
                    })

                })
                .catch(e => rej(e));
        });

    }

};

/**
 * Creates a realm instance
 * @returns {ProgressPromise}
 */
const realmFactory = () : {...any} => {
    "use strict";

    return Realm
        .open({
            path: './database/panthalassa',
            schema: [schemata.ProfileSchema]
        })

};

module.exports = {
    factory: () => {
        "use strict";

        const r = realmFactory();

        return {
            write: write(r),
            query: query(r)
        }

    },
    raw: {
        write,
        query,
        realmFactory
    }
};
