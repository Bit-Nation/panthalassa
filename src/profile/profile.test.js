const profile = require('./profile');
const errors = require('./../../lib/errors');
const { spawn } = require('child_process');

describe('profile', () => {
    "use strict";

    describe('setProfile', () => {

        /**
         * A profile has three parameters, a pseudo name, a description and a image
         * This is a functional test
         */
        test('create profile', () => {

            // Kill the database
            spawn.execSync('npm run db:flush');

            return expect(new Promise((res, rej) => {

                profile
                    .setProfile('pseudoName', 'I am a florian', 'base64...')
                    .then(result => {
                        return profile.getProfile();
                    })
                    .then(profile => {
                        res(profile);
                    })
                    .catch(err => {
                        rej(err);
                    })

            }))
                .resvoles
                .toEqual({
                    pseudo: 'pseudoName',
                    description: 'I am a florian',
                    image: 'base64...'
                })

        });

        test('try to update profile', () => {

            // Kill the database
            spawn.execSync('npm run db:flush');

            const profilePromise = profile

                //Create profile
                .setProfile('pseudoName', 'I am a florian', 'base64...')

                //After saving fetch it and return the promise
                .then(_ => {

                    return profile.getProfile();

                })

                //The profile content should match since get profile
                //fetched the profile
                .then(profileContent => {

                    //Assert profile matches
                    expect(profileContent).toEqual({
                        pseudo: 'pseudoName',
                        description: 'I am a florian',
                        image: 'base64...'
                    });

                    //Update profile after ensure that it was written to the db
                    return profile
                        .setProfile({
                            pseudo: 'pseudoNameUpdated',
                            description: 'I am a florian',
                            image: 'base64...'
                        });

                })

                // Fetch the updated profile
                .then(_ => {

                    return profile.getProfile()

                });

            return expect(profilePromise)
                .resvoles
                .toEqual({
                    pseudo: 'pseudoNameUpdated',
                    description: 'I am a florian',
                    image: 'base64...'
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

            return expect(profile.getProfile())
                .resolves.toEqual({
                    pseudo: 'pedsa',
                    description: 'i am a programmer',
                    image: 'base64....',
                    version: '1.0.0'
                });

        });

        /**
         * Try to fetch profile that doesn't exist
         */
        test('get profile that do not exist', () => {

            return expect(profile.getProfile())
                .rejects
                .toEqual(new errors.NoProfilePresent())

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

});

