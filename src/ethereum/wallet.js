//@flow

import {WalletInterface} from "../specification/wallet";
import type {Balance} from "../specification/wallet";
import type {DB} from "../database/db";
import type {EthUtilsInterface} from "./utils";
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
