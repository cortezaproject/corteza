# Organisations

Organisations represent a top-level grouping entity. There may be many organisations defined in a single deployment.

## Add new organisation

## Archive organisation

## Update organisation details

## Read organisation details

## Remove organisation

## Search organisations

# Teams

An organisation may have many teams. Teams may have many channels available. Access to channels may be shared between teams.

## Add new team

## Archive team

## Update team details

## Move team to different organisation

## Read team details

## Remove team

## Search teams

# Channels

A channel is a representation of a sequence of messages. It has meta data like channel subject. Channels may be public, private or direct (between two users).

## Add new channel

## Archive channel

## Update channel details

## Move channel to different team or organisation

## Read channel details

## Remove channel

## Search channels

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

## Add new message

## Update message details

## Flag message for user (bookmark)

## Pin message to channel (public bookmark)

## Read message details

## Remove message

## Search messages

# Members

## Add new member

## Update member details

## Read member details

## Remove member

## Search members (Directory)

# Files

The Files API is an abstraction over messages that have been sent with a file attachment.

## Add new file

## Read file

## Remove file

## Search files
