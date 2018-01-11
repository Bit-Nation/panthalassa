// @flow
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
