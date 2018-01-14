// @flow

import type {SecureStorageInterface} from '../specification/secureStorageInterface';
import type {PrivateKeyType} from '../specification/privateKey';
import type {TxData} from '../specification/tx';
import {AbortedSigningOfTx, InvalidPrivateKeyError, InvalidChecksumAddress} from '../errors';
import type {OsDependenciesInterface} from '../specification/osDependencies';

const crypto = require('crypto-js');
const ethereumjsUtils = require('ethereumjs-util');
const errors = require('./../errors');
const aes = require('crypto-js/aes');
const EventEmitter = require('eventemitter3');
const EthTx = require('ethereumjs-tx');
const bip39 = require('bip39');
const ethJsUtils = require('ethereumjs-util');
const assert = require('assert');
import {AMOUNT_OF_ADDRESSES_CHANGED} from '../events'

const PRIVATE_ETH_KEY_PREFIX = 'PRIVATE_ETH_KEY#';

/**
 * @typedef EthUtilsInterface
 * @property {function} createPrivateKey
 * @property {function(privateKey:string, pw: ?string, pwConfirm: ?string)} savePrivateKey
 * @property {function} allKeyPairs returns @todo
 * @property {function(address:string) : PrivateKeyType} getPrivateKey fetch private key by address
 * @property {function(address:string)} deletePrivateKey
 * @property {function(privateKey: PrivateKeyType, reason: string, topic: string)} decryptPrivateKey
 * @property {function(txData: TxData, privkey: string)} signTx
 * @property {function(address:string)} normalizeAddress
 * @property {function(pk: string)} privateKeyToMnemonic
 * @property {function(mnemonic: string)} mnemonicToPrivateKey
 * @property {function(mnemonic: string)} mnemonicValid
 *
 */
export interface EthUtilsInterface {

    /**
     * Creates private key and return as hex string
     */
    createPrivateKey: () => Promise<string>,

    /**
     * Save private key with an optional password
     */
    savePrivateKey: (privateKey: string, pw: ?string, pwConfirm: ?string) => Promise<void>,

    // @todo change this method and the doc's
    allKeyPairs: () => Promise<Map<string, PrivateKeyType>>,

    /**
     * Fetch private key by address. Make sure to normalize the address.
     * Will be rejected if private key was not found.
     */
    getPrivateKey: (address: string) => Promise<PrivateKeyType>,

    /**
     * Delete private key by address. Make sure to normalize the address.
     */
    deletePrivateKey: (address: string) => Promise<void>,

    /**
     * This method decrypt's an private key. Have a look at the readme in this folder
     * to see how to use this method.
     */
    decryptPrivateKey: (privateKey: PrivateKeyType, reason: string, topic: string) => Promise<string>,

    /**
     * Sign eth transaction data. have a look at the readme in this folder to see
     * how to use this method.
     */
    signTx: (txData: TxData, privkey: string) => Promise<EthTx>,

    /**
     * Normalize an ethereum address
     */
    normalizeAddress: (address: string) => string,

    /**
     * Normalize an ethereum private key
     */
    normalizePrivateKey: (privateKey: string) => string,

    /**
     * Transform private key to list of words
     */
    privateKeyToMnemonic: (pk: string) => Array<string>,

    /**
     * Mnemonic to private key
     */
    mnemonicToPrivateKey: (mnemonic: string) => string,

    /**
     * Validates a mnemonic
     */
    mnemonicValid: (mnemonic: string) => boolean

}

/**
 * @name ethereum/utils.js
 * @desc ethereum utils factory
 * @param {object} ss secure storage
 * @param {object} ee event emitter
 * @param {object} osDeps operating system dependencies
 * @return {EthUtilsInterface}
 */
export default function utilsFactory(ss: SecureStorageInterface, ee: EventEmitter, osDeps: OsDependenciesInterface): EthUtilsInterface {
    const ethUtilsImpl:EthUtilsInterface = {
        createPrivateKey: () => new Promise((res, rej) => {
            osDeps.crypto.randomBytes(32)
                .then((privateKey) => {
                    if (!ethereumjsUtils.isValidPrivate(Buffer.from(privateKey, 'hex'))) {
                        return rej(new errors.InvalidPrivateKeyError());
                    }

                    return res(privateKey);
                })
                .catch(rej);
        }),
        savePrivateKey: (privateKey: string, pw: ?string, pwConfirm: ?string): Promise<void> => new Promise((res, rej) => {
            privateKey = ethUtilsImpl.normalizePrivateKey(privateKey);

            const privateKeyBuffer = Buffer.from(privateKey, 'hex');

            // Reject promise if private key is not a valid hey private key
            if (!ethJsUtils.isValidPrivate(privateKeyBuffer)) {
                return rej(new errors.InvalidPrivateKeyError);
            }

            const addressOfPrivateKey = ethUtilsImpl.normalizeAddress(ethJsUtils.privateToAddress(privateKeyBuffer).toString('hex'));

            const pk:PrivateKeyType = {
                encryption: '',
                value: privateKey,
                encrypted: false,
                version: '1.0.0',
            };

            // Reject promise if one of the passwords is entered AND if they don't match
            if ('undefined' !== typeof pw || 'undefined' !== typeof pwConfirm) {
                if (pw !== pwConfirm) {
                    return rej(new errors.PasswordMismatch);
                }

                // Special chars mach pattern
                const specialCharsPattern = /[ !@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/;

                // $FlowFixMe From a logical point of view the password can't be null / undefined here
                if (specialCharsPattern.test(pw) || specialCharsPattern.test(pwConfirm)) {
                    return rej(new errors.PasswordContainsSpecialChars());
                }

                // Upgrade PK to encrypted one
                pk.value = aes.encrypt(privateKey, pw).toString();
                pk.encryption = 'AES-256';
                pk.encrypted = true;
            }

            // Save the private key
            ss.set(PRIVATE_ETH_KEY_PREFIX+addressOfPrivateKey, JSON.stringify(pk))
                .then(result => {

                    ee.emit(AMOUNT_OF_ADDRESSES_CHANGED);

                    res(result);

                })
                .catch(rej);
        }),
        allKeyPairs: () => new Promise((res, rej) => {
            ss
                .fetchItems((key: string) => key.indexOf(PRIVATE_ETH_KEY_PREFIX) !== -1)
                .then((keys) => {
                    let transformedKeys:Map<string, PrivateKeyType> = new Map();

                    Object.keys(keys).map((key) => {
                        // We only accept string's since. the private key is an stringified object
                        if (typeof keys[key] !== 'string') {
                            return rej(new Error(`Value of key: '${key}' is not an string`));
                        }

                        transformedKeys.set(key.split(PRIVATE_ETH_KEY_PREFIX).pop(), JSON.parse(keys[key]));
                    });

                    res(transformedKeys);
                })
                .catch((err) => rej(err));
        }),
        getPrivateKey: (address: string) => new Promise((res, rej) => {
            const key = PRIVATE_ETH_KEY_PREFIX+ethUtilsImpl.normalizeAddress(address);

            ss
                .has(key)
                .then((hasPrivateKey) => {
                    if (false === hasPrivateKey) {
                        rej(new errors.NoEquivalentPrivateKey());
                        return;
                    }

                    return ss.get(key);
                })
                .then(function(privKey: any) {
                    res(JSON.parse(privKey));
                })
                .catch((err) => rej(err));
        }),
        deletePrivateKey: (address: string) => new Promise((res, rej) => {
            const key = PRIVATE_ETH_KEY_PREFIX+ethUtilsImpl.normalizeAddress(address);

            ss
                .has(key)
                .then((hasPrivateKey) => {
                    if (false === hasPrivateKey) {
                        rej(new errors.NoEquivalentPrivateKey());
                        return;
                    }

                    return ss.remove(key);
                })
                .then((result) => res(result))
                .catch((err) => rej(err));
        }),
        decryptPrivateKey: (privateKey: PrivateKeyType, reason: string, topic: string) => new Promise((mRes, mRej) => {
            // break if the algo is unknown
            if (privateKey.encryption !== 'AES-256') {
                return mRej(new errors.InvalidEncryptionAlgorithm());
            }

            /**
             * @desc Call this to decrypt the password
             * @param {string} pw password
             * @return {Promise<any>}
             */
            function successor(pw: string): Promise<void> {
                return new Promise((res, rej) => {
                    const decryptedPrivateKey = crypto
                        .AES
                        .decrypt(privateKey.value.toString(), pw)
                        .toString(crypto.enc.Utf8);

                    // When aes decryption failes a empty string is returned
                    if ('' === decryptedPrivateKey) {
                        rej(new errors.FailedToDecryptPrivateKeyPasswordInvalid);
                        return;
                    }

                    // Check if decrypted key is valid
                    if (!ethJsUtils.isValidPrivate(Buffer.from(decryptedPrivateKey, 'hex'))) {
                        rej(new errors.DecryptedValueIsNotAPrivateKey());
                        return;
                    }

                    res();
                    mRes(decryptedPrivateKey);
                });
            }

            // Call this to kill the decryption proccess
            const killer = () => {
                mRej(new errors.CanceledAction());
            };

            ee.emit('eth:decrypt-private-key', {
                successor: successor,
                killer: killer,
                reason: reason,
                topic: topic,
            });
        }),
        signTx: (txData: TxData, privKey: string): Promise<EthTx> => new Promise((res, rej) => {
            // Private key as buffer
            const pKB = Buffer.from(privKey, 'hex');

            // reject if private key is invalid
            if (!ethJsUtils.isValidPrivate(pKB)) {
                return rej(new InvalidPrivateKeyError());
            }

            // Sign transaction
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
                abort: () => rej(new AbortedSigningOfTx()),
            });
        }),
        normalizeAddress(address: string): string {
            const checksumAddress:string = ethereumjsUtils.toChecksumAddress(address);

            if (!ethereumjsUtils.isValidChecksumAddress(checksumAddress)) {
                throw new InvalidChecksumAddress(address);
            }

            return checksumAddress;
        },
        normalizePrivateKey(privateKey: string): string {
            if (!ethereumjsUtils.isValidPrivate(Buffer.from(privateKey, 'hex'))) {
                throw new InvalidPrivateKeyError();
            }

            return privateKey;
        },
        privateKeyToMnemonic: (privateKey: string): Array<string> => {
            assert.equal(true, ethJsUtils.isValidPrivate(Buffer.from(privateKey, 'hex')), 'Expected valid private key');

            return bip39.entropyToMnemonic(privateKey).split(' ');
        },
        mnemonicToPrivateKey: (mnemonic: string): string => bip39.mnemonicToEntropy(mnemonic),
        mnemonicValid: bip39.validateMnemonic,
    };

    return ethUtilsImpl;
}
