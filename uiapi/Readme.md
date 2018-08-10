# Events

- Messages

    - `MESSAGE:PERSISTED`
        - `db_id` (database id of message as string)
        - `partner` (hex encoded ed25519 public key)
        - `content` raw message content (will be "" if DApp message)
        - `created_at` unix timestamp

    - `MESSAGE:DELIVERED`
        - `db_id` (database id of message as string)
        - `partner` (hex encoded ed25519 public key)
        - `content` raw message content (will be "" if DApp message)
        - `created_at` unix timestamp


    - `MESSAGE:RECEIVED`
        - `db_id` (database id of message as string)
        - `partner` (hex encoded ed25519 public key)
        - `content` raw message content (will be "" if DApp message)
        - `created_at` unix timestamp

- DApp

    - `DAPP:PERSISTED`
        - `dapp_signing_key` hex encoded signing key used to sign the DApp