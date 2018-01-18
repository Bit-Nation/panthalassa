# DB

The database (`database/db.js`) is an proxy between realm and the application.
You can call `write` and `query` on every DBInterface object.
Those methods expect an function dthat should be executed in the `query` or `write` context.