# Events

- Messages

    - `MESSAGE:PERSISTED`
        - `message_id`
        - `partner` (hex encoded ed25519 public key)
        - `message` plain message protobuf encoded as string

    - `MESSAGE:DELIVERED`
        - `message_id`
        - `partner` (hex encoded ed25519 public key)
        - `message` plain message protobuf encoded as string


    - `MESSAGE:RECEIVED`
        - `message_id`
        - `partner` (hex encoded ed25519 public key)
        - `message` plain message protobuf encoded as string
