/* eslint-disable */

import {getAccounts, signTx} from './PanthalassaProvider';
import type {PrivateKeyType} from '../specification/privateKey';
const EthTx = require('ethereumjs-tx');

describe('getAccounts', () => {
    test('success', (done) => {
        const address_one = '0x465868366a0f45748f24d8979a98c2118e71b2bc';

        const address_two = '0x26e75307fc0c021472feb8f727839531f112f317';

        const ethUtils = {
            allKeyPairs: () => new Promise((res, rej) => {

                const m = new Map();
                m.set('0x465868366a0f45748f24d8979a98c2118e71b2bc', '');
                m.set('0x26e75307fc0c021472feb8f727839531f112f317', '');

                res(m);

            }),
        };

        const cb = (error, addresses) => {
            expect(error).toBeNull();

            expect(addresses).toEqual([
                address_one,
                address_two,
            ]);

            // This will mark the test as done
            done();
        };

        getAccounts(ethUtils)(cb);
    });

    test('error', (done) => {
        class TestError extends Error {}

        const ethUtils = {
            allKeyPairs: () => {
                return new Promise((res, rej) => {
                    // Reject promise with test error
                    rej(new TestError());
                });
            },
        };

        const cb = (error, addresses) => {
            expect(error).toEqual(new TestError());

            expect(addresses).toBeNull();

            // This will mark the test as done
            done();
        };

        getAccounts(ethUtils)(cb);
    });
});

describe('signTx', () => {
    // Sample tx data
    const txDataMock = {
        nonce: '0x03',
        gas: '0x5208',
        from: '0xae481410716b6d087261e0d69480b4cb9305c624',
        to: '0x814944ed940f27eb40330882a24baad21c30818e',
        value: '0x1',
        gasPrice: '0x4a817c800',
    };

    // sample private key
    const privateKey:string = 'affd0b4039708432bb2759fc747bf7b9b1fbdab71bf86eab6d812ae83419b708';

    // signed transaction
    const signedTxStrMock = 'f864038504a817c80082520894814944ed940f27eb40330882a24baad21c30818e01801ba063a5002e8054f7c95e4520ad4ef7739e8d66adc3a11d511b53b15388d6cd8c84a0212ccf0f79cc23a1f53aa8f90e8210633bceb2c85d6797bd0acfdec874c5b092';

    test('success', (done) => {
        // Mock signed transaction
        const signedTxMock = new EthTx(txDataMock);
        signedTxMock.sign(Buffer.from(privateKey, 'hex'));

        // Mock eth utils
        const ethUtils = {

            getPrivateKey: (address) => {
                expect(address).toBe(txDataMock.from);

                return new Promise((res, rej) => {
                    // Return PrivateKeyType json object
                    res({
                        encryption: '',
                        value: privateKey,
                        encrypted: false,
                        version: '1.0.0',
                    });
                });
            },

            signTx: (txData, privateKey) => {
                expect(txData).toBe(txDataMock);
                expect(privateKey).toBe(privateKey);

                return new Promise((res, rej) => {
                    res(signedTxMock);
                });
            },

        };

        // Mock callback
        const cb = (error, signedTx) => {
            // No error since this is a success test
            expect(error).toBeNull();

            // Expect hex string
            expect(signedTx).toBe('0x'+signedTxStrMock);

            done();
        };

        signTx(ethUtils)(txDataMock, cb);
    });

    /**
     * Test error handling
     */
    describe('error', () => {
        /**
         * Test error handling when smth goes wrong in the ethUtils getPrivateKey method
         */
        test('EthUtils - getPrivateKey', (done) => {
            class TestErrorGetPrivateKey extends Error {}

            // Mock eth utils
            const ethUtils = {

                getPrivateKey: (address) => {
                    expect(address).toBe(txDataMock.from);

                    return new Promise((res, rej) => {
                        rej(new TestErrorGetPrivateKey());
                    });
                },

            };

            // Mock callback
            const cb = (error, signedTx) => {
                // No error since this is a success test
                expect(error).toEqual(new TestErrorGetPrivateKey());

                // Expect hex string
                expect(signedTx).toBeNull();

                done();
            };

            signTx(ethUtils)(txDataMock, cb);
        });

        /**
         * Test error handling when smth goes wrong in the ethUtils signTx method
         */
        test('EthUtils - signTx', (done) => {
            class FailedToSignError extends Error {}

            // Mock eth utils
            const ethUtils = {

                getPrivateKey: (address) => {
                    expect(address).toBe(txDataMock.from);

                    return new Promise((res, rej) => {
                        // Return PrivateKeyType json object
                        res({
                            encryption: '',
                            value: privateKey,
                            encrypted: false,
                            version: '1.0.0',
                        });
                    });
                },

                signTx: (txData, privateKey) => {
                    expect(txData).toBe(txDataMock);
                    expect(privateKey).toBe(privateKey);

                    return new Promise((res, rej) => {
                        rej(new FailedToSignError());
                    });
                },

            };

            // Mock callback
            const cb = (error, signedTx) => {
                // No error since this is a success test
                expect(error).toBeInstanceOf(FailedToSignError);

                // Expect hex string
                expect(signedTx).toBeNull();

                done();
            };

            signTx(ethUtils)(txDataMock, cb);
        });
    });
});
