# Supported calls

A response should always be a json object with a `error` key and a `payload` key.
If there is no error after processing the request let it me an empty string. The payload can be an empty string as well.
IMPORTANT: payload has to be a serialized json object in the case there is a result that should be send.

### Double ratchet key Store

#### DR:KEY_STORE:GET
> Fetch a double ratchet key from the client.

Type: `DR:KEY_STORE:GET`

Data:
- `key` (string)
- `msg_num` (uint)

Response:

```
{
    error: "",
    payload: "{
        key: ""
    }"
}
```

#### DR:KEY_STORE:PUT
> Safe a key to the double ratchet key store

Type: `DR:KEY_STORE:PUT`

Data:
- `index_key` (string) hex string
- `msg_number` (uint) number of the message
- `msg_key` (string) encrypted message key

Response:

```
{
    error: "",
    payload: ""
}
```