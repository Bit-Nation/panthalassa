const secureStorageSpecification = require('./secureStorageSpecification');
const errors = require('./../errors');

/**
 * Tests the secure storage specification
 */
describe('secureStorageSpecification', () => {
    "use strict";

    /**
     * Test that check pass if all methods are implemented
     */
    test('satisfy specification', () => {

        // Required methods for the secure storage
        const requiredMethods = [
            'set',
            'get',
            'remove',
            'has',
            'isPasswordValid',
            'destroyStorage'
        ];

        // Expect that listed required methods are the same
        // than the exported once form secureStorageSpecification.requiredMethods
        expect(secureStorageSpecification.requiredMethods)
            .toBe(requiredMethods);

        // Create secure storage implementation
        const secureStorageImplementation = {};

        requiredMethods.forEach((method) => {

            secureStorageImplementation[method] = () => {};

        });

        // Expect that the created secure storage implementation is ok
        expect(secureStorageSpecification.check(secureStorageImplementation)).toBeTruthy();

    });

    /**
     * Test secure storage implementation with missing methods
     */
    test('does not satisfy specification', () => {

        expect(secureStorageSpecification({})).toThrow(errors.UnsatisfiedSecureStorageImplementationError)

    });

});