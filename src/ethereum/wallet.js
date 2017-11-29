//@flow

import {WalletInterface} from "../specification/wallet";
import {DB} from "../database/db";
import {EthUtilsInterface} from "./utils";
import type {AccountBalanceType} from '../database/schemata';
const Web3 = require('web3');
const ethereumJsUtils = require('ethereumjs-util');

export function ethSend(ethUtils:EthUtilsInterface, web3:Web3) {

    return (from:string, to:string, amount:number, gasLimit:number = 21000, gasPrice:number = 20000000000) : Promise<{...mixed}> => {

        return new Promise((res, rej) => {

            //Will throw error if invalid address so we need to catch it
            try{
                from = ethUtils.normalizeAddress(from);

                to = ethUtils.normalizeAddress(to);
            }catch (e){
                return rej(e);
            }

            web3.eth.sendTransaction({
                from: from,
                to: to,
                value: web3.utils.toWei(amount, 'ether'),
                gasLimit: gasLimit,
                gasPrice: gasPrice
            }, (error, txReceipt) => {

                if(error){
                    return rej(error);
                }

                res(txReceipt);

            });

        });

    }

}

export function ethBalance(db:DB, ethUtils:EthUtilsInterface) {

    return (address:string) : Promise<AccountBalanceType | null> => {

        return new Promise((res, rej) => {

            try {
                ethUtils.normalizeAddress(address);
            }catch (e){
                rej(e);
            }

            db.query((realm) => {

                const balances = realm.objects('AccountBalance').filtered(`id == "${address}_ETH"`);

                if (balances.length <= 0){
                    return res(null);
                }

                if (balances.length === 1) {
                    return res(balances[0]);
                }

                rej(`Expected balances.length to be '<=1'. Got: ${balances.length}`);

            })

        });

    }

}
