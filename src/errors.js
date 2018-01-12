// @flow

/**
 * @ignore
 */
export class InvalidPrivateKeyError extends Error {}

/**
 * @ignore
 */
export class PasswordMismatch extends Error {}

/**
 * @ignore
 */
export class PasswordContainsSpecialChars extends Error {}

/**
 * @ignore
 */
export class NoEquivalentPrivateKey extends Error {}

/**
 * @ignore
 */
export class InvalidEncryptionAlgorithm extends Error {}

/**
 * @ignore
 */
export class FailedToDecryptPrivateKeyPasswordInvalid extends Error {}

/**
 * @ignore
 */
export class CanceledAction extends Error {}

/**
 * @ignore
 */
export class DecryptedValueIsNotAPrivateKey extends Error {}

/**
 * @ignore
 */
export class NoProfilePresent extends Error {}

/**
 * @ignore
 */
export class NoPublicProfilePresent extends Error {}

/**
 * @ignore
 */
export class AbortedSigningOfTx extends Error {}

/**
 * @ignore
 */
export class InvalidChecksumAddress extends Error {
    /**
     *
     * @param {string} address
     */
    constructor(address: string) {
        super('Address: '+address+' is invalid');
    }
}
