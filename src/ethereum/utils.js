//@flow

import type {SecureStorage} from "../specification/secureStorageInterface";
import type { PrivateKeyType } from '../specification/privateKey'

const crypto = require('crypto');
const ethereumjsUtils = require('ethereumjs-util');
const errors = require('./../errors');
const aes = require('crypto-js/aes');
const eventEmitter = require('eventemitter3');

const PRIVATE_ETH_KEY_PREFIX = 'PRIVATE_ETH_KEY#';

/**
 * Creates a new private key
 * @param crypto
 * @param isValidPrivateKey
 * @returns {function()}
 */
export function createPrivateKey(crypto:{...any}, isValidPrivateKey: (key: Buffer) => boolean) : (() => Promise<string>){
    "use strict";

    return () : Promise<string> => {

        return new Promise((res, rej) => {

            crypto.randomBytes(32, function(err, privKey){

                if(err){
                    rej(err);
                }

                if(!isValidPrivateKey(privKey)){
                    rej(new errors.InvalidPrivateKeyError());
                }

                res(privKey.toString('hex'));

            });

        })

    }

}

/**
 *
 * @param secureStorage
 * @param ethjsUtils
 * @param aes
 * @returns {function(string, string, string)}
 */
export function savePrivateKey(secureStorage: SecureStorage, ethjsUtils: ethereumjsUtils, aes: aes) : ((privateKey:string, pw:?string, pwConfirm:?string) => Promise<void>){
    "use strict";

    return (privateKey: string, pw: ?string, pwConfirm: ?string) : Promise<void> => {

        return new Promise((res, rej) => {

            //Reject promise if private key is not a valid hey private key
            if(!ethjsUtils.isValidPrivate(Buffer.from(privateKey, 'hex'))){

                rej(new errors.InvalidPrivateKeyError);
                return;

            }

            privateKey = ethjsUtils.addHexPrefix(privateKey);

            const addressOfPrivateKey = ethjsUtils
                    .toChecksumAddress(ethjsUtils.privateToAddress(privateKey)
                    .toString('hex'));

            //Reject promise if one of the passwords is entered AND if they don't match
            if('undefined' !== typeof pw || 'undefined' !== typeof pwConfirm){

                if(pw !== pwConfirm){
                    rej(new errors.PasswordMismatch);
                    return;
                }

                //Special chars mach pattern
                const specialCharsPattern = /[ !@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/;

                // $FlowFixMe From a logical point of view the password can't be null / undefined here
                if(specialCharsPattern.test(pw) || specialCharsPattern.test(pwConfirm)){
                    rej(new errors.PasswordContainsSpecialChars());
                    return;
                }

                const pk:PrivateKeyType = {
                    encryption: 'AES-256',
                    value: aes.encrypt(privateKey, pw).toString(),
                    encrypted: true,
                    version: '1.0.0'
                };

                //Save the private key
                secureStorage.set(
                    PRIVATE_ETH_KEY_PREFIX+addressOfPrivateKey,
                    JSON.stringify(pk)
                )
                    .then(result => res(result))
                    .catch(err => rej(err));

            }

            //Save the private key
            //@Todo make the json data set a type (maybe)
            secureStorage.set(
                PRIVATE_ETH_KEY_PREFIX+addressOfPrivateKey,
                JSON.stringify({
                    encryption : '',
                    value: privateKey,
                    encrypted: false,
                    version: '1.0.0'
                })
            )
                .then(result => res(result))
                .catch(err => rej(err));

        });

    };

}

/**
 * Fetch all keyPairs
 * @param secureStorage
 * @returns {function()}
 */
export function allKeyPairs(secureStorage:SecureStorage) : (() => Promise<*>){
    "use strict";

    return () : Promise<*> => {

        return new Promise((res, rej) => {
            secureStorage

                .fetchItems((key:string, value:string) => {

                    // Filter eth keys
                    return key.indexOf(PRIVATE_ETH_KEY_PREFIX) !== -1;

                })
                .then(keyValuePairs => {

                    //transform key value pairs. remove key eth prefix and transform json string to json
                    keyValuePairs
                        .map((keyValuePair) => {
                            //Transform keypair
                            keyValuePair.key.split(PRIVATE_ETH_KEY_PREFIX).pop();
                            keyValuePair.value = JSON.parse(keyValuePair.value)
                        });

                })
                .catch(err => rej(err));
        })

    }

}

/**
 * Fetches a private key based on the
 * @param secureStorage
 * @returns {function(string)}
 */
export function getPrivateKey(secureStorage:SecureStorage) : ((address:string) => Promise<PrivateKeyType>){
    "use strict";

    return (address:string) : Promise<{...any}> => {

        return new Promise((res, rej) => {

            const key = PRIVATE_ETH_KEY_PREFIX+address;

            secureStorage
                .has(key)
                .then(hasPrivateKey => {

                    if(false === hasPrivateKey){
                        rej(new errors.NoEquivalentPrivateKey());
                        return;
                    }

                    return secureStorage.get(key);

                })
                .then(function(privKey:any){
                    res(JSON.parse(privKey))
                })
                .catch(err => rej(err));

        });

    }

}

/**
 *
 * @param secureStorage
 * @returns {function(string)}
 */
export function deletePrivateKey(secureStorage:SecureStorage) : ((address:string) => Promise<void>){
    "use strict";

    return (address:string) : Promise<void> => {

        return new Promise((res, rej) => {

            const key = PRIVATE_ETH_KEY_PREFIX+address;

            secureStorage
                .has(key)
                .then(hasPrivateKey => {

                    if(false === hasPrivateKey){
                        rej(new errors.NoEquivalentPrivateKey());
                        return;
                    }

                    return secureStorage.remove(key);

                })
                .then(result => res(result))
                .catch(err => rej(err));

        });

    }

}

/**
 * Decrypt the private key. Will emit an event that contains method's to solve this problem
 * @param pubEE
 * @param crypto
 * @param ethjsUtils
 * @returns {function({}, string, string)}
 */
export function decryptPrivateKey(pubEE:eventEmitter, crypto: any, ethjsUtils: ethereumjsUtils): ((privateKey: {value: string}, reason: string, topic: string) => Promise<string>){
    "use strict";

    return (privateKey: PrivateKeyType, reason:string, topic: string) : Promise<string> => {

        return new Promise((mRes, mRej) => {

            //break if the algo is unknown
            if(privateKey.encryption !== 'AES-256'){
                mRej(new errors.InvalidEncryptionAlgorithm());
                return;
            }

            //Call this to decrypt the password
            function successor(pw:string) : Promise<void>{

                return new Promise((res, rej) => {

                    const decryptedPrivateKey = crypto
                        .AES
                        .decrypt(privateKey.value.toString(), pw)
                        .toString(crypto.enc.Utf8);

                    //When aes decryption failes a empty string is returned
                    if('' === decryptedPrivateKey){
                        rej(new errors.FailedToDecryptPrivateKeyPasswordInvalid);
                        return;
                    }

                    //Check if decrypted key is valid
                    if(!ethjsUtils.isValidPrivate(Buffer.from(decryptedPrivateKey, 'hex'))){
                        rej(new errors.DecryptedValueIsNotAPrivateKey());
                        return;
                    }

                    res();
                    mRes(decryptedPrivateKey);

                });

            }

            //Call this to kill the decryption proccess
            const killer = () => {
                mRej(new errors.CanceledAction());
            };

            pubEE.emit('eth:decrypt-private-key', {
                successor: successor,
                killer: killer,
                reason: reason,
                topic: topic
            })

        });

    }

}

export interface EthUtils {
    createPrivateKey: () => Promise<string>,
    savePrivateKey: (privateKey:string, pw:?string, pwConfirm:?string) => Promise<void>,
    allKeyPairs: () => Promise<*>,
    getPrivateKey: (address:string) => Promise<{...any}>,
    deletePrivateKey: (address:string) => Promise<void>,
    decryptPrivateKey: (privateKey: {value: string}, reason: string, topic: string) => Promise<string>
}

/**
 * Returns eth utils implementation
 * @param ss
 * @param ee
 * @returns {EthUtils}
 */
export default function (ss:SecureStorage, ee:eventEmitter) : EthUtils {

    const ethUtilsImplementation:EthUtils = {
        createPrivateKey: createPrivateKey(crypto, ethereumjsUtils.isValidPrivate),
        savePrivateKey: savePrivateKey(ss, ethereumjsUtils, aes),
        allKeyPairs: allKeyPairs(ss),
        getPrivateKey: getPrivateKey(ss),
        deletePrivateKey: deletePrivateKey(ss),
        decryptPrivateKey: decryptPrivateKey(ee, crypto, ethereumjsUtils)
    };

    return ethUtilsImplementation;

}
