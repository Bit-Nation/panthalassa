# Specifications

Some Interfaces / Types don't really fit into a specific folder / fil (like the SecureStorageInterface, which is used to access the keychain and must be implemented for each env).

### JsonRpcNode
Basically an interface for the nodes we are using.

### OsDependencies
A few functions we need (like randomBytes) need to be implemented for each env. `specification/osDependencies.js` provide and interface for that.

### PrivateKey
And type for private key that holds some meta data about it

### PublicProfile
A type for the public profile

### Tx
Type for the transaction data