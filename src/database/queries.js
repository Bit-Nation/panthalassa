/**
 * Profile search query
 * @param realm
 */
import type {ProfileObject} from "./schemata";

export function findProfiles(realm) : Array<ProfileObject> {

    return realm.objects('Profile')

}
