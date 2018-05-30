# Device Api
> The device api is responsible for communication between the client and panthalassa

## Supported Call's
You can finda list of suported calls [here](./rpc)

## Implementing the Upstream
The upstream is responsible for sending data from panthalassa to the client.
A client is e.g. [Pangea](https://github.com/Bit-Nation/BITNATION-Pangea-mobile) our mobile interface.
Panthalassa will call the `send` function of the upstream implementation with rpc call's.
DON'T forget to send a response back to panthalassa.
You can send a response by calling the `SendResponse` on panthalassa with the `id` from the rpc call and the requires response data.

## Using the API
Every call will be structured like this:
```
{
    type: string,
    id:   string,
    data: string
}
```

__type__
Is the `type` of the call. E.g. `DHT_PUT` which tell you to put a value in the DHT.

__id__
The id of the call. You HAVE TO call send your response to the call with this id.

__data__
The data contains the payload. For a call of type `DHT_PUT` it data would be: `{key: "0xf...", value: "base64..."}`