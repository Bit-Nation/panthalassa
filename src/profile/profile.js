//@flow
const querys = require('./../../lib/database/queries');

/**
 * Check if the user has a profile
 * @param db
 * @param query
 * @returns {function()}
 */
const hasProfile = (db:any, query: (realm:any) => Array<*>) : (() => Promise<boolean>) => {
    "use strict";
    return () : Promise<boolean> => {

        return new Promise((res, rej) => {

            db.query(query)
                .then(profiles => {

                    if(profiles.length >= 1){
                        res(true);
                        return;
                    }

                    res(false);

                })
                .catch(e => rej(e));

        });

    }
};

/**
 *
 * @param db
 * @returns {{hasProfile: (function()), raw: {hasProfile: (function(any, *))}}}
 */
module.exports = (db:any) => {
    "use strict";
    return {
        hasProfile: hasProfile(db, querys.findProfiles),
        raw: {
            hasProfile
        }
    }
};
