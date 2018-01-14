/* eslint-disable */
import utils from "../ethereum/utils";

const execSync = require('child_process').execSync;
import database, {DB} from '../database/db';
import profile from './../profile/profile';
import {NoProfilePresent} from './../errors';
const {describe, expect, test} = global;
import type {PublicProfile} from './../specification/publicProfile';
import {ProfileType} from '../database/schemata';

const DATABASE_PATH = 'database/panthalassa';

describe('profile', () => {
    'use strict';

    describe('setProfile', () => {
        /**
         * A profile has three parameters, a pseudo name, a description and a image
         * This is a functional test
         */
        test('create profile', () => {
            // Kill the database
            execSync('npm run db:flush');

            const db:DB = database(DATABASE_PATH);

            const p = profile(db);

            const expectedProfile:ProfileType = {
                id: 1,
                pseudo: 'pseudoName',
                description: 'I am a florian',
                image: 'base64...',
                version: '1.0.0',
            };

            // This is a promise used for the expect statement
            const testPromise = new Promise((res, rej) => {
                p
                    .setProfile('pseudoName', 'I am a florian', 'base64...')
                    .then((_) => {
                        return p.getProfile();
                    })
                    .then((profile) => res(JSON.stringify(profile)))
                    .catch((err) => rej(err));
            });

            return expect(testPromise)
                .resolves
                .toBe(JSON.stringify(expectedProfile));
        });

        test('try to update profile', () => {
            // Kill the database
            execSync('npm run db:flush');

            const db = database(DATABASE_PATH);

            const p = profile(db);

            const testPromie = new Promise((res, rej) => {
                p

                    // Create profile
                    .setProfile('pseudoName', 'I am a florian', 'base64...')

                    // After saving fetch it and return the promise
                    .then((_) => {
                        return p.getProfile();
                    })

                    // The profile content should match since get profile
                    // fetched the profile
                    .then((profileContent) => {
                        // Assert profile matches
                        expect({
                            pseudo: profileContent.pseudo,
                            description: profileContent.description,
                            image: profileContent.image,
                            id: profileContent.id,
                        })
                            .toEqual({
                                pseudo: 'pseudoName',
                                description: 'I am a florian',
                                image: 'base64...',
                                id: 1,
                            });

                        // Update profile after ensure that it was written to the db
                        return p.setProfile('pseudoNameUpdated', 'I am a florian updated', 'base64...new image');
                    })

                    // Fetch the updated profile
                    .then((_) => {
                        return p.getProfile();
                    })

                    .then((profile) => {
                        res(JSON.stringify({
                            id: profile.id,
                            pseudo: profile.pseudo,
                            description: profile.description,
                            image: profile.image,
                            version: profile.version,
                        }));
                    });
            });

            const expectedProfile:ProfileObject = {
                id: 1,
                pseudo: 'pseudoNameUpdated',
                description: 'I am a florian updated',
                image: 'base64...new image',
                version: '1.0.0',
            };

            return expect(testPromie)
                .resolves
                .toEqual(JSON.stringify(expectedProfile));
        });
    });

    /**
     * Fetches own profile
     */
    describe('getProfile', () => {
        /**
         * Fetch my existing profile successfully
         */
        test('get profile that exist', () => {
            // Kill the database
            execSync('npm run db:flush');

            const db:DB = database(DATABASE_PATH);

            const p = profile(db);

            const testPromise = new Promise((res, rej) => {
                p
                    // Make sure that no profile is present
                    .hasProfile()

                    // Set Profile
                    .then((hasProfile) => {
                        expect(hasProfile).toBeFalsy();
                        return p.setProfile('pedsa', 'i am a programmer', 'base64....');
                    })

                    // fetch profile
                    .then((_) => {
                        expect(_).toBeUndefined();
                        return p.getProfile();
                    })

                    .then((profile) => {
                        res(JSON.stringify(profile));
                    })

                    .catch((err) => rej(err));
            });

            const expectedProfile:ProfileObject = {
                id: 1,
                pseudo: 'pedsa',
                description: 'i am a programmer',
                image: 'base64....',
                version: '1.0.0',
            };

            return expect(testPromise)
                .resolves
                .toEqual(JSON.stringify(expectedProfile));
        });

        /**
         * Try to fetch profile that doesn't exist
         */
        test('get profile that do not exist', () => {
            // Kill the database
            execSync('npm run db:flush');

            const db:DB = database(DATABASE_PATH);

            const p = profile(db);

            return expect(p.getProfile())
                .rejects
                .toEqual(new NoProfilePresent());
        });
    });

    /**
     * Get my public profile
     * The "getPublicProfile" is different from the "getProfile" method
     */
    describe('getPublicProfile', () => {
        /**
         * Fetching my public profile should resolve in an error
         * if no profile exist
         */
        test('try to fetch profile that does not exist', () => {

            const p = profile(
                {},
                null
            );

            p.getProfile = () => {
                return new Promise((res, rej) => {
                    rej(new NoProfilePresent());
                });
            };

            return expect(p.getPublicProfile())
                .rejects
                .toBeInstanceOf(NoProfilePresent);

        });

        test('fetch my existing public profile', () => {
            // Mock allKeyPairs method since it will called.
            // The keys will be added to the address.
            const ethUtils = {
                allKeyPairs: () => {
                    return new Promise((res, rej) => {
                        const m = new Map();
                        m.set('0x7ed1e469fcb3ee19c0366d829e291451be638e59', '');
                        m.set('0xe0b70147149b4232a3aa58c6c1cd192c9fef385d', '');
                        res(m);
                    });
                },
            };

            // Mock the function which fetches profile
            const getProfile = () => {
                return new Promise((res, rej) => {
                    res({
                        pseudo: 'peasded',
                        description: 'I am a description',
                        image: 'base64....',
                        version: '1.0.0',
                    });
                });
            };

            // The expected public profile
            const expectedPublicProfile:PublicProfile = {
                pseudo: 'peasded',
                description: 'I am a description',
                image: 'base64....',
                ethAddresses: [
                    '0x7ed1e469fcb3ee19c0366d829e291451be638e59',
                    '0xe0b70147149b4232a3aa58c6c1cd192c9fef385d',
                ],
                version: '1.0.0',
            };

            const p = profile(null, ethUtils);

            p.getProfile = getProfile;

            return expect(p.getPublicProfile())
                .resolves
                .toEqual(expectedPublicProfile);
        });
    });

    describe('hasProfile', () => {
        test('true', () => {

            const db = database(DATABASE_PATH);

            //Since hasProfile will query the database under the hood we just mock the database
            db.query = () => new Promise(function (res, rej) {
                //Resolve with an list of objects since
                res([{}]);
            });

            return expect(profile(db).hasProfile())
                .resolves
                .toBeTruthy();
        });

        test('false', () => {
            // Kill the database
            execSync('npm run db:flush');

            return expect(profile(database(DATABASE_PATH)).hasProfile())
                .resolves
                .toBeFalsy();
        });

        test('error during fetch', () => {
            class TestError extends Error {}

            const db = database(DATABASE_PATH);

            db.query = () => new Promise((res, rej) => {

                rej(new TestError());

            });

            const p = profile(db);

            return expect(p.hasProfile())
                .rejects
                .toEqual(new TestError());
        });
    });
});

