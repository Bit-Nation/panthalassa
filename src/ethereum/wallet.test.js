const web3 = require('web3');
import {ethSend} from './wallet'
//@Todo Replace this with the real wallet. This is just a dummy
const fakeWallet = {
    ethSend : (from, to, amount) => {
        "use strict";

    },
    ethBalance: (address) => {
        "use strict";

    },
    ethSync: (address) => {
        "use strict";

    },
    syncCurrencies: (address) => {

    }
};

describe('wallet', () => {
    "use strict";

    describe('ethBalance', () => {

        test('never synced', () => {

            const address = '';

            //Will be empty object it was not synchronised before
            return expect(fakeWallet.ethBalance(address)).toEqual({});

        });

        test('synced some time ago', () => {

            const address = '0x687422eEA2cB73B5d3e242bA5456b782919AFc85';

            return expect(fakeWallet.ethBalance(address)).toEqual({
                synced_at : 1511185212,
                wei: '168179030063160961914893',
                address: '0x687422eEA2cB73B5d3e242bA5456b782919AFc85'
            })

        });

    });

    describe('ethSend', () => {

        test('send eth successfully', done => {

            const fromAddress = '0x9901c66f2d4b95f7074b553da78084d708beca70';

            const toAddress = '0x71d271f8b14adef568f8f28f1587ce7271ac4ca5';

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

                    })
                },
                utils: {
                    toWei: jest.fn((eth) => {

                        return web3.utils.toWei(eth, 'ether');

                    })
                }
            };

            const ethUtilsMock = {
                normalizeAddress: (address) => address
            };

            const sendPromise = ethSend(ethUtilsMock, web3Mock)(fromAddress, toAddress, '1');

            sendPromise
                .then(txR => {

                    expect(txR).toBe(txReceipt);

                    //sendTransaction should have been called since it's used to transfer eth
                    expect(web3Mock.eth.sendTransaction).toHaveBeenCalledTimes(1);

                    //toWei should have been called to transform eth to wei
                    expect(web3Mock.utils.toWei).toHaveBeenCalledTimes(1);

                    done();

                })

        });

        test('failed to send eth', () => {

            class TestError extends Error{}

            const fromAddress = '';

            const toAddress = '';

            //The error will be from web3
            return expect(fakeWallet.sendEth(fromAddress, toAddress, '1'))
                .resolves
                .toBe(new TestError());

        })

    });

    describe('ethSync', () => {

        test('success', () => {

            const address = '';

            //Will be resolved as "void" if successfull synced
            return expect(fakeWallet.ethSync(address))
                .resolves
                .toBeUndefined();

        });

        test('error', () => {

            const address = '';

            class TestError extends Error{}

            //Will be resolved with error that was thrown by other code
            return expect(fakeWallet.ethSync(address))
                .resolves
                .toEqual(new TestError());

        });

    });

    test('syncCurrencies', () => {

        const address = '';

        //syncCurrencies sync's eth and pat. Expect to get back the
        //result of ethSync and patSync
        return expect(fakeWallet.syncCurrencies(address))
            .resolves
            .toEqual([
                undefined,
                undefined
            ])

    })

});