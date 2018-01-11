// @flow

import {DB} from '../database/db';
import {EthUtilsInterface} from './utils';
import type {AccountBalanceType} from '../database/schemata';
const Web3 = require('web3');

export interface WalletInterface {

    /**
     * from is the senders ethereum address
     * to is the receiver ethereum address
     * amount in ether NOT in wei
     */
    ethSend: (from: string, to: string, amount: number, gasLimit: number, gasPrice: number) => Promise<{...mixed}>,

    /**
     * Get balance of account.
     * Will resolve in object or null
     */
    ethBalance: (address: string) => Promise<AccountBalanceType | null>,

    /**
     * Sync balance of specific ethereum address
     */
    ethSync: (address: string) => Promise<void>,
}

/**
 *
 * @param {object} ethUtils
 * @param {object} web3
 * @param {object} db
 * @return {WalletInterface}
 */
export default function(ethUtils: EthUtilsInterface, web3: Web3, db: DB): WalletInterface {
    const walletImpl:WalletInterface = {
        ethSend: (from: string, to: string, amount: number, gasLimit: number = 21000, gasPrice: number = 20000000000): Promise<{...mixed}> => {
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
                    ethUtils.normalizeAddress(address);
                } catch (e) {
                    rej(e);
                }

                db.query((realm) => {
                    const balances = realm.objects('AccountBalance').filtered(`id == "${address}_ETH"`);

                    if (balances.length <= 0) {
                        return res(null);
                    }

                    if (balances.length === 1) {
                        return res(balances[0]);
                    }

                    rej(`Expected balances.length to be '<=1'. Got: ${balances.length}`);
                });
            });
        },
        ethSync: (address: string): Promise<void> => new Promise((res, rej) => {
            try {
                ethUtils.normalizeAddress(address);
            } catch (e) {
                return rej(e);
            }

            web3.eth.getBalance(address, (error, balance) => {
                if (error) {
                    return rej(error);
                }

                if ('string' !== typeof balance) {
                    return rej(new Error('Fetched balance is not a string'));
                }

                db.write((realm) => {
                    realm.create('AccountBalance', {
                        id: address+'_ETH',
                        address: address,
                        currency: 'ETH',
                        synced_at: Date.now(),
                        amount: balance,
                    }, true);

                    res();
                });
            });
        }),
    };

    return walletImpl;
}
