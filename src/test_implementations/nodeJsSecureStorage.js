// @flow
import {SecureStorage} from '../specification/secureStorageInterface';

// ///////////////////////////////////////////////////////////////////////////////////
// This key value storage is not save. It's only intend to be used for development //
// ///////////////////////////////////////////////////////////////////////////////////

/**
 *
 * @return {SecureStorage}
 */
export default function(): SecureStorage {
    let storeage = {};

    const secureStorageImplementation : SecureStorage = {

        set: (key: string, value: any): Promise<void> => new Promise((res, rej) => {
            storeage[key] = value;
            res();
        }),

        get: (key: string): Promise<mixed> => new Promise((res, rej) => res(storeage[key])),

        remove: (key: string): Promise<void> => new Promise((res, rej): void => {
            delete storeage[key];

            res();
        }),

        has: (key: string): Promise<boolean> => new Promise((res, rej) => {
            if ('undefined' === typeof storeage[key]) {
                return res(false);
            }

            res(true);
        }),

        fetchItems: (filter: (key: string, value: mixed) => boolean): Promise<{}> => new Promise((res, rej) => {
            const filterdItems:{} = {};

            Object.keys(storeage).filter(function(key) {
               if (filter(key, storeage[key]) === true) {
                   filterdItems[key] = storeage[key];
               }
            });

            res(filterdItems);
        }),

        destroyStorage: () => new Promise((res, rej) => {
            storeage = {};

            res();
        }),

    };

    return secureStorageImplementation;
}
