// @flow

// /////////////////////////////////////////////////////////
// ATTENTION !!! Everyime you update the schema,         //
//               update the relating interfaces as well. //
// /////////////////////////////////////////////////////////

// A note on e.g. ProfileType (and all types exported from this file).
// In the project you will often see smth like this:

// (this is an example from the queries)
// findProfiles(realm) : Array<ProfileObject>

// The value returned by realm is NOT directly a instance of an object that implement this interface,
// BUT the signature is exactly the same.

// It's ok to do this, since after the compilation from flow -> js all interfaces
// and types are striped and they are all objects. So this interface is here to
// support the developers.

/**
 * @typedef ProfileType
 * @property {number} id
 * @property {string} pseudo
 * @property {string} description
 * @property {string} image
 * @property {string} version
 */
export type ProfileType = {
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
        version: 'string',
    },
};

/**
 * @typedef AccountBalanceType
 * @property {string} id
 * @property {string} address
 * @property {string} amount
 * @property {number} synced_at
 * @property {string} currency
 */
export type AccountBalanceType = {
    id: string,
    address: string,
    // Amount is in wei
    amount: string,
    synced_at: number,
    currency: string
}

export const AccountBalanceSchema = {
    name: 'AccountBalance',
    primaryKey: 'id',
    properties: {
        id: 'string',
        address: 'string',
        amount: 'string',
        synced_at: 'date',
        currency: 'string',
    },
};

/**
 * @typedef MessageJobType
 * @property {number} id
 * @property {string} heading
 * @property {string} text
 * @property {number} version
 * @property {date} created_at
 */
export type MessageJobType = {
    id: number,
    heading: string,
    text: string,
    version: number,
    created_at: number
}

export const MessageJobSchema = {
    name: 'MessageJob',
    primaryKey: 'id',
    properties: {
        id: 'int',
        heading: 'string',
        text: 'string',
        version: 'date',
        created_at: 'date'
    }
};