// @flow

const errors = require('./../errors');

const requiredMethods: Array<string> = [
    'set',
    'get',
    'remove',
    'has',
    'isPasswordValid',
    'destroyStorage'
];

module.exports = (() : {} => {
    "use strict";
    return {
        requiredMethods,
        check: (implementation: {}) : boolean => {

            for(var c = 0; c < requiredMethods.length; c++){

                const checkMethod = requiredMethods[c];

                if('undefined' === typeof implementation[checkMethod]){
                    throw new errors.UnsatisfiedSecureStorageImplementationError(checkMethod);
                }

            }

            return true;

        }
    }
})();

