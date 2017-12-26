## Utils (EthUtils)

### How to use?
Import the default function from `src/ethereum/utils.js` and call it with an object that satisfies the [Secure Storage Interface](./../specification/secureStorageInterface.js) as the first parameter and with an instance of [Event Emitter 3](https://www.npmjs.com/package/eventemitter3) as the second parameter.

### API
Have a look at the [EthUtilsInterface](./utils.js).

__Some of the more complex method's__

#### `decryptPrivateKey`

The `decryptPrivateKey` method is responsible for decrypting a private key if it was encrypted. Since we don't know the password we have to ask the user for the password. `decryptPrivateKey` will emit an event called `eth:decrypt-private-key` with the following data:
```js
{
    successor,
    killer,
    reason,
    topic
}
```

- `topic` is could be something like 'Decrypt Private Key'

- `reason` could be something like 'Send money from 0x... to 0x...'

- `successor` is a function that takes a password (string) as it's only parameter. If the password is valid the returned promise will be resolved if the password is invalid the promise will be rejected.

- `killer` is a function that is used to abort the decryption process.

The idea of this event is, that the backend can inform the client about the need to decrypt the private key. The client can then call the successor / killer to interact with the backend. 

#### `signTx`

The `signTx` method work's like the `decryptPrivateKey` method expect for that this method is responsible for signing an transaction. When called it will emit an event called `eth:tx:sign` with the following data:
```flow
{
    tx,
    txData,
    confirm,
    abort
}
```

- `tx` is an instance of [EthTx](https://www.npmjs.com/package/ethereumjs-tx)
- `txData` is an object that look's like that:
```js
{
  nonce: '0x00',
  gasPrice: '0x09184e72a000', 
  gasLimit: '0x2710',
  to: '0x0000000000000000000000000000000000000000', 
  value: '0x00', 
  data: '0x7f7465737432000000000000000000000000000000000000000000000000000000600057',
  // EIP 155 chainId - mainnet: 1, ropsten: 3 
  chainId: 3
}
```

- `confirm` is an function. Call it to confirm the signing. 
- `abort` abort signing. Call it to abort the signing. 

Keep in mind to listen to the `eth:tx:sign` event. If you don't call `confirm` or `abort` the returned promise will never be resolved.

## Wallet

### How to use?
Import the exported default function from `src/ethereum/wallet.js` and call it with an object that satisfies the EthUtilsInterface, an [Web3 object](web3.js) and an [DB object](../database/db.js)

### API
Have a look at the [WalletInterface](wallet.js).

## Web3

### How to use?
Import the exported default function from `src/ethereum/web3` and call it with object that satisfies the [JsonRpcNodeInterface](../specification/jsonRpcNode.js), an instance of [EventEmitter](https://www.npmjs.com/package/eventemitter3) and an object that satisfies the [EthUtilsInterface](utils.js) interface. 

### Api
- 

## PanthalassaProvider
> The PanthalassaProvider is an web3 provider that contain's custom logic, e.g. for signing an transaction.

### How to use?
Import the exported default class and instantiate it with an object that satisfies the [EthUtilsInterface](utils.js) and with an url (string) to the json rpc.

### Api
- 