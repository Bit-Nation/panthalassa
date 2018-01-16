// @flow
const Realm = require('realm');
const schemata = require('./schemata');

/**
 * @typedef {Object} DBInterface
 * @property {function} query query the realm database.
 * @property {function} write executes an write on the database
 */
export interface DBInterface {

    /**
     * Expect an query callback that will receive an instance of realm.
     * Should return the realm query result
     */
    query(queryAction: (realm: Realm) => Realm.Results) : Promise<Realm.Results>;

    /**
     * Expect an callback that that will receive an instance of realm.
     * The callback should return nothing.
     */
    write(writeAction: (realm: Realm) => void) : Promise<void>;

}

/**
 * @module database/db.js
 * @param {string} path
 * @return {DBInterface}
 */
export default function dbFactory(path: string): DBInterface {
    const realm = Realm
        .open({
            path: path,
            schema: [
                schemata.ProfileSchema,
                schemata.AccountBalanceSchema,
                schemata.MessageJobSchema
            ],
        });

    const dbImplementation : DBInterface = {

        query: (queryAction: (realm) => any): Promise<*> => {
            return new Promise((res, rej) => {
                realm
                    .then((r) => {
                        res(queryAction(r));
                    })
                    .catch((e) => rej(e));
            });
        },

        write: (writeAction: (realm: any) => void): Promise<void> => {
            return new Promise((res, rej) => {
                'use strict';

                realm
                    .then((r) => {
                        r.write(() => {
                            writeAction(r);
                            res();
                        });
                    })
                    .catch((e) => rej(e));
            });
        },

    };

    return dbImplementation;
}
