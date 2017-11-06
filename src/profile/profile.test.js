import {} from './../errors';
const execSync = require('child_process').execSync;
import {DB, factory} from "../database/db";
import {} from './../database/queries';
import profile, {hasProfile} from './../profile/profile';
import {NoProfilePresent} from './../errors';

describe('profile', () => {
    "use strict";

    describe('setProfile', () => {

        /**
         * A profile has three parameters, a pseudo name, a description and a image
         * This is a functional test
         */
        test('create profile', () => {

            // Kill the database
            execSync('npm run db:flush');

            const db:DB = factory();

            const p = profile(db);

            //This is a promise used for the expect statement
            const testPromise = new Promise((res, rej) => {

                p
                    .setProfile('pseudoName', 'I am a florian', 'base64...')
                    .then(_ => {
                        return p.getProfile();
                    })
                    .then(profile => {
                        res({
                            pseudo: profile.pseudo,
                            description: profile.description,
                            image: profile.image,
                            id: profile.id
                        });
                    })
                    .catch(err => rej(err));

            });

            return expect(testPromise)
                .resolves
                .toEqual({
                    pseudo: 'pseudoName',
                    description: 'I am a florian',
                    image: 'base64...',
                    id: 1
                })

        });

        test('try to update profile', () => {

            // Kill the database
            execSync('npm run db:flush');

            const db = factory();

            const p = profile(db);

            const testPromie = new Promise((res, rej) => {

                p

                    //Create profile
                    .setProfile('pseudoName', 'I am a florian', 'base64...')

                    //After saving fetch it and return the promise
                    .then(_ => {

                        return p.getProfile();

                    })

                    //The profile content should match since get profile
                    //fetched the profile
                    .then(profileContent => {

                        const profileAsObj = {
                            pseudo: profileContent.pseudo,
                            description: profileContent.description,
                            image: profileContent.image,
                            id: profileContent.id
                        };

                        //Assert profile matches
                        expect(profileAsObj).toEqual({
                            pseudo: 'pseudoName',
                            description: 'I am a florian',
                            image: 'base64...',
                            id: 1
                        });

                        //Update profile after ensure that it was written to the db
                        return p.setProfile('pseudoNameUpdated', 'I am a florian', 'base64...');

                    })

                    // Fetch the updated profile
                    .then(_ => {

                        return p.getProfile()

                    })

                    .then(profile => {
                        res({
                            pseudo: profile.pseudo,
                            description: profile.description,
                            image: profile.image,
                            id: profile.id
                        })
                    })

            });

            return expect(testPromie)
                .resolves
                .toEqual({
                    pseudo: 'pseudoNameUpdated',
                    description: 'I am a florian',
                    image: 'base64...',
                    id: 1
                });

        })

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

            const db:DB = factory();

            const p = profile(db);

            const testPromise = new Promise((res, rej) => {

                p
                    //Make sure that no profile is present
                    .hasProfile()

                    //Set Profile
                    .then(hasProfile => {
                        expect(hasProfile).toBeFalsy();
                        return p.setProfile('pedsa', 'i am a programmer', 'base64....');
                    })

                    //fetch profile
                    .then(_ => {
                        expect(_).toBeUndefined();
                        return p.getProfile();
                    })
                    .then(profile => {
                        res({
                            id: profile.id,
                            pseudo: profile.pseudo,
                            description: profile.description,
                            image: profile.image
                        });
                    })
                    .catch(err => rej(err))

            });

            return expect(testPromise)
                .resolves
                .toEqual({
                    pseudo: 'pedsa',
                    description: 'i am a programmer',
                    image: 'base64....',
                    id: 1
                });

        });

        /**
         * Try to fetch profile that doesn't exist
         */
        test('get profile that do not exist', () => {

            // Kill the database
            execSync('npm run db:flush');

            const db:DB = factory();

            const p = profile(db);

            return expect(p.getProfile())
                .rejects
                .toEqual(new NoProfilePresent())

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

            return expect(profile.getPublicProfile())
                .rejects
                .toEqual(new errors.NoPublicProfilePresent())

        });

        test('fetch my existing public profile', () => {

            return expect(profile.getPublicProfile())
                .resolves
                .toEqual({
                    pseudo: 'peasded',
                    description: 'I am a description',
                    image: 'base64....',

                    //Public eth addresses
                    ethAddresses: [
                        "0x2a65aca4d5fc5b5c859090a6c34d164135398226"
                    ],

                    //Mesh network keys (p2p lib identity)
                    meshKeys: [
                        "QmUKydZyhZmt2x5VpLJSojRarhRrC4k9QYpkYNf23sWy98"
                    ],

                    //Rsa key's for signing messages
                    identKey : [
                        "QmczUvgj46cvf4wWv8y3Z7RFiKyTUGSx3cSdL4Tqo5aevT"
                    ],

                    //Version of the profile
                    version: '1.0.0'
                });

        });

    });

    describe('hasProfile', () => {

        test('true', () => {

            const fakeQuery = () => {
                return [
                    //Since we count the object's returned by the query
                    //it's ok to return empty objects as a dummy
                    {}
                ]
            };

            return expect(hasProfile(factory(), fakeQuery)())
                .resolves
                .toBeTruthy();

        });

        test('false', () => {

            // Kill the database
            execSync('npm run db:flush');

            return expect(profile(factory()).hasProfile())
                .resolves
                .toBeFalsy();

        });

        test('error during fetch', () => {

            class TestError extends Error{}

            return expect(hasProfile(factory(), () => { throw new TestError()})())
                .rejects
                .toEqual(new TestError());

        });

    });

});

