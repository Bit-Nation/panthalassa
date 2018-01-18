// @flow
//@todo move to profile
export type PublicProfile = {

    // Pseudo name of the user
    name: string,

    // Description of the user
    description: string,

    // profile image as base64
    image: string,

    // version
    version: string,

    // Ethereum addresses of the user
    ethAddresses: Array<string>

};
