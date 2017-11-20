describe('wallet', () => {
    "use strict";

    describe('ethBalance', () => {

        test('never synced', () => {

            const result = '';

            expect(result).toEqual({});

        });

        test('synced some time ago', () => {

            const result = '';

            expect(result).toEqual({
                addresss: '0x687422eEA2cB73B5d3e242bA5456b782919AFc85',
                synced_at : 1511185212,
                wei: '168179030063160961914893'
            })

        });

    })

});