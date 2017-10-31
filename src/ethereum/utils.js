//@flow

const crypto = require('crypto');
const Promise = require('promise');
const ethereumjsUtils = require('ethereumjs-util');
const errors = require('./../errors');
const aes = require('crypto-js/aes');

const PRIVATE_ETH_KEY_PREFIX = 'PRIVATE_ETH_KEY#';

/**
 * Creates a new private key
 * @param crypto
 * @param isValidPrivateKey
 * @returns {function()}
 * Todo change the promise * to real typehint
 */
const createPrivateKey = (crypto, isValidPrivateKey: (key: Buffer) => boolean) : (() => Promise<*>) => {
    "use strict";

    //Todo change the promise * to real typehint
    return () : Promise<*> => {

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

};

/**
 *
 * @param secureStorage
 * @param ethjsUtils
 * @param aes
 * @returns {function(string, string, string)}
 */
const savePrivateKey = (secureStorage: any, ethjsUtils: ethereumjsUtils, aes: any) : ((privateKey:string, pw:?string, pwConfirm:?string) => Promise<*>)  => {
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
                }

                //Save the private key
                secureStorage.set(
                    PRIVATE_ETH_KEY_PREFIX+addressOfPrivateKey,
                    JSON.stringify({
                        encryption: 'AES-256',
                        value: aes.encrypt(privateKey, pw).toString(),
                        encrypted: true,
                        version: '1.0.0'
                    })
                )
                    .then(result => res(result))
                    .catch(err => rej(err));

                return;
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

};

module.exports = (secureStorage:any) : {} => {
    "use strict";

    return {
        createPrivateKey: createPrivateKey(crypto, ethereumjsUtils.isValidPrivate),
        savePrivateKey: savePrivateKey(secureStorage, ethereumjsUtils, aes),
        raw: {
            createPrivateKey: createPrivateKey,
            savePrivateKey: savePrivateKey
        }
    }

};
