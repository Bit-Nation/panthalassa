// @flow

/**
 *
 */
export class InvalidPrivateKeyError extends Error {}

/**
 *
 */
export class PasswordMismatch extends Error {}

/**
 *
 */
export class PasswordContainsSpecialChars extends Error {}

/**
 *
 */
export class NoEquivalentPrivateKey extends Error {}

/**
 *
 */
export class InvalidEncryptionAlgorithm extends Error {}

/**
 *
 */
export class FailedToDecryptPrivateKeyPasswordInvalid extends Error {}

/**
 *
 */
export class CanceledAction extends Error {}

/**
 *
 */
export class DecryptedValueIsNotAPrivateKey extends Error {}

/**
 *
 */
export class NoProfilePresent extends Error {}

/**
 *
 */
export class NoPublicProfilePresent extends Error {}

/**
 *
 */
export class AbortedSigningOfTx extends Error {}

/**
 *
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
