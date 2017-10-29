const errors = require('./../lib/errors');

describe('errors', () => {
    "use strict";

    /**
     * Test if error is exported
     */
    test('InvalidPrivateKeyError', () => {

        expect(errors.InvalidPrivateKeyError).toBeDefined();

    });

    /**
     * Test if error is exported and if error message is correct
     */
    test('UnsatisfiedSecureStorageImplementationError', () => {

        expect(errors.UnsatisfiedSecureStorageImplementationError).toBeDefined();

        const error = new errors.UnsatisfiedSecureStorageImplementationError('foo');

        expect(error.message)
            .toBe('Missing method: "foo" in secure storage implementation');

    });

});