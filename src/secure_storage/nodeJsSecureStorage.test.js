const nodeJsSecureStorage = require('./../../lib/secure_storage/nodeJsSecureStorage');
const fs = require('fs');

describe('nodeJsSecureStorage', () => {
    "use strict";

    test('set', () => {

        const ss = nodeJsSecureStorage('./lib/'+Math.random());

        //Expecting to resolve in undefined since the set returns a void promise
        return expect(ss.set('name', 'Florian')).resolves.toBeUndefined();

    });

    test('get', () => {

        const ss = nodeJsSecureStorage('./lib/'+Math.random());

        const name = ss
            .set('name', 'Florian')
            .then(res => {
                return ss.get('name');
            });

        return expect(name).resolves.toBe('Florian');

    });

    test('remove', () => {

        const ss = nodeJsSecureStorage('./lib/'+Math.random());

        const remove = ss
            .set('key', 'value')
            .then(res => {
                expect(res).toBeUndefined();
                return ss
                    .remove('key')
            });


        //Expect the promise to resolve in undefined if it succeed
        return expect(remove).resolves.toBeUndefined();

    });

    test('has - true', () => {

        const ss = nodeJsSecureStorage('./lib/'+Math.random());

        //Expect a promise that resolves with true if the key is present
        return expect(ss
            .set('key', 'value')
            .then(res => {
                return ss
                    .has('key');
            })
        ).resolves.toBeTruthy();

    });

    test('has - false', () => {
        const ss = nodeJsSecureStorage('./lib/'+Math.random());

        //Expext the promise to resolve in false if key is not present
        return expect(ss.has('h')).resolves.toBeFalsy();
    });

    //Test successfully creation and deletion of secure storage
    test('destroyStorage', () => {

        const path = './lib/'+Math.random();

        const ss = nodeJsSecureStorage('./lib/'+Math.random());

        expect(fs.existsSync(path))
            .toBeTruthy();

        //Expect promise to resolve in undefined if the storage get's destroyed
        return expect(new Promise((res, rej) => {

            ss.destroyStorage()
                .then(result => {
                    res(result);
                })
                .catch(err => rej(err));

        })).resolves.toBeUndefinded();

    });

});
