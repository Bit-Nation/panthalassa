//@flow

/**
 *
 * @param realm
 * @returns {function()}
 */
const hasProfile = (realm:any) : (() => Promise<boolean>) => {
    "use strict";
    return () : Promise<boolean> => {

        return new Promise((res, rej) => {

            try {
                const profiles = realm.objects('Profile');

                if(profiles.length >= 1){
                    res(true);
                    return;
                }

                res(false);

            }catch (e){
                rej(e);
            }

        });

    }
};

/**
 *
 * @param realm
 * @returns {{hasProfile}}
 */
module.exports = (realm:any) => {
    "use strict";
    return {
        hasProfile: hasProfile(realm),
        raw: {
            hasProfile
        }
    }
};
