Messages represent individual messages in the chat system. Messages are typed, indicating the event which triggered the message.

Currently expected message types are:

| Name | Description |
| ---- | ----------- |
| CREATE | The first message when the channel is created |
| TOPIC | A member changed the topic of the channel |
| RENAME | A member renamed the channel |
| MESSAGE | A member posted a message to the channel |
| FILE | A member uploaded a file to the channel |

The following event types may be sent with a message event:

| Name | Description |
| ---- | ----------- |
| CREATED | A message has been created on a channel |
| EDITED | A message has been edited by the sender |
| REMOVED | A message has been removed by the sender |
