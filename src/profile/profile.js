//@flow
import {findProfiles} from './../database/queries'
import {DB} from "../database/db";

export interface Profile {

    hasProfile() : Promise<boolean>;
    setProfile(pseudo:string, description:string, image:string) : Promise<void>;
}

/**
 * Set profile data
 * @param db
 * @returns {function(string, string, string)}
 */
export function setProfile(db:DB) : (pseudo:string, description:string, image:string) => Promise<void> {

    return (pseudo:string, description:string, image:string) : Promise<void> => {

        return new Promise((res, rej) => {

            db.write((realm:any) => {

                //Since a user can create only one profile
                //we will updated the existing one if it exist

                const profiles = findProfiles(realm);

                //Create profile if no exist
                if(profiles.length === 0){

                    let id = profiles.length;

                    realm.create('Profile', {
                        id: id++,
                        pseudo: pseudo,
                        description: description,
                        image: image,
                    });

                    res();
                    return;
                }

                //Updated existing profile
                const profile = profiles[0];

                profile.pseudo = pseudo;
                profile.description = description;
                profile.image = image;

                res();
            });

        });

    }

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

        hasProfile: hasProfile(db, findProfiles),

        setProfile: setProfile(db)

    };

    return profileImplementation;

}
