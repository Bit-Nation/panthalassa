const customRealm = require('./../../lib/database/realm');
const Realm = require('realm');

describe('realm', () => {
    "use strict";

    test('db return promise', () => {

        return expect(customRealm.db).toBeInstanceOf(Promise);

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

describe('query', () => {
    "use strict";

    test('successful', () => {

        //Since the realm db is nothing more than an object we can just create a fake obj
        const realmFakeDb = {"I am the realm fake db": true};

        //Test query
        const query = (realm) => {

            expect(realm).toBe(realmFakeDb);

            //Return some data which represent a query result
            return [
                'cat',
                'dog'
            ]

        };

        return expect(customRealm.raw.query(realmFakeDb)(query))
            .resolves
            .toEqual([
                'cat',
                'dog'
            ]);

    });

    test('with error', () => {

        class MyTestError extends Error{}

        //Since the realm db is nothing more than an object we can just create a fake obj
        const realmFakeDb = {"I am the realm fake db": true};

        //Test query that throws error
        const query = (realm) => {

            expect(realm).toBe(realmFakeDb);

            throw new MyTestError();

        };

        //Expect query promise rejection
        return expect(customRealm.raw.query(realmFakeDb)(query))
            .rejects
            .toEqual(new MyTestError());

    });

});
