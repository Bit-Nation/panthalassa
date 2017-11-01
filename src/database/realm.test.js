const customRealm = require('./../../lib/database/realm');
const Realm = require('realm');

describe('realm', () => {
    "use strict";

    test('db return promise', () => {

        expect(customRealm.db).toBeInstanceOf(Promise);

    })

});