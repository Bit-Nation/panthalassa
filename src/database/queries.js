// @flow

import type {ProfileObject} from './schemata';

export default {
    findProfiles: (realm): Array<ProfileObject> => realm.objects('Profile'),
};
