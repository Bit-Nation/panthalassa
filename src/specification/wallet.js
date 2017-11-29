//@flow

export type Balance = {
    account: string,
    balance: string,
    timestamp: number
}

export interface WalletInterface {

    //Resolves with success value from web.eth.sendTransaction()
    ethSend: (from:string, to:string, amount:string, gasLimit:number, gasPrice:number) => Promise<{...mixed}>,

    ethBalance: (address:string) => Promise<Balance>,

    ethSync: (address:string) => Promise<void>,

    syncCurrencies: (address:string) => Promise<[typeof undefined]>

}