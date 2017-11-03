//@flow

const Realm = require('realm');

// Profile
const ProfileSchema  = {
    name: 'Profile',
    primaryKey: 'id',
    properties: {
        id: 'int',
        pseudo: 'string',
        description: 'string',
        image: 'string'
    },
};

const db = (realm:Realm) => {
    "use strict";

    return realm.open({
        path: './database/panthalassa',
        schema: [
            ProfileSchema
        ]
    })

};

/**
 *
 * @param realm
 * @returns {function(*)}
 */
const query = (realm:any) : ((queryAction: (realm:any) => any) => Promise<any>) => {
    "use strict";

    return (queryAction: (realm) => any) : Promise<*> => {

        return new Promise((res, rej) => {

            try{
                res(queryAction(realm));
            }catch (e){
                rej(e);
            }

        });

    }

};

module.exports = {

    db: db(Realm),

    write: (writeAction: (realm:any) => void) : Promise<*> => {

        return new Promise((res, rej) => {
            "use strict";

            module.exports
                .db
                .then(realm => {

                    realm
                        .write(() => {
                            writeAction(realm);
                            res();
                        })

                })
                .catch(err => rej(err))

        });

    },

    query: query(module.exports.db),

    raw: {
        db,
        ProfileSchema,
        query
    }

};
