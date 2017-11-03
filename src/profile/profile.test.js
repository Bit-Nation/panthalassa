const profile = require('./profile');
const errors = require('./../../lib/errors');

describe('profile', () => {
    "use strict";

    /**
     * A profile has three parameters, a pseudo name, a description and a image
     */
    test('setProfile', () => {

        return expect(profile.setProfile('pseudoName', 'I am a florian', 'base64...')).resvoles.toBeUndefined();

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
                .toBe(new errors.NoProfilePresent())

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
                    pseudo: '',
                    description: '',
                    image: 'base64....',
                    ethAddresses: [
                        "0x2a65aca4d5fc5b5c859090a6c34d164135398226"
                    ],
                    meshKeys: [
                        "QmUKydZyhZmt2x5VpLJSojRarhRrC4k9QYpkYNf23sWy98"
                    ],
                    ids : [
                        "QmczUvgj46cvf4wWv8y3Z7RFiKyTUGSx3cSdL4Tqo5aevT"
                    ],
                    version: '1.0.0'
                });

        });

    });

});

