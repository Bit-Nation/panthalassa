const customRealm = require('./../../lib/database/realm');
const Realm = require('realm');

describe('realm', () => {
    "use strict";

    test('db return promise', () => {

        expect(customRealm.db).toBeInstanceOf(Promise);

    });

    /**
     * Test database write function. The write methods takes a function
     * that will be executed in the realm write context
     */
    test('db write action successfully', () => {

        return expect(customRealm.write((realm) => {})).resolves.toBeUndefined();

    });

    test('db write rejection', () => {

        class MyTestError extends Error{};

        return expect(customRealm.write((realm) => { throw new MyTestError() })).rejects.toEqual(new MyTestError);

    });

});