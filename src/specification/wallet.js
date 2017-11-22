//@flow

export type Balance = {
    account: string,
    balance: string,
    timestamp: number
}

export interface WalletInterface {

    ethSend: (from:string, to:string) => Promise<void>,

    ethBalance: (address:string) => Promise<Balance>,

    ethSync: (address:string) => Promise<void>,

    patSend: (from:string, to:string) => Promise<void>,

    patBalance: (address:string) => Promise<Balance>,

    patSync: (address:string) => Promise<void>,

    syncCurrencies: (address:string) => Promise<[typeof undefined]>

}