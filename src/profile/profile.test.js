/* eslint-disable */
import utils from "../ethereum/utils";

const execSync = require('child_process').execSync;
import database, {DB} from '../database/db';
import profile from './../profile/profile';
import {NoProfilePresent} from './../errors';
const {describe, expect, test} = global;
import type {PublicProfile} from './../specification/publicProfile';
import {ProfileType} from '../database/schemata';

function dbPath() {
    return './database/'+Math.random()
}

describe('profile', () => {
    'use strict';

    describe('setProfile', () => {
        /**
         * A profile has three parameters, a pseudo name, a description and a image
         * This is a functional test
         */
        test('create profile', (done) => {

            const db:DB = database(dbPath());

            const p = profile(db);

            const expectedProfile:ProfileType = {
                id: 1,
                description: 'I am a florian',
                image: 'base64...',
                version: '1.0.0',
                name: '',
                location: '',
                latitude: '',
                longitude: ''
            };

            p
                .setProfile(expectedProfile)
                .then(_ => p.getProfile())
                .then((profile:ProfileType) => {

                    expect(profile.name).toBe(expectedProfile.name);
                    expect(profile.id).toBe(expectedProfile.id);
                    expect(profile.description).toBe(expectedProfile.description);
                    expect(profile.image).toBe(expectedProfile.image);
                    expect(profile.location).toBe(expectedProfile.location);
                    expect(profile.latitude).toBe(expectedProfile.latitude);
                    expect(profile.version).toBe(expectedProfile.version);

                    done();

                })
                .catch((err) => rej(err));

        });

        test('try to update profile', (done) => {
            // Kill the database

            const db = database(dbPath());

            const p = profile(db);

            const expectedProfile:ProfileType = {
                id: 1,
                description: 'I am a florian',
                image: 'base64...',
                version: '1.0.0',
                name: 'Florian',
                location: '',
                latitude: '',
                longitude: ''
            };

            p
                // Create profile
                .setProfile(expectedProfile)

                // After saving fetch it and return the promise
                .then(_ => p.getProfile())

                // The profile content should match since get profile
                // fetched the profile
                .then((profileContent:ProfileType) => {

                    expect(profileContent.name).toBe('Florian');

                    profileContent.name = 'Jaspy';

                    // Update profile after ensure that it was written to the db
                    return p.setProfile(profileContent);
                })

                // Fetch the updated profile
                .then(_ => p.getProfile())
                .then((profile) => {
                    expect(profile.name).toBe('Jaspy');
                    done();
                });

        });
    });

    /**
     * Fetches own profile
     */
    describe('getProfile', () => {
        /**
         * Fetch my existing profile successfully
         */
        test('get profile that exist', (done) => {

            const db:DB = database(dbPath());

            const p = profile(db);

            const expectedProfile:ProfileType = {
                id: 1,
                description: 'I am a florian',
                image: 'base64...',
                version: '1.0.0',
                name: 'Florian',
                location: '',
                latitude: '',
                longitude: ''
            };

            p
                // Make sure that no profile is present
                .hasProfile()

                // Set Profile
                .then(hasProfile => {
                    expect(hasProfile).toBeFalsy();
                    return p.setProfile(expectedProfile);
                })

                // fetch profile
                .then(_ => {
                    expect(_).toBeUndefined();
                    return p.getProfile();
                })

                .then((profile:ProfileType) => {
                    expect(profile.name).toBe(expectedProfile.name);
                    expect(profile.id).toBe(expectedProfile.id);
                    expect(profile.description).toBe(expectedProfile.description);
                    expect(profile.image).toBe(expectedProfile.image);
                    expect(profile.version).toBe(expectedProfile.version);
                    expect(profile.location).toBe(expectedProfile.location);
                    expect(profile.latitude).toBe(expectedProfile.latitude);
                    expect(profile.longitude).toBe(expectedProfile.longitude);
                    done();
                })

                .catch(done.fail);

        });

        /**
         * Try to fetch profile that doesn't exist
         */
        test('get profile that do not exist', () => {
            // Kill the database

            const db:DB = database(dbPath());

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
                        name: 'peasded',
                        description: 'I am a description',
                        image: 'base64....',
                        version: '1.0.0',
                    });
                });
            };

            // The expected public profile
            const expectedPublicProfile:PublicProfile = {
                name: 'peasded',
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

            const db = database(dbPath());

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

            return expect(profile(database(dbPath())).hasProfile())
                .resolves
                .toBeFalsy();
        });

        test('error during fetch', () => {
            class TestError extends Error {}

            const db = database(dbPath());

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

