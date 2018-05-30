# Supported calls

#### PRE_KEY_PUT

Type: `PRE_KEY_PUT`
> Save a public key with a corresponding private key

Data:
- `public_key` (string)
- `private_key` (string)

#### PRE_KEY_FETCH
> Fetch a pre key based on the public key

Type: `PRE_KEY_FETCH`

Data:
- `public_key` (string)