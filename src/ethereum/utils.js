//@flow

import type {SecureStorage} from "../specification/secureStorageInterface";
import type { PrivateKeyType } from '../specification/privateKey'
import type {TxData} from '../specification/tx';
import {AbortedSigningOfTx, InvalidPrivateKeyError, InvalidChecksumAddress} from "../errors";
import type {OsDependenciesInterface} from "../specification/osDependencies";

const crypto = require('crypto-js');
const ethereumjsUtils = require('ethereumjs-util');
const errors = require('./../errors');
const aes = require('crypto-js/aes');
const EventEmitter = require('eventemitter3');
const EthTx = require('ethereumjs-tx');
const bip39 = require('bip39');
const ethJsUtils = require('ethereumjs-util');
const assert = require('assert');

const PRIVATE_ETH_KEY_PREFIX = 'PRIVATE_ETH_KEY#';

/**
 * Ethereum Utils Interface
 */
export interface EthUtilsInterface {

    /**
     * Creates private key and return as hex string
     */
    createPrivateKey: () => Promise<string>,

    /**
     * Save private key with an optional password
     */
    savePrivateKey: (privateKey:string, pw:?string, pwConfirm:?string) => Promise<void>,

    //@todo change this method and the doc's
    allKeyPairs: () => Promise<{}>,

    /**
     * Fetch private key by address. Make sure to normalize the address.
     * Will be rejected if private key was not found.
     */
    getPrivateKey: (address:string) => Promise<PrivateKeyType>,

    /**
     * Delete private key by address. Make sure to normalize the address.
     */
    deletePrivateKey: (address:string) => Promise<void>,

    /**
     * This method decrypt's an private key. Have a look at the readme in this folder
     * to see how to use this method.
     */
    decryptPrivateKey: (privateKey: PrivateKeyType, reason: string, topic: string) => Promise<string>,

    /**
     * Sign eth transaction data. have a look at the readme in this folder to see
     * how to use this method.
     */
    signTx: (txData:TxData, privkey:string) => Promise<EthTx>,

    /**
     * Normalize an ethereum address
     */
    normalizeAddress: (address:string) => string,

    /**
     * Normalize an ethereum private key
     */
    normalizePrivateKey: (privateKey:string) => string,

    /**
     * Transform private key to list of words
     */
    privateKeyToMnemonic: (pk:string) => Array<string>,

    /**
     * Mnemonic to private key
     */
    mnemonicToPrivateKey: (mnemonic:string) => string,

    /**
     * Validates a mnemonic
     */
    mnemonicValid: (mnemonic:string) => boolean

}

/**
 *
 * @param address
 * @returns {string}
 */
export function normalizeAddress(address:string) : string {

    const checksumAddress:string = ethereumjsUtils.toChecksumAddress(address);

    if(!ethereumjsUtils.isValidChecksumAddress(checksumAddress)){
        throw new InvalidChecksumAddress(address);
    }

    return checksumAddress;

}

/**
 *
 * @param privateKey
 * @returns {string}
 */
export function normalizePrivateKey(privateKey:string) : string {

    if(!ethereumjsUtils.isValidPrivate(Buffer.from(privateKey, 'hex'))){
        throw new InvalidPrivateKeyError();
    }

    return privateKey;
}

/**
 * @description Generate's a new secure private key
 * @param osDeps
 * @param isValidPrivateKey
 * @return {function()}
 */
export function createPrivateKey(osDeps:OsDependenciesInterface, isValidPrivateKey: (key: Buffer) => boolean) : (() => Promise<string>){
    "use strict";

    return () : Promise<string> => {

        return new Promise((res, rej) => {

            osDeps.crypto.randomBytes(32)
                .then(privateKey => {

                    if(!isValidPrivateKey(Buffer.from(privateKey, 'hex'))){
                        return rej(new errors.InvalidPrivateKeyError());
                    }

                    return res(privateKey);

                })
                .catch(rej);


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

            privateKey = normalizePrivateKey(privateKey);

            const privateKeyBuffer = Buffer.from(privateKey, 'hex');

            //Reject promise if private key is not a valid hey private key
            if(!ethjsUtils.isValidPrivate(privateKeyBuffer)){

                rej(new errors.InvalidPrivateKeyError);
                return;

            }

            const addressOfPrivateKey = normalizeAddress(ethjsUtils.privateToAddress(privateKeyBuffer).toString('hex'));

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
export function allKeyPairs(secureStorage:SecureStorage) : (() => Promise<{}>){
    "use strict";

    return () : Promise<{}> => new Promise((res, rej) => secureStorage
        .fetchItems((key:string) => key.indexOf(PRIVATE_ETH_KEY_PREFIX) !== -1)
        .then(keys => {

            const transformedKeys = {};

            Object.keys(keys).map(key => {

                //We only accept string's since. the private key is an stringified object
                if(typeof keys[key] !== 'string'){
                    return rej(new Error(`Value of key: '${key}' is not an string`));
                }

                transformedKeys[key.split(PRIVATE_ETH_KEY_PREFIX).pop()] = JSON.parse(keys[key])

            });

            res(transformedKeys);

        })
        .catch(err => rej(err)))

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

            const key = PRIVATE_ETH_KEY_PREFIX+normalizeAddress(address);

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

            const key = PRIVATE_ETH_KEY_PREFIX+normalizeAddress(address);

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
export function decryptPrivateKey(pubEE:EventEmitter, crypto: any, ethjsUtils: ethereumjsUtils): ((privateKey: {value: string}, reason: string, topic: string) => Promise<string>){
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

/**
 * Sign a transaction
 * @param isPrivateKey
 * @param ee
 * @returns {function(TxData, string)}
 */
export function signTx(isPrivateKey: (privKey:Buffer) => boolean, ee: EventEmitter) : (txData:TxData, privKey:string) => Promise<EthTx> {

    return (txData:TxData, privKey:string) : Promise<EthTx> => new Promise((res, rej) => {

        //Private key as buffer
        const pKB = Buffer.from(privKey, 'hex');

        //reject if private key is invalid
        if(!isPrivateKey(pKB)){
            return rej(new InvalidPrivateKeyError());
        }

        //Sign transaction
        const tx = new EthTx(txData);

        /**
         * client need's to react to this event
         * in order to sign the transaction
         */
        ee.emit('eth:tx:sign', {
            tx: tx,
            txData: txData,
            confirm: () => {
                tx.sign(pKB);
                res(tx);
            },
            abort: () => rej(new AbortedSigningOfTx())
        });

    })

}

/**
 *
 * @param ss
 * @param ee
 * @param osDeps
 * @return {EthUtilsInterface}
 */
export default function (ss:SecureStorage, ee:EventEmitter, osDeps:OsDependenciesInterface) : EthUtilsInterface {

    const ethUtilsImplementation:EthUtilsInterface = {
        createPrivateKey: createPrivateKey(osDeps, ethereumjsUtils.isValidPrivate),
        savePrivateKey: savePrivateKey(ss, ethereumjsUtils, aes),
        allKeyPairs: allKeyPairs(ss),
        getPrivateKey: getPrivateKey(ss),
        deletePrivateKey: deletePrivateKey(ss),
        decryptPrivateKey: decryptPrivateKey(ee, crypto, ethereumjsUtils),
        signTx: signTx(ethereumjsUtils.isValidPrivate, ee),
        normalizeAddress: normalizeAddress,
        normalizePrivateKey: normalizePrivateKey,
        privateKeyToMnemonic: (privateKey:string) : Array<string> => {

            assert.equal(true, ethJsUtils.isValidPrivate(Buffer.from(privateKey, 'hex')), 'Expected valid private key');

            return bip39.entropyToMnemonic(privateKey).split(' ');

        },
        mnemonicToPrivateKey: (mnemonic:string) : string => bip39.mnemonicToEntropy(mnemonic),
        mnemonicValid: bip39.validateMnemonic
    };

    return ethUtilsImplementation;

}
