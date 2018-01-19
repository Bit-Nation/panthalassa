// @flow

import {DBInterface} from '../database/db';
import {EthUtilsInterface} from './utils';
import type {AccountBalanceType} from '../database/schemata';
const Web3 = require('web3');

/**
 * @typedef {Object} WalletInterface
 * @property {function(from:string, to:string, amount:string)} ethSend send ether
 * @property {function} ethBalance fetch the balance
 * @property {function} ethSync sync the ethereum accounts
 */
export interface WalletInterface {

    /**
     * @desc Send eth from A to B
     * @param {string} the address of the account where to send from
     * @param {string} address where to send
     * @param {number} amount of eth you would like to send (in eth)
     * @param {number} the limit of gas for this transaction (default is 21000
     * @param {number} gas price default is 20000000000
     * @return {Promise<string>} resolves in transaction hash
     */
    ethSend: (from: string, to: string, amount: number, gasLimit: number, gasPrice: number) => Promise<string>,

    /**
     * @desc Fetch the eth balance of one of my owned accounts
     * @param {string} address of the account you would like to sync
     * @return {Promise<AccountBalanceType>}
     */
    ethBalance: (address: string) => Promise<AccountBalanceType | null>,

    /**
     * @desc Sync balance of specific ethereum address
     * @param {string} address of the account you would like to update
     * @return {Promise} promise that will resolve with void
     */
    ethSync: (address: string) => Promise<void>,
}

/**
 * @module ethereum/wallet.js
 * @desc Factory that returns an wallet object
 * @param {EthUtilsInterface} ethUtils
 * @param {object} web3
 * @param {DBInterface} db
 * @return {WalletInterface}
 */
export default function walletFactory(ethUtils: EthUtilsInterface, web3: Web3, db: DBInterface): WalletInterface {
    const walletImpl:WalletInterface = {
        ethSend: (from: string, to: string, amount: number, gasLimit: number, gasPrice: number): Promise<string> => {
            gasLimit = gasLimit || 21000;
            gasPrice = gasPrice || 20000000000;

            return new Promise((res, rej) => {
                // Will throw error if invalid address so we need to catch it
                try {
                    from = ethUtils.normalizeAddress(from);

                    to = ethUtils.normalizeAddress(to);
                } catch (e) {
                    return rej(e);
                }

                web3.eth.sendTransaction({
                    from: from,
                    to: to,
                    value: web3.toWei(amount, 'ether'),
                    gasLimit: gasLimit,
                    gasPrice: gasPrice,
                }, (error, txReceipt) => {
                    if (error) {
                        return rej(error);
                    }

                    res(txReceipt);
                });
            });
        },
        ethBalance: (address: string): Promise<AccountBalanceType | null> => {
            return new Promise((res, rej) => {
                try {
                    address = ethUtils.normalizeAddress(address);
                } catch (e) {
                    return rej(e);
                }

                db
                    .query((realm) => {
                        const balances = realm.objects('AccountBalance').filtered(`id == "${address}_ETH"`);

                        if (balances.length <= 0) {
                            return res(null);
                        }

                        if (balances.length === 1) {
                            return res(balances[0]);
                        }

                        rej(`Expected balances.length to be '<=1'. Got: ${balances.length}`);
                    })
                    .then(res)
                    .catch(rej);
            });
        },
        ethSync: (address: string): Promise<void> => new Promise((res, rej) => {
            try {
                address = ethUtils.normalizeAddress(address);
            } catch (e) {
                return rej(e);
            }

            web3.eth.getBalance(address, (error, balance) => {
                if (error) {
                    return rej(error);
                }

                //Transform balance to string (will be in wei)
                balance = balance.toString(10);

                if ('string' !== typeof balance) {
                    return rej(new Error('Fetched balance is not a string'));
                }

                db
                    .write((realm) => {
                        realm.create('AccountBalance', {
                            id: address+'_ETH',
                            address: address,
                            currency: 'ETH',
                            synced_at: new Date(),
                            amount: web3.fromWei(balance, 'ether'),
                        }, true);
                    })
                    .then(_ => res())
                    .catch(rej);
            });
        }),

    };

    return walletImpl;
}
