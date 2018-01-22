/**
 * @desc should be dispatched when the client goes online
 * @type {string}
 */
export const APP_ONLINE = 'app:online';

/**
 * @desc should be dispatched when the client goes offline
 * @type {string}
 */
export const APP_OFFLINE = 'app:offline';

/**
 * @desc Should be dispatched when the amount of addresses changes
 * @type {string}
 */
export const AMOUNT_OF_ADDRESSES_CHANGED = 'amount_of_addresses_changed';

/**
 * @desc Should be dispatched when a job was added to the messaging queue
 * @type {string}
 */
export const MESSAGING_QUEUE_JOB_ADDED = 'messaging_queue:job:added';

/**
 * @desc Is emitted when an job as added to the transaction queue
 * @type {string}
 */
export const TRANSACTION_QUEUE_JOB_ADDED = 'transaction_queue:job:added';

/**
 * @desc Emitted when there is an transaction to sign
 * @type {string}
 */
export const ETH_TX_SIGN = 'eth:tx:sign';

/**
 * @desc Is emitted when a nation is created
 * @type {string}
 */
export const NATION_CREATE = 'nation:create';
