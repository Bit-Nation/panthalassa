// @flow
import queries from './../database/queries';
import {DBInterface} from '../database/db';
import {NoProfilePresent} from '../errors';
import type {PublicProfile} from '../specification/publicProfile.js';
import type {ProfileType} from '../database/schemata';
import type {EthUtilsInterface} from '../ethereum/utils';
import type {PrivateKeyType} from '../specification/privateKey';
export const PROFILE_VERSION = '1.0.0';

/**
 * @typedef ProfileInterface
 * @property {function() : Promise<boolean>} hasProfile
 * @property {function(pseudo: string, description: string, image: string) : Promise<void>} setProfile
 * @property {function() : Promise<ProfileType>} getProfile
 * @property {function() : Promise<PublicProfile>} getPublicProfile
 */
export interface ProfileInterface {

    hasProfile() : Promise<boolean>;
    setProfile(pseudo: string, description: string, image: string) : Promise<void>;
    getProfile() : Promise<ProfileType>;
    getPublicProfile(): Promise<PublicProfile>

}

/**
 * @param {DBInterface} db object that satisfy the DBInterface
 * @param {EthUtilsInterface} ethUtils object that satisfy the EthUtilsInterface
 * @return {ProfileInterface}
 */
export default function profileFactory(db: DBInterface, ethUtils: EthUtilsInterface): ProfileInterface {
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

                const profiles:Array<ProfileType> = queries.findProfiles(realm);

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

        getProfile: (): Promise<ProfileType> => new Promise((res, rej) => {
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

                const ethAddresses:Map<string, PrivateKeyType> = await ethUtils.allKeyPairs();

                // Public profile
                const pubProfile:PublicProfile = {
                    pseudo: sp.pseudo,
                    description: sp.description,
                    image: sp.image,
                    ethAddresses: Array.from(ethAddresses.keys()),
                    version: '1.0.0',
                };

                res(pubProfile);
            } catch (e) {
                rej(e);
            }
        }),

    };

    return profileImplementation;
}
