//@flow
import {findProfiles} from './../database/queries'
import {DB} from "../database/db";
import {NoProfilePresent} from "../errors";
import {SecureStorage} from "../specification/secureStorageInterface";
import type {PublicProfile} from '../specification/publicProfile.js'

export interface Profile {

    hasProfile() : Promise<boolean>;
    setProfile(pseudo:string, description:string, image:string) : Promise<void>;
    getProfile() : Promise<{...any}>;
    getPublicProfile(): Promise<PublicProfile>

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
                        id: id +1,
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
 * Fetch profile
 * @param db
 * @param query
 * @returns {function()}
 */
export function getProfile(db:DB, query: (realm:any) => Array<{...any}>) : (() => Promise<{...any}>) {

    return () : Promise<{...any}> => {

        return new Promise((res, rej) => {

            db.query(query)
                .then(profiles => {

                    if(profiles.length <= 0){
                        rej(new NoProfilePresent());
                        return;
                    }

                    res(profiles[0]);

                })
                .catch(err => rej(err));

        });

    }

}

/**
 *
 * @param ethUtils
 * @param getProfile
 * @returns {function()}
 */
export function getPublicProfile(ethUtils:{...any}, getProfile: () => Promise<{...any}>) : () => Promise<PublicProfile> {

    return () : Promise<PublicProfile> => {

        return new Promise(async function(res, rej){

            try{

                //Fetch saved profile
                const sp = await getProfile();

                //Public profile
                const pubProfile:PublicProfile = {
                    pseudo: sp.pseudo,
                    description: sp.description,
                    image: sp.image,
                    ethAddresses: [],
                    version: '1.0.0'
                };

                //Fetch all keypairs
                const keyPairs = await ethUtils.allKeyPairs();

                keyPairs.map(keyPair => {
                    pubProfile.ethAddresses.push(keyPair.key);
                });

                res(pubProfile);

            }catch (e){

                rej(e);

            }

        });

    }

}

/**
 *
 * @param db
 */
export default function (db:DB, ethutils:{...any}) : Profile {

    const profileImplementation : Profile = {

        hasProfile: hasProfile(db, findProfiles),

        setProfile: setProfile(db),

        getProfile: getProfile(db, findProfiles),

        getPublicProfile: getPublicProfile(ethutils, getProfile(db, findProfiles))

    };

    return profileImplementation;

}
