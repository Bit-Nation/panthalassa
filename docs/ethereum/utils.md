# Utils

The `ethereum/utils.js` file contains a many method to do ethereum related things. It e.g. used for
- generating a private key
- saving a private key
- fetching all the available key pairs
- to get an specific private key by its address
- to delete a private key by it's address
- to decrypt an private key
- signing transactions (more on that below)
- normalizing addresses and private keys
- transforming an private key (hex format) to an mnemonic and vise versa
- validating mnemonic's

#### singing transaction
Signing an transaction is a central functionality in our application.

panga libs                   user

    Send singing Request
------------------------------->

    User approve / abort
<------------------------------

In order to for this workflow to work, you need to listen for the `eth:tx:sign` event. Have a look at the code example in the interface.