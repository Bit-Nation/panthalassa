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
    setProfile(profile: ProfileType) : Promise<void>;
    getProfile() : Promise<ProfileType>;
    getPublicProfile(): Promise<PublicProfile>

}

export type ProfileInputType = {

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

        setProfile: (profile: ProfileType): Promise<void> => new Promise((res, rej) => {
            db
                .write((realm: any) => {
                    // Since we only support one profile at the moment, we can just set this always to 1
                    profile.id = 1;

                    realm.create('Profile', profile, true);
                })
                .then(res)
                .catch(rej);
        }),

        getProfile: (): Promise<ProfileType> => new Promise((res, rej) => {
            db
                .query(queries.findProfiles)
                // Fetch the first profile or reject if user has no profiles
                .then((profiles) => {
                    if (profiles.length <= 0) {
                        return rej(new NoProfilePresent());
                    }

                    const profile:ProfileType = {
                        id: profiles[0].id,
                        name: profiles[0].name,
                        description: profiles[0].description,
                        location: profiles[0].location,
                        latitude: profiles[0].latitude,
                        longitude: profiles[0].longitude,
                        image: profiles[0].image,
                        version: profiles[0].version,
                    };

                    res(profile);
                })

                .catch((err) => rej(err));
        }),

        getPublicProfile: () => new Promise(async function(res, rej) {
            try {
                // Fetch saved profile
                const sp:ProfileType = await profileImplementation.getProfile();

                const ethAddresses:Map<string, PrivateKeyType> = await ethUtils.allKeyPairs();

                // Public profile
                const pubProfile:PublicProfile = {
                    name: sp.name,
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
