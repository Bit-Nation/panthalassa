// @flow

class InvalidPrivateKeyError extends Error{}

class PasswordMismatch extends Error{}

class PasswordContainsSpecialChars extends Error{}

// Is thrown when there is no private key for an address
class NoEquivalentPrivateKey extends Error{}

class InvalidEncryptionAlgorithm extends Error{}

class FailedToDecryptPrivateKeyPasswordInvalid extends Error{}

class CanceledAction extends Error{}

class DecryptedValueIsNotAPrivateKey extends Error{}

class NoProfilePresent extends Error{}

class NoPublicProfilePresent extends Error{}

/**
 * Is thrown if a method on the secure storage implementation is missing
 */
class UnsatisfiedSecureStorageImplementationError extends Error{

    constructor(missingMethodName: string) {

        super('Missing method: "'+missingMethodName+'" in secure storage implementation');

    }

}

module.exports = {
    InvalidPrivateKeyError,
    UnsatisfiedSecureStorageImplementationError,
    PasswordMismatch,
    PasswordContainsSpecialChars,
    NoEquivalentPrivateKey,
    InvalidEncryptionAlgorithm,
    FailedToDecryptPrivateKeyPasswordInvalid,
    CanceledAction,
    DecryptedValueIsNotAPrivateKey,
    NoProfilePresent,
    NoPublicProfilePresent
};
