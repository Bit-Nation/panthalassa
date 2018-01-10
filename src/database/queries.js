// @flow

import type {ProfileObject} from './schemata';

/**
 *
 * @param {object} realm
 * @return {Realm.Results<any> | * | {$ref}}
 */
export function findProfiles(realm): Array<ProfileObject> {
    return realm.objects('Profile');
}
