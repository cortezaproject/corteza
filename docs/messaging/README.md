# User activity

| Method | Endpoint | Purpose |
| ------ | -------- | ------- |
| `POST` | `/activity/` | Sends user's activity to all subscribers; globally or per channel/message. |

## Sends user's activity to all subscribers; globally or per channel/message.

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/activity/` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| channelID | uint64 | POST | Channel ID, if set, activity will be send only to subscribed users | N/A | NO |
| messageID | uint64 | POST | Message ID, if set, channelID must be set as well | N/A | NO |
| kind | string | POST | Arbitrary string | N/A | YES |

---




# Attachments

| Method | Endpoint | Purpose |
| ------ | -------- | ------- |
| `GET` | `/attachment/{attachmentID}/original/{name}` | Serves attached file |
| `GET` | `/attachment/{attachmentID}/preview.{ext}` | Serves preview of an attached file |

## Serves attached file

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/attachment/{attachmentID}/original/{name}` | HTTP/S | GET | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| download | bool | GET | Force file download | N/A | NO |
| sign | string | GET | Signature | N/A | YES |
| userID | uint64 | GET | User ID | N/A | YES |
| name | string | PATH | File name | N/A | YES |
| attachmentID | uint64 | PATH | Attachment ID | N/A | YES |

## Serves preview of an attached file

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/attachment/{attachmentID}/preview.{ext}` | HTTP/S | GET | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| ext | string | PATH | Preview extension/format | N/A | YES |
| attachmentID | uint64 | PATH | Attachment ID | N/A | YES |
| sign | string | GET | Signature | N/A | YES |
| userID | uint64 | GET | User ID | N/A | YES |

---




# Channels

A channel is a representation of a sequence of messages. It has meta data like channel subject. Channels may be public, private or group.

| Method | Endpoint | Purpose |
| ------ | -------- | ------- |
| `GET` | `/channels/` | List channels |
| `POST` | `/channels/` | Create new channel |
| `PUT` | `/channels/{channelID}` | Update channel details |
| `PUT` | `/channels/{channelID}/state` | Update channel state |
| `PUT` | `/channels/{channelID}/flag` | Update channel membership flag |
| `DELETE` | `/channels/{channelID}/flag` | Remove channel membership flag |
| `GET` | `/channels/{channelID}` | Read channel details |
| `GET` | `/channels/{channelID}/members` | List channel members |
| `PUT` | `/channels/{channelID}/members/{userID}` | Join channel |
| `DELETE` | `/channels/{channelID}/members/{userID}` | Remove member from channel |
| `POST` | `/channels/{channelID}/invite` | Join channel |
| `POST` | `/channels/{channelID}/attach` | Attach file to channel |

## List channels

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/channels/` | HTTP/S | GET | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| query | string | GET | Search query | N/A | NO |

## Create new channel

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/channels/` | HTTP/S | POST | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| name | string | POST | Name of Channel | N/A | YES |
| topic | string | POST | Subject of Channel | N/A | NO |
| type | string | POST | Channel type | N/A | NO |
| members | []string | POST | Initial members of the channel | N/A | NO |

## Update channel details

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/channels/{channelID}` | HTTP/S | PUT | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| channelID | uint64 | PATH | Channel ID | N/A | YES |
| name | string | POST | Name of Channel | N/A | NO |
| topic | string | POST | Subject of Channel | N/A | NO |
| type | string | POST | Channel type | N/A | NO |
| organisationID | uint64 | POST | Move channel to different organisation | N/A | NO |

## Update channel state

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/channels/{channelID}/state` | HTTP/S | PUT | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| channelID | uint64 | PATH | Channel ID | N/A | YES |
| state | string | POST | Valid values: delete, undelete, archive, unarchive | N/A | YES |

## Update channel membership flag

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/channels/{channelID}/flag` | HTTP/S | PUT | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| channelID | uint64 | PATH | Channel ID | N/A | YES |
| flag | string | POST | Valid values: pinned, hidden, ignored | N/A | YES |

## Remove channel membership flag

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/channels/{channelID}/flag` | HTTP/S | DELETE | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| channelID | uint64 | PATH | Channel ID | N/A | YES |

## Read channel details

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/channels/{channelID}` | HTTP/S | GET | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| channelID | uint64 | PATH | Channel ID | N/A | YES |

## List channel members

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/channels/{channelID}/members` | HTTP/S | GET | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| channelID | uint64 | PATH | Channel ID | N/A | YES |

## Join channel

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/channels/{channelID}/members/{userID}` | HTTP/S | PUT | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| channelID | uint64 | PATH | Channel ID | N/A | YES |
| userID | uint64 | PATH | Member ID | N/A | NO |

## Remove member from channel

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/channels/{channelID}/members/{userID}` | HTTP/S | DELETE | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| channelID | uint64 | PATH | Channel ID | N/A | YES |
| userID | uint64 | PATH | Member ID | N/A | NO |

## Join channel

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/channels/{channelID}/invite` | HTTP/S | POST | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| channelID | uint64 | PATH | Channel ID | N/A | YES |
| userID | []uint64 | POST | User ID | N/A | NO |

## Attach file to channel

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/channels/{channelID}/attach` | HTTP/S | POST | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| channelID | uint64 | PATH | Channel ID | N/A | YES |
| replyTo | uint64 | POST | Upload as a reply | N/A | NO |
| upload | *multipart.FileHeader | POST | File to upload | N/A | YES |

---




# Commands

| Method | Endpoint | Purpose |
| ------ | -------- | ------- |
| `GET` | `/commands/` | List of available commands |

## List of available commands

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/commands/` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |

---




# Messages

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

| Method | Endpoint | Purpose |
| ------ | -------- | ------- |
| `POST` | `/channels/{channelID}/messages/` | Post new message to the channel |
| `POST` | `/channels/{channelID}/messages/command/{command}/exec` | Execute command |
| `GET` | `/channels/{channelID}/messages/mark-as-read` | Manages read/unread messages in a channel or a thread |
| `PUT` | `/channels/{channelID}/messages/{messageID}` | Edit existing message |
| `DELETE` | `/channels/{channelID}/messages/{messageID}` | Delete existing message |
| `POST` | `/channels/{channelID}/messages/{messageID}/replies` | Reply to a message |
| `POST` | `/channels/{channelID}/messages/{messageID}/pin` | Pin message to channel (public bookmark) |
| `DELETE` | `/channels/{channelID}/messages/{messageID}/pin` | Pin message to channel (public bookmark) |
| `POST` | `/channels/{channelID}/messages/{messageID}/bookmark` | Bookmark a message (private bookmark) |
| `DELETE` | `/channels/{channelID}/messages/{messageID}/bookmark` | Remove boomark from message (private bookmark) |
| `POST` | `/channels/{channelID}/messages/{messageID}/reaction/{reaction}` | React to a message |
| `DELETE` | `/channels/{channelID}/messages/{messageID}/reaction/{reaction}` | Delete reaction from a message |

## Post new message to the channel

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/channels/{channelID}/messages/` | HTTP/S | POST | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| message | string | POST | Message contents (markdown) | N/A | YES |
| channelID | uint64 | PATH | Channel ID | N/A | YES |

## Execute command

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/channels/{channelID}/messages/command/{command}/exec` | HTTP/S | POST | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| command | string | PATH | Command to be executed | N/A | YES |
| channelID | uint64 | PATH | Channel ID | N/A | YES |
| input | string | POST | Arbitrary command input | N/A | NO |
| params | []string | POST | Command parameters | N/A | NO |

## Manages read/unread messages in a channel or a thread

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/channels/{channelID}/messages/mark-as-read` | HTTP/S | GET | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| channelID | uint64 | PATH | Channel ID | N/A | YES |
| threadID | uint64 | POST | ID of thread (messageID)  | N/A | NO |
| lastReadMessageID | uint64 | POST | ID of the last read message | N/A | NO |

## Edit existing message

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/channels/{channelID}/messages/{messageID}` | HTTP/S | PUT | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| messageID | uint64 | PATH | Message ID | N/A | YES |
| channelID | uint64 | PATH | Channel ID | N/A | YES |
| message | string | POST | Message contents (markdown) | N/A | YES |

## Delete existing message

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/channels/{channelID}/messages/{messageID}` | HTTP/S | DELETE | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| messageID | uint64 | PATH | Message ID | N/A | YES |
| channelID | uint64 | PATH | Channel ID | N/A | YES |

## Reply to a message

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/channels/{channelID}/messages/{messageID}/replies` | HTTP/S | POST | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| messageID | uint64 | PATH | Message ID | N/A | YES |
| channelID | uint64 | PATH | Channel ID | N/A | YES |
| message | string | POST | Message contents (markdown) | N/A | YES |

## Pin message to channel (public bookmark)

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/channels/{channelID}/messages/{messageID}/pin` | HTTP/S | POST | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| messageID | uint64 | PATH | Message ID | N/A | YES |
| channelID | uint64 | PATH | Channel ID | N/A | YES |

## Pin message to channel (public bookmark)

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/channels/{channelID}/messages/{messageID}/pin` | HTTP/S | DELETE | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| messageID | uint64 | PATH | Message ID | N/A | YES |
| channelID | uint64 | PATH | Channel ID | N/A | YES |

## Bookmark a message (private bookmark)

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/channels/{channelID}/messages/{messageID}/bookmark` | HTTP/S | POST | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| messageID | uint64 | PATH | Message ID | N/A | YES |
| channelID | uint64 | PATH | Channel ID | N/A | YES |

## Remove boomark from message (private bookmark)

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/channels/{channelID}/messages/{messageID}/bookmark` | HTTP/S | DELETE | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| messageID | uint64 | PATH | Message ID | N/A | YES |
| channelID | uint64 | PATH | Channel ID | N/A | YES |

## React to a message

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/channels/{channelID}/messages/{messageID}/reaction/{reaction}` | HTTP/S | POST | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| messageID | uint64 | PATH | Message ID | N/A | YES |
| reaction | string | PATH | Reaction | N/A | YES |
| channelID | uint64 | PATH | Channel ID | N/A | YES |

## Delete reaction from a message

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/channels/{channelID}/messages/{messageID}/reaction/{reaction}` | HTTP/S | DELETE | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| messageID | uint64 | PATH | Message ID | N/A | YES |
| reaction | string | PATH | Reaction | N/A | YES |
| channelID | uint64 | PATH | Channel ID | N/A | YES |

---




# Permissions

| Method | Endpoint | Purpose |
| ------ | -------- | ------- |
| `GET` | `/permissions/effective` | Effective rules for current user |

## Effective rules for current user

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/permissions/effective` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| resource | string | GET | Show only rules for a specific resource | N/A | NO |

---




# Search entry point

| Method | Endpoint | Purpose |
| ------ | -------- | ------- |
| `GET` | `/search/messages` | Search for messages |

## Search for messages

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/search/messages` | HTTP/S | GET | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| channelID | []uint64 | GET | Filter by channels | N/A | NO |
| afterMessageID | uint64 | GET | ID of the first message in the list (exclusive) | N/A | NO |
| beforeMessageID | uint64 | GET | ID of the last message in the list (exclusive) | N/A | NO |
| fromMessageID | uint64 | GET | ID of the first message in the list (inclusive) | N/A | NO |
| toMessageID | uint64 | GET | ID of the last message the list (inclusive) | N/A | NO |
| threadID | []uint64 | GET | Filter by thread message ID | N/A | NO |
| userID | []uint64 | GET | Filter by one or more user | N/A | NO |
| type | []string | GET | Filter by message type (text, inlineImage, attachment, ...) | N/A | NO |
| pinnedOnly | bool | GET | Return only pinned messages | N/A | NO |
| bookmarkedOnly | bool | GET | Only bookmarked messages | N/A | NO |
| limit | uint | GET | Max number of messages | N/A | NO |
| query | string | GET | Search query | N/A | NO |

---




# Status

| Method | Endpoint | Purpose |
| ------ | -------- | ------- |
| `GET` | `/status/` | See all current statuses |
| `POST` | `/status/` | Set user's status |
| `DELETE` | `/status/` | Clear status |

## See all current statuses

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/status/` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |

## Set user's status

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/status/` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| icon | string | POST | Status icon | N/A | NO |
| message | string | POST | Status message | N/A | NO |
| expires | string | POST | Clear status when it expires (eg: when-active, afternoon, tomorrow 1h, 30m, 1 PM, 2019-05-20) | N/A | NO |

## Clear status

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/status/` | HTTP/S | DELETE |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |

---




# Webhooks

| Method | Endpoint | Purpose |
| ------ | -------- | ------- |
| `GET` | `/webhooks/` | List created webhooks |
| `POST` | `/webhooks/` | Create webhook |
| `POST` | `/webhooks/{webhookID}` | Attach file to channel |
| `GET` | `/webhooks/{webhookID}` | Get webhook details |
| `DELETE` | `/webhooks/{webhookID}` | Delete webhook |

## List created webhooks

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/webhooks/` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| channelID | uint64 | GET | Channel ID | N/A | NO |
| userID | uint64 | GET | Owner user ID | N/A | NO |

## Create webhook

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/webhooks/` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| channelID | uint64 | POST | Channel ID | N/A | YES |
| kind | types.WebhookKind | POST | Kind (incoming, outgoing) | N/A | YES |
| trigger | string | POST | Trigger word | N/A | NO |
| url | string | POST | POST URL | N/A | NO |
| username | string | POST | Default user name | N/A | NO |
| avatar | *multipart.FileHeader | POST | Default avatar | N/A | NO |
| avatarURL | string | POST | Default avatar (from URL) | N/A | NO |

## Attach file to channel

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/webhooks/{webhookID}` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| webhookID | uint64 | PATH | Webhook ID | N/A | YES |
| channelID | uint64 | POST | Channel ID | N/A | YES |
| kind | types.WebhookKind | POST | Kind (incoming, outgoing) | N/A | YES |
| trigger | string | POST | Trigger word | N/A | NO |
| url | string | POST | POST URL | N/A | NO |
| username | string | POST | Default user name | N/A | NO |
| avatar | *multipart.FileHeader | POST | Default avatar | N/A | NO |
| avatarURL | string | POST | Default avatar (from URL) | N/A | NO |

## Get webhook details

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/webhooks/{webhookID}` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| webhookID | uint64 | PATH | Webhook ID | N/A | YES |

## Delete webhook

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/webhooks/{webhookID}` | HTTP/S | DELETE |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| webhookID | uint64 | PATH | Webhook ID | N/A | YES |

---




# Webhooks (Public)

| Method | Endpoint | Purpose |
| ------ | -------- | ------- |
| `DELETE` | `/webhooks/{webhookID}/{webhookToken}` | Delete webhook |
| `POST` | `/webhooks/{webhookID}/{webhookToken}` | Create a message from a webhook payload |

## Delete webhook

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/webhooks/{webhookID}/{webhookToken}` | HTTP/S | DELETE |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| webhookID | uint64 | PATH | Webhook ID | N/A | YES |
| webhookToken | string | PATH | Authentication token | N/A | YES |

## Create a message from a webhook payload

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/webhooks/{webhookID}/{webhookToken}` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| username | string | GET | Custom username for webhook message | N/A | NO |
| avatarURL | string | GET | Custom avatar picture for webhook message | N/A | NO |
| content | string | GET | Message contents | N/A | YES |
| webhookID | uint64 | PATH | Webhook ID | N/A | YES |
| webhookToken | string | PATH | Authentication token | N/A | YES |

---