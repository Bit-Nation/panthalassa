// @flow

import type {ProfileObject} from './schemata';
const Realm = require('realm');

export default {
    findProfiles: (realm: Realm): Array<ProfileObject> => realm.objects('Profile'),
};
