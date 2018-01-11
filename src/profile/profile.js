// @flow
import queries from './../database/queries';
import {DB} from '../database/db';
import {NoProfilePresent} from '../errors';
import type {PublicProfile} from '../specification/publicProfile.js';
import {ProfileObject} from '../database/schemata';
import type {EthUtilsInterface} from '../ethereum/utils';
export const PROFILE_VERSION = '1.0.0';

/**
 * @typedef ProfileInterface
 * @property {function} hasProfile
 * @property {function(pseudo: string, description: string, image: string)} setProfile
 * @property {function} getProfile
 * @property {function} getPublicProfile
 */
export interface ProfileInterface {

    hasProfile() : Promise<boolean>;
    setProfile(pseudo: string, description: string, image: string) : Promise<void>;
    getProfile() : Promise<ProfileObject>;
    getPublicProfile(): Promise<PublicProfile>

}

/**
 *
 * @param {object} db object that implements DBInterface
 * @param {object} ethUtils
 * @return {ProfileInterface}
 */
export default function(db: DB, ethUtils: EthUtilsInterface): ProfileInterface {
    const profileImplementation : ProfileInterface = {
        hasProfile: () => new Promise((res, rej) => {
            db.query(queries.findProfiles)
                .then((profiles) => {
                    if (profiles.length >= 1) {
                        res(true);
                        return;
                    }

                    res(false);
                })
                .catch((e) => rej(e));
        }),

        setProfile: (pseudo: string, description: string, image: string): Promise<void> => new Promise((res, rej) => {
            db.write((realm: any) => {
                // Since a user can create only one profile
                // we will updated the existing one if it exist

                const profiles:Array<ProfileObject> = queries.findProfiles(realm);

                // Create profile if no exist
                if (profiles.length === 0) {
                    realm.create('Profile', {
                        id: profiles.length +1,
                        pseudo: pseudo,
                        description: description,
                        image: image,
                        version: PROFILE_VERSION,
                    });

                    res();
                    return;
                }

                // Updated existing profile
                const profile = profiles[0];

                profile.pseudo = pseudo;
                profile.description = description;
                profile.image = image;

                res();
            });
        }),

        getProfile: (): Promise<ProfileObject> => new Promise((res, rej) => {
            db
                .query(queries.findProfiles)
                // Fetch the first profile or reject if user has no profiles
                .then((profiles) => {
                    if (profiles.length <= 0) {
                        rej(new NoProfilePresent());
                        return;
                    }

                    res(profiles[0]);
                })

                .catch((err) => rej(err));
        }),

        getPublicProfile: () => new Promise(async function(res, rej) {
            try {
                // Fetch saved profile
                const sp = await profileImplementation.getProfile();

                // Public profile
                const pubProfile:PublicProfile = {
                    pseudo: sp.pseudo,
                    description: sp.description,
                    image: sp.image,
                    ethAddresses: [],
                    version: '1.0.0',
                };

                // Fetch all keypairs
                const keyPairs:{} = await ethUtils.allKeyPairs();

                Object.keys(keyPairs).map((key) => pubProfile.ethAddresses.push(key));

                res(pubProfile);
            } catch (e) {
                rej(e);
            }
        }),

    };

    return profileImplementation;
}
