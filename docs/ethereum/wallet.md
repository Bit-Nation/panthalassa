# Pangea wallet

In order to make use of ethereum (e.g. create a Nation) you need ETH. In order to get it you need wallet functionality. This functionality lives in `ethereum/wallet.js`. On a low level the wallet makes requests to the JSON rpc via Web3. The wallet don't check if you have an active internet connection, so make sure to check it your self when using the wallet.