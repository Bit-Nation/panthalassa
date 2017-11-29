//@flow

///////////////////////////////////////////////////////////
// ATTENTION !!! Everyime you update the schema,         //
//               update the relating interfaces as well. //
///////////////////////////////////////////////////////////

//Profile

/**
 * A note on this ProfileObject interface.
 * In the project you will often see smth like this:
 *
 * (this is an example from the queries)
 * findProfiles(realm) : Array<ProfileObject>
 *
 * The value returned by realm is NOT directly a instance of an object that implement this interface,
 * BUT the signature is exactly the same.
 *
 * It's ok to do this, since after the compilation from flow -> js all interfaces
 * and types are striped and they are all objects. So this interface is here to
 * support the developers.
 */
export interface ProfileObject {
    id: number,
    pseudo: string,
    description: string,
    image: string,
    version: string
}

export const ProfileSchema = {
    name: 'Profile',
    primaryKey: 'id',
    properties: {
        id: 'int',
        pseudo: 'string',
        description: 'string',
        image: 'string',
        version: 'string'
    },
};

/**
 * AccountBalance
 */
export type AccountBalanceType = {
    id: string,
    address:string,
    //Amount is in wei
    amount:string,
    synced_at:number,
    currency:string
}

export const AccountBalanceSchema = {
    name: 'AccountBalance',
    primaryKey: 'id',
    properties: {
        id: 'string',
        address: 'string',
        amount: 'string',
        synced_at: 'date',
        currency: 'string'
    }
};