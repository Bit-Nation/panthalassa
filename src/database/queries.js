/**
 * Profile search query
 * @param realm
 */
function findProfiles(realm) {

    return realm.objects('Profile')

}

module.exports = {
    findProfiles
};