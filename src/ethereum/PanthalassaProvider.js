//@flow

import type {EthUtilsInterface} from "./utils";
import type {PrivateKeyType} from "../specification/privateKey";
import type {TxData} from "../specification/tx";
const EthTx = require('ethereumjs-tx');

const ZeroProvider = require('web3-provider-engine/zero');

/**
 * fetch all accounts
 * @param ethUtils
 * @returns {function(*)}
 */
export function getAccounts(ethUtils:EthUtilsInterface) : (cb: (error:any, addresses:any) => void) => void{

    return (cb:(error:any, addresses:any) => void) : void => {

        ethUtils.allKeyPairs()
            .then(keyPairs => {

                const addresses:Array<string> = [];

                keyPairs.map((keyPair) => addresses.push(keyPair.key));

                cb(null, addresses);

            })
            .catch(error => {

                cb(error, null)

            })

    }

}

/**
 * sign's transaction
 * @param ethUtils
 * @returns {function(TxData, *)}
 */
export function signTx(ethUtils:EthUtilsInterface) : (txData:TxData, cb:(error:any, signedTx:any) => void) => void {

    return (txData:TxData, cb:(error:any, signedTx:any) => void) : void => {

        ethUtils
            .getPrivateKey(txData.from)
            .then(async function(privateKey:PrivateKeyType){

                try{

                    let pk:string = privateKey.value;

                    if(privateKey.encrypted){
                        pk = await ethUtils.decryptPrivateKey(privateKey, 'Please decrypt your private key in order to sign the transaction', 'Sign transaction')
                    }

                    ethUtils.signTx(txData, pk)

                        .then(function(signedTx:EthTx){
                            cb(null, '0x'+signedTx.serialize().toString('hex'));
                        })

                        .catch(e => cb(e, null));

                }catch(e) {
                    cb(e, null);
                }

            })
            .catch(e => cb(e, null))

    }

}

export default class PanthalassaProvider extends ZeroProvider {

    constructor(ethUtils:EthUtilsInterface, rpcUrl: string){

        super({
            getAccounts: getAccounts(ethUtils),
            signTransaction: signTx(ethUtils),
            rpcUrl: rpcUrl
        })

    }

}
