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

    raw: {
        db,
        ProfileSchema
    }

};
