# Attachments

## Serves attached file

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/attachment/{attachmentID}/original/{name}` | HTTP/S | GET | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| download | bool | GET | Force file download | N/A | NO |
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




# Channels

A channel is a representation of a sequence of messages. It has meta data like channel subject. Channels may be public, private or group.

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

## All messages (channel history)

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/channels/{channelID}/messages/` | HTTP/S | GET | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| lastMessageID | uint64 | GET |  | N/A | NO |
| channelID | uint64 | PATH | Channel ID | N/A | YES |

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

## Returns all replies to a message

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/channels/{channelID}/messages/{messageID}/replies` | HTTP/S | GET | Client ID, Session ID |

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




# Permissions

## List default permissions

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/permissions/` | HTTP/S | GET | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |

## Retrieve current permission settings

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/permissions/{teamID}` | HTTP/S | GET | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| resource | string | GET | Permissions resource | N/A | YES |
| teamID | uint64 | PATH | Team ID | N/A | YES |

## Update permission settings

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/permissions/{teamID}` | HTTP/S | POST | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| teamID | uint64 | PATH | Team ID | N/A | YES |
| permissions | []rules.Rules | POST | List of rules to set | N/A | YES |

## List resources for given scope

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/permissions/scopes/{scope}` | HTTP/S | GET | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| scope | string | PATH | Scope | N/A | YES |




# Search entry point

## Search for messages

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/search/messages` | HTTP/S | GET | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| inChannel | uint64 | GET | Search only in one channel | N/A | NO |
| fromUser | uint64 | GET | Search only from one user | N/A | NO |
| firstID | uint64 | GET | Paging; return newer messages only (higher id) | N/A | NO |
| lastID | uint64 | GET | Paging; return older messages only (lower id) | N/A | NO |
| query | string | GET | Search query | N/A | NO |