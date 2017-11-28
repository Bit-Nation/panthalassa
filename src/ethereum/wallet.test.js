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

        test('send eth successfully', () => {

            const fromAddress = '';

            const toAddress = '';

            return expect(fakeWallet.sendEth(fromAddress, toAddress, '1'))
                .resolves
                .toBeUndefined();

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