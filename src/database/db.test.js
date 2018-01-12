/* eslint-disable */

import database from './db';
const execSync = require('child_process').execSync;

describe('write', () => {
    'use strict';

    /**
     * A database write should be void and will always result in a void promise
     */
    test('successfully', () => {
        // Kill the database
        execSync('npm run db:flush');

        const db = database();

        function writeAction(realm) {
            expect(realm).toBeDefined();

            return 'I wrote the profile';
        }

        return expect(db.write(writeAction))
            .resolves
            .toBeUndefined();
    });

    test('error', () => {
        // Kill the database
        execSync('npm run db:flush');

        class CustomError extends Error {}

        const db = database();

        function writeAction(realm) {
            expect(realm).toBeDefined();

            throw new CustomError();
        }

        return expect(db.write(writeAction))
            .rejects
            .toEqual(new CustomError);
    });
});

describe('query', () => {
    'use strict';

    test('successfully', () => {
        // Kill the database
        execSync('npm run db:flush');

        const db = database();

        function searchPets(realm) {
            expect(realm).toBeDefined();

            return [
                'dog',
                'cat',
            ];
        }

        return expect(db.query(searchPets))
            .resolves
            .toEqual(['dog', 'cat']);
    });

    test('error', () => {
        // Kill the database
        execSync('npm run db:flush');

        const db = database();

        class CustomError extends Error {}

        function searchPets(realm) {
            expect(realm).toBeDefined();


            throw new CustomError();
        }

        return expect(db.query(searchPets))
            .rejects
            .toEqual(new CustomError());
    });
});
