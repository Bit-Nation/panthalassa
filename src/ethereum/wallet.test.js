/* eslint-disable */

import utils from './utils';
import {InvalidChecksumAddress} from './../errors';
import wallet from './wallet';
const BigNumber = require('bignumber.js');

const Web3 = require('web3');

describe('wallet', () => {
    'use strict';

    describe('ethBalance', () => {
        test('never synced', (done) => {
            const address = '0xfbb1b73c4f0bda4f67dca266ce6ef42f520fbb98';

            const filtered = jest.fn((filterString) => {
                expect(filterString).toBe(`id == "${address}_ETH"`);

                return [];
            });

            const realmMock = {
                objects: jest.fn((collection) => {
                    expect(collection).toBe('AccountBalance');

                    return {
                        filtered: filtered,
                    };
                }),
            };

            const dbMock = {
                query: (cb) => {
                    cb(realmMock);
                },
            };

            const u = utils();

            const w = wallet(u, null, dbMock);

            w.ethBalance(address)
                .then((_) => {
                    expect(_).toBeNull();

                    expect(realmMock.objects).toHaveBeenCalledTimes(1);
                    expect(filtered).toHaveBeenCalledTimes(1);

                    done();
                });
        });

        test('synced some time ago', (done) => {
            const address = '0x687422eEA2cB73B5d3e242bA5456b782919AFc85';

            const realmMock = {
                objects: jest.fn((collection) => {
                    expect(collection).toBe('AccountBalance');

                    return {
                        filtered: jest.fn(() => {
                            return [
                                {
                                    address: address,
                                    amount: '1000000000000000',
                                    synced_at: 1246624643444,
                                    currency: 'ETH',
                                    id: address+'_ETH',
                                },
                            ];
                        }),
                    };
                }),
            };

            const dbMock = {
                query: (cb) => {
                    cb(realmMock);
                },
            };

            const utl = utils();

            const w = wallet(utl, null, dbMock);

            w.ethBalance(address)
                .then((balance) => {
                    expect(balance).toEqual({
                        address: address,
                        amount: '1000000000000000',
                        synced_at: 1246624643444,
                        currency: 'ETH',
                        id: address+'_ETH',
                    });

                    done();
                });
        });

        test('invlid address', (done) => {

            const u = utils();

            const w = wallet(u);

            w.ethBalance('invalid_address')
                .catch((error) => {
                    expect(error).toBeInstanceOf(InvalidChecksumAddress);

                    done();
                });
        });
    });

    describe('ethSend', () => {
        test('send eth successfully', (done) => {

            const u = utils();

            const fromAddress = '0x9901C66F2d4b95F7074b553DA78084D708BECA70';

            const toAddress = '0x71d271f8B14adEf568F8f28f1587ce7271AC4Ca5';

            const txReceipt = {};

            const web3Mock = {
                eth: {
                    sendTransaction: jest.fn((txData, cb) => {
                        expect(txData.from).toBe(fromAddress);
                        expect(txData.to).toBe(toAddress);
                        expect(txData.value).toBe('1000000000000000000');
                        expect(txData.gasLimit).toBe(21000);
                        expect(txData.gasPrice).toBe(20000000000);

                        cb(null, txReceipt);
                    }),
                },
                toWei: jest.fn((eth) => {
                    const w3 = new Web3();

                    return w3.toWei(eth);
                }),
            };

            const w = wallet(u, web3Mock);

            w.ethSend(fromAddress, toAddress, '1')
                .then((txR) => {
                    expect(txR).toBe(txReceipt);

                    // sendTransaction should have been called since it's used to transfer eth
                    expect(web3Mock.eth.sendTransaction).toHaveBeenCalledTimes(1);

                    // toWei should have been called to transform eth to wei
                    expect(web3Mock.toWei).toHaveBeenCalledTimes(1);

                    done();
                });
        });

        test('failed to send eth', (done) => {
            class TestError extends Error {}

            const fromAddress = '0x9901C66F2d4b95F7074b553DA78084D708BECA70';

            const toAddress = '0x71d271f8B14adEf568F8f28f1587ce7271AC4Ca5';

            const web3Mock = {
                eth: {
                    sendTransaction: jest.fn((txData, cb) => {
                        expect(txData.from).toBe(fromAddress);
                        expect(txData.to).toBe(toAddress);
                        expect(txData.value).toBe('1000000000000000000');
                        expect(txData.gasLimit).toBe(21000);
                        expect(txData.gasPrice).toBe(20000000000);

                        cb(new TestError(), null);
                    }),
                },
                toWei: jest.fn((eth) => {
                    const w3 = new Web3();

                    return w3.toWei(eth, 'ether');
                }),
            };

            const u = utils();

            const w = wallet(u, web3Mock);

            w.ethSend(fromAddress, toAddress, '1')
                .catch((error) => {
                    // sendTransaction should have been called since it's used to transfer eth
                    expect(web3Mock.eth.sendTransaction).toHaveBeenCalledTimes(1);

                    // toWei should have been called to transform eth to wei
                    expect(web3Mock.toWei).toHaveBeenCalledTimes(1);

                    expect(error).toBeInstanceOf(TestError);

                    done();
                });
        });

        // Test if an invalid from address is reported.
        test('invalid from address', () => {

            const u = utils();

            const sendPromise = wallet(u).ethSend('I_AM_AN_INVALID_ADDRESS', null, '1');

            return expect(sendPromise).rejects.toEqual(new InvalidChecksumAddress('I_AM_AN_INVALID_ADDRESS'));
        });

        test('invalid to address', () => {

            const u = utils();

            const sendPromise = wallet(u).ethSend('I_AM_AN_INVALID_TO_ADDRESS', null, '1');

            return expect(sendPromise).rejects.toEqual(new InvalidChecksumAddress('I_AM_AN_INVALID_TO_ADDRESS'));
        });
    });

    describe('ethSync', () => {
        test('success', (done) => {
            const address = '0x9901c66f2d4b95f7074b553da78084d708beca70';

            const realm = {
                create: jest.fn((schemaName, data, update) => {
                    expect(update).toBe(true);
                    expect(schemaName).toBe('AccountBalance');

                    expect(data.id).toBe(address+'_ETH');
                    expect(data.address).toBe(address);
                    expect(data.currency).toBe('ETH');
                    expect('number' === typeof data.synced_at).toBeTruthy();
                    expect(data.amount).toBe('0.000001');
                }),
            };

            const dbMock = {
                write: jest.fn((cb) => {
                    cb(realm);
                }),
            };

            const web3 = new Web3();

            web3.eth.getBalance = jest.fn((addr, cb) => {
                expect(addr).toBe(address);

                cb(null, new BigNumber('1000000000000'));
            });

            wallet(utils(), web3, dbMock).ethSync(address)
                .then((_) => {
                    expect(web3.eth.getBalance).toHaveBeenCalledTimes(1);
                    expect(realm.create).toHaveBeenCalledTimes(1);
                    expect(dbMock.write).toHaveBeenCalledTimes(1);

                    done();
                })
                .catch((error) => {
                    throw error;
                });
        });

        test('invalid address', (done) => {
            wallet(utils()).ethSync('i_am_an_invalid_address')
                .catch((error) => {
                    expect(error).toBeInstanceOf(InvalidChecksumAddress);
                    done();
                });
        });

        test('error', (done) => {
            const address = '0x9901c66f2d4b95f7074b553da78084d708beca70';

            class TestError extends Error {}

            const web3Mock = {
                eth: {
                    getBalance: jest.fn((address, cb) => {
                        cb(new TestError(), null);
                    }),
                },
            };

            wallet(utils(), web3Mock).ethSync(address)
                .catch((error) => {
                    expect(error).toBeInstanceOf(TestError);

                    done();
                });
        });
    });
});
