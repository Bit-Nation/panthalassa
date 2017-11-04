//@flow
import {findProfiles} from './../database/queries'
import {DB} from "../database/db";

export interface Profile {

    hasProfile() : Promise<boolean>

}

/**
 *
 * @param db
 * @param query
 * @returns {function()}
 */
export function hasProfile(db:DB, query: (realm:any) => Array<{...any}>) : (() => Promise<boolean>) {
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
}

/**
 *
 * @param db
 */
export default function (db:DB) : Profile {

    const profileImplementation : Profile = {

        hasProfile: hasProfile(db, findProfiles)

    };

    return profileImplementation;

}
