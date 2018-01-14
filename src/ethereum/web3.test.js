import web3Factory from './web3';

describe('web3', () => {
    test('create web3 for offline usage with default address', (done) => {
        const ethUtilsMock = {
            allKeyPairs: () => new Promise((res, rej) => {
                const m = new Map();
                m.set('0xdb311bcf9a0e9b651006498ba69ad26e862c5a06', {});
                m.set('0x6589d76e6408e67c726903d5d777d18f03d184a1', {});
                res(m);
            }),
        };

        const nodeMock = {
            start: () => new Promise((res, rej) => res()),
        };

        web3Factory(nodeMock, ethUtilsMock, false)
            .then(function(web3) {
                // Should be first account from allKeyPairs method
                expect(web3.eth.defaultAccount).toBe('0xdb311bcf9a0e9b651006498ba69ad26e862c5a06');

                expect(web3.currentProvider).toBeUndefined();

                done();
            })
            .catch(done.fail);
    });

    test('create web3 for offline usage without an address', (done) => {
        const ethUtilsMock = {
            allKeyPairs: () => new Promise((res, rej) => res(new Map())),
        };

        const nodeMock = {
            start: () => new Promise((res, rej) => res()),
        };

        web3Factory(nodeMock, ethUtilsMock, false)
            .then(function(web3) {
                // Should be first account from allKeyPairs method
                expect(web3.eth.defaultAccount).toBeUndefined();

                expect(web3.currentProvider).toBeUndefined();

                done();
            })
            .catch(done.fail);
    });

    test('create web3 for online usage with default address', (done) => {
        const ethUtilsMock = {
            allKeyPairs: () => new Promise((res, rej) => {
                const m = new Map();
                m.set('0xdb311bcf9a0e9b651006498ba69ad26e862c5a06', {});
                m.set('0x6589d76e6408e67c726903d5d777d18f03d184a1', {});
                res(m);
            }),
        };

        const nodeMock = {
            start: () => new Promise((res, rej) => res()),
            url: 'http://',
        };

        web3Factory(nodeMock, ethUtilsMock, true)
            .then(function(web3) {
                // Should be first account from allKeyPairs method
                expect(web3.eth.defaultAccount).toBe('0xdb311bcf9a0e9b651006498ba69ad26e862c5a06');

                expect(web3.currentProvider).toBeDefined();

                done();
            })
            .catch(done.fail);
    });

    test('create web3 for offline usage without an address', (done) => {
        const ethUtilsMock = {
            allKeyPairs: () => new Promise((res, rej) => res(new Map())),
        };

        const nodeMock = {
            start: () => new Promise((res, rej) => res()),
            url: 'http://',
        };

        web3Factory(nodeMock, ethUtilsMock, true)
            .then(function(web3) {
                // Should be first account from allKeyPairs method
                expect(web3.eth.defaultAccount).toBeUndefined();

                expect(web3.currentProvider).toBeDefined();

                done();
            })
            .catch(done.fail);
    });
});
