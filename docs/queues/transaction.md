# Transaction Queue

There are some web3 interactions (like the nation creation) which can't be done in one web3 request. We just save the nation data in a TransactionJob and submit the first part of the data to the blockchain, wait for the result and depending on the result we submit the second part and so on.