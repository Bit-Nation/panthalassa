import nodeJsSecureStorage from "./nodeJsSecureStorage";

const fs = require('fs');
import './nodeJsSecureStorage';
import {set} from './nodeJsSecureStorage'

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

        const ss = nodeJsSecureStorage('./lib/'+Math.random());

        //Expect promise to resolve in undefined if the storage get's destroyed
        return expect(ss.destroyStorage()).resolves.toBeUndefined();

    });

    //Fetch items and filter for eth key's
    test('fetchItems', () => {

        const ss = nodeJsSecureStorage('./lib/'+Math.random());

        const filteredResults = ss
            .set('eth#1', 'eth_key_1')
            .then(() => {
                return ss.set('mesh#1', 'mesh_key_1');
            })
            .then(() => {
                return ss.set('eth#3', 'eth_key_3');
            })
            .then(() => {

                return ss.fetchItems((key, value) => {

                    return key.indexOf('eth') !== -1;

                });

            });

        return expect(filteredResults).resolves.toEqual([
            {
                key: 'eth#1',
                value: 'eth_key_1'
            },
            {
                key: 'eth#3',
                value: 'eth_key_3'
            }
        ]);

    })

});
