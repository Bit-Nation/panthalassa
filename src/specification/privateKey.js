// @flow
// @todo move to ethereum/utils.js
/**
 * @typedef PrivateKeyType
 * @property {string} encryption should be an string of the encryption algo used
 * @property {string} value
 * @property {boolean} encrypted
 * @property {string} version
 */
export type PrivateKeyType = {

    // used encryption algo as a string
    encryption: string,

    // private key (encrypted or plain)
    value: string,

    // true if so false if not
    encrypted: boolean,

    // version
    version: string

}
