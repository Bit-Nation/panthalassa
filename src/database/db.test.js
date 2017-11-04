const database = require('./db');

describe('write', () => {
    "use strict";

    /**
     * A database write should be void and will always result in a void promise
     */
    test('successfully', () => {

        const db = database.factory();

        function writeAction(realm){

            expect(realm).toBeDefined();

            return "I wrote the profile";

        }

        return expect(db.write(writeAction))
            .resolves
            .toBeUndefined();

    });

    test('error', () => {

        class CustomError extends Error{}

        const db = database.factory();

        function writeAction(realm){

            expect(realm).toBeDefined();

            throw new CustomError();

        }

        return expect(db.write(writeAction))
            .rejects
            .toEqual(new CustomError);


    });

});

describe('query', () => {
    "use strict";

    test('successfully', () => {

        const db = database.factory();

        function searchPets(realm){

            expect(realm).toBeDefined();

            return [
                'dog',
                'cat'
            ]

        }

        return expect(db.query(searchPets))
            .resolves
            .toEqual(['dog', 'cat']);

    });

    test('error', () => {

        const db = database.factory();

        class CustomError extends Error{}

        function searchPets(realm){

            expect(realm).toBeDefined();


            throw new CustomError();

        }

        return expect(db.query(searchPets))
            .rejects
            .toEqual(new CustomError());

    });

});
