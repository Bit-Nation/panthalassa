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
Let the key be an empty string if there is no key found for the request.

```
{
    error: "",
    payload: "{
        key: ""
    }"
}
```
