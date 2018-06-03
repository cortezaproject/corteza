The intent of the WebSocket API is to provide a bidirectional communication channel between the browser and the
API services. The listed API calls may be encapsulated in a JSON payload and transmitted over websockets.

A pre-requirement to connect to the websocket endpoint is to provide a valid Session ID, which is returned over
the `/member/login` endpoint. After connecting to the endpoint, the user may issue chat-specific events.

These are custom messages which can be sent over WebSockets APIs:

- `join` - join a channel
- `part` - part a channel
- `message` - send message to channel or user
- `typing` - send indicator of typing
- `presence` - send presence status (do not disturb, idle/away, available)

These messages are automatically sent by the server upon connect:

- `channel` - list of channels,
  - `members` - list of members in channel,
  - `messages` - a predefined backlog for a channel (more messages are retrieved by issuing `/message/search` API calls)
- `message` - individual new messages
