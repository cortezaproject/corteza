# Organisations

Organisations represent a top-level grouping entity. There may be many organisations defined in a single deployment.

## Archive organisation

### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| id | uint64 | POST | ID of Organisation | N/A | YES |

## Update organisation details

### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| id | uint64 | POST | ID of Organisation | N/A | NO |
| name | string | POST | Name of Organisation | N/A | YES |

## Read organisation details

### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| id | uint64 | GET | ID of Organisation | N/A | YES |

## Remove organisation

### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| id | uint64 | GET | ID of Organisation | N/A | YES |

## Search organisations

### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| query | string | GET | Search query | N/A | NO |

# Teams

An organisation may have many teams. Teams may have many channels available. Access to channels may be shared between teams.

## Archive team

### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| id | uint64 | POST | ID of Organisation | N/A | YES |

## Update team details

### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| id | uint64 | POST | ID of Team | N/A | NO |
| name | string | POST | Name of Team | N/A | YES |
| members | []uint64 | POST | Team member IDs | N/A | NO |

## Move team to different organisation

### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| id | uint64 | POST | ID of Organisation | N/A | YES |
| organisation_id | uint64 | POST | ID of Organisation | N/A | YES |

## Read team details and memberships

### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| id | uint64 | GET | ID of Organisation | N/A | YES |

## Remove team

### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| id | uint64 | GET | ID of Organisation | N/A | YES |

## Search teams

### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| query | string | GET | Search query | N/A | NO |

# Channels

A channel is a representation of a sequence of messages. It has meta data like channel subject. Channels may be public, private or direct (between two users).

## Archive channel

### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| id | uint64 | POST | ID of Channel | N/A | YES |

## Update channel details

### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| id | uint64 | POST | ID of Channel | N/A | NO |
| name | string | POST | Name of Channel | N/A | YES |
| topic | string | POST | Subject of Channel | N/A | YES |

## Merge one team into another

### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| destination | uint64 | POST | Destination Channel ID | N/A | YES |
| source | uint64 | POST | Source Channel ID | N/A | YES |

## Move channel to different team or organisation

### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| id | uint64 | POST | ID of Channel | N/A | YES |

## Read channel details

### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| id | uint64 | GET | ID of Channel | N/A | YES |

## Remove channel

### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| id | uint64 | GET | ID of Channel | N/A | YES |

## Search channels

### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| query | string | GET | Search query | N/A | NO |

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

### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| id | uint64 | POST | ID of Message | N/A | NO |
| channel_id | uint64 | POST | ID of Channel where to post message | N/A | NO |
| contents | string | POST | Message contents (markdown) | N/A | YES |

## Flag message for user (bookmark)

### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| id | uint64 | POST | ID of Message | N/A | YES |

## Pin message to channel (public bookmark)

### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| id | uint64 | POST | ID of Message | N/A | YES |

## Read message details

### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| channel_id | uint64 | POST | ID of Channel to read messages from | N/A | YES |

## Remove message

### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| id | uint64 | GET | ID of Message | N/A | YES |

## Search messages

### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| query | string | GET | Search string to match against messages | N/A | NO |
| message_type | string | GET | Limit results to message type | N/A

Values:

- `history`
- `message`
- `attachment`
- `media`
| NO |

# Members

## Search members (Directory)

### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |

# Files

The Files API is an abstraction over messages that have been sent with a file attachment.

## Add new file

### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |

## Read file

### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |

## Remove file

### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |

## Search files

### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
