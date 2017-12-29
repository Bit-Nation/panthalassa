//@flow

export interface Crypto {

    randomBytes: (length:number) => Promise<string>

}

/**
 * @description Used to abstract some of the dep's that are only available on e.g. node. Done to be OS agnostic.
 */
export interface OsDependenciesInterface {

    crypto:Crypto

}