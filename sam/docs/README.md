# Organisations

Organisations represent a top-level grouping entity. There may be many organisations defined in a single deployment.

## List organisations

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/organisation/` | HTTP/S | GET | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| query | string | GET | Search query | N/A | NO |

## Create organisation

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/organisation/` | HTTP/S | PUT | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| name | string | POST | Organisation Name | N/A | YES |

## Update organisation details

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/organisation/{id}` | HTTP/S | POST | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| id | uint64 | PATH | Organisation ID | N/A | NO |
| name | string | POST | Organisation Name | N/A | YES |

## Remove organisation

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/organisation/{id}` | HTTP/S | DELETE | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| id | uint64 | PATH | Organisation ID | N/A | YES |

## Read organisation details

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/organisation/{id}` | HTTP/S | GET | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| id | uint64 | GET | Organisation ID | N/A | YES |

## Archive organisation

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/organisation/{id}/archive` | HTTP/S | POST | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| id | uint64 | PATH | Organisation ID | N/A | YES |




# Teams

An organisation may have many teams. Teams may have many channels available. Access to channels may be shared between teams.

## List teams

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/team/` | HTTP/S | GET | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| query | string | GET | Search query | N/A | NO |

## Update team details

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/team/` | HTTP/S | PUT | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| name | string | POST | Name of Team | N/A | YES |
| members | []uint64 | POST | Team member IDs | N/A | NO |

## Update team details

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/team/{id}` | HTTP/S | POST | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| id | uint64 | PATH | Team ID | N/A | YES |
| name | string | POST | Name of Team | N/A | NO |
| members | []uint64 | POST | Team member IDs | N/A | NO |

## Read team details and memberships

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/team/{id}` | HTTP/S | GET | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| id | uint64 | PATH | Organisation ID | N/A | YES |

## Remove team

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/team/{id}` | HTTP/S | DELETE | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| id | uint64 | PATH | Organisation ID | N/A | YES |

## Archive team

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/team/{id}/archive` | HTTP/S | POST | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| id | uint64 | PATH | Organisation ID | N/A | YES |

## Move team to different organisation

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/team/{id}/move` | HTTP/S | POST | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| id | uint64 | PATH | Organisation ID | N/A | YES |
| organisation_id | uint64 | POST | Organisation ID | N/A | YES |

## Merge one team into another

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/team/{id}/merge` | HTTP/S | POST | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| id | uint64 | PATH | Source Channel ID | N/A | YES |
| destination | uint64 | POST | Destination Channel ID | N/A | YES |




# Channels

A channel is a representation of a sequence of messages. It has meta data like channel subject. Channels may be public, private or direct (between two users).

## List channels

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/channel/` | HTTP/S | GET | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| query | string | GET | Search query | N/A | NO |

## Create new channel

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/channel/` | HTTP/S | PUT | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| name | string | POST | Name of Channel | N/A | YES |
| topic | string | POST | Subject of Channel | N/A | NO |

## Update channel details

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/channel/edit` | HTTP/S | POST | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| id | uint64 | POST | Channel ID | N/A | YES |
| name | string | POST | Name of Channel | N/A | NO |
| topic | string | POST | Subject of Channel | N/A | NO |
| archive | bool | POST | Request channel to be archived or unarchived | N/A | NO |
| organisationId | uint64 | POST | Move channel to different organisation | N/A | NO |

## Read channel details

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/channel/read` | HTTP/S | GET | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| id | uint64 | GET | Channel ID | N/A | YES |

## Remove channel

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/channel/delete` | HTTP/S | DELETE | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| id | uint64 | GET | Channel ID | N/A | YES |




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

## New message / edit message

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/message/edit` | HTTP/S | POST | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| id | uint64 | POST | Message ID | N/A | NO |
| channel_id | uint64 | POST | Channel ID where to post message | N/A | NO |
| contents | string | POST | Message contents (markdown) | N/A | YES |

## Attach file to message

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/message/attach` | HTTP/S | PUT | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |

## Remove message

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/message/remove` | HTTP/S | DELETE | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| id | uint64 | GET | Message ID | N/A | YES |

## Read message details

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/message/read` | HTTP/S | GET | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| channel_id | uint64 | POST | Channel ID to read messages from | N/A | YES |

## Search messages

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/message/search` | HTTP/S | GET | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| query | string | GET | Search string to match against messages | N/A | NO |
| message_type | string | GET | Limit results to message type | N/A<br><br>Values:<br><br><ul><li>`history`</li><li>`message`</li><li>`attachment`</li><li>`media`</li> | NO |

## Pin message to channel (public bookmark)

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/message/pin` | HTTP/S | POST | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| id | uint64 | POST | Message ID | N/A | YES |

## Flag message for user (bookmark)

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/message/flag` | HTTP/S | POST | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| id | uint64 | POST | Message ID | N/A | YES |




# Users

## User login

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/user/login` | HTTP/S | POST | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| username | string | POST | Username or email | N/A | YES |
| password | string | POST | Password for user | N/A | YES |

## Search users (Directory)

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/user/search` | HTTP/S | GET | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| query | string | GET | Search query to match against users | N/A | NO |




# Websocket

## Bidirectional websocket chat connection

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/websocket/client` | WebSockets | GET | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |