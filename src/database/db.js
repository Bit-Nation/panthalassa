//@flow
const Realm = require('realm');
const schemata = require('./schemata');

/**
 * Interface for database
 */
export interface DB {
    query(queryAction: (realm:Realm) => Realm.Results) : Promise<any>;
    write(writeAction: (realm:Realm) => void) : Promise<void>;
}

/**
 *
 * @param realm
 * @returns {function(*)}
 */
export function query(realm:any) : ((queryAction: (realm:any) => any) => Promise<any>){
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

}

/**
 * Executes a writeAction
 * @param realm
 * @returns {function(*)}
 */
export function write(realm: {...any}) : ((writeAction: (realm:any) => void) => Promise<void>){

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

}

export default function () : DB {

    const realm = Realm
        .open({
            path: './database/panthalassa',
            schema: [schemata.ProfileSchema, schemata.AccountBalanceSchema]
        });

    const dbImplementation : DB = {
        query: query(realm),
        write: write(realm)
    };

    return dbImplementation;
}
