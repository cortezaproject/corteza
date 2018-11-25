# Authentication

## Check JWT token

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/auth/check` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |

## Delete JWT token (Sign Out)

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/auth/check` | HTTP/S | DELETE |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |




# Organisations

Organisations represent a top-level grouping entity. There may be many organisations defined in a single deployment.

## List organisations

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/organisations/` | HTTP/S | GET | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| query | string | GET | Search query | N/A | NO |

## Create organisation

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/organisations/` | HTTP/S | POST | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| name | string | POST | Organisation Name | N/A | YES |

## Update organisation details

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/organisations/{id}` | HTTP/S | PUT | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| id | uint64 | PATH | Organisation ID | N/A | NO |
| name | string | POST | Organisation Name | N/A | YES |

## Remove organisation

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/organisations/{id}` | HTTP/S | DELETE | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| id | uint64 | PATH | Organisation ID | N/A | YES |

## Read organisation details

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/organisations/{id}` | HTTP/S | GET | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| id | uint64 | GET | Organisation ID | N/A | YES |

## Archive organisation

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/organisations/{id}/archive` | HTTP/S | POST | Client ID, Session ID |

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
| `/teams/` | HTTP/S | GET | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| query | string | GET | Search query | N/A | NO |

## Update team details

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/teams/` | HTTP/S | POST | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| name | string | POST | Name of Team | N/A | YES |
| members | []uint64 | POST | Team member IDs | N/A | NO |

## Update team details

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/teams/{teamID}` | HTTP/S | PUT | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| teamID | uint64 | PATH | Team ID | N/A | YES |
| name | string | POST | Name of Team | N/A | NO |
| members | []uint64 | POST | Team member IDs | N/A | NO |

## Read team details and memberships

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/teams/{teamID}` | HTTP/S | GET | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| teamID | uint64 | PATH | Team ID | N/A | YES |

## Remove team

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/teams/{teamID}` | HTTP/S | DELETE | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| teamID | uint64 | PATH | Team ID | N/A | YES |

## Archive team

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/teams/{teamID}/archive` | HTTP/S | POST | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| teamID | uint64 | PATH | Team ID | N/A | YES |

## Move team to different organisation

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/teams/{teamID}/move` | HTTP/S | POST | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| teamID | uint64 | PATH | Team ID | N/A | YES |
| organisationID | uint64 | POST | Team ID | N/A | YES |

## Merge one team into another

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/teams/{teamID}/merge` | HTTP/S | POST | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| teamID | uint64 | PATH | Source Team ID | N/A | YES |
| destination | uint64 | POST | Destination Team ID | N/A | YES |

## Add member to a team

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/teams/{teamID}/memberAdd` | HTTP/S | POST | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| teamID | uint64 | PATH | Source Team ID | N/A | YES |
| userID | uint64 | POST | User ID | N/A | YES |

## Remove member from a team

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/teams/{teamID}/memberRemove` | HTTP/S | POST | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| teamID | uint64 | PATH | Source Team ID | N/A | YES |
| userID | uint64 | POST | User ID | N/A | YES |




# Users

## Search users (Directory)

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/users/` | HTTP/S | GET | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| query | string | GET | Search query to match against users | N/A | NO |

## Create user

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/users/` | HTTP/S | POST | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| email | string | POST | Email | N/A | YES |
| username | string | POST | Username | N/A | YES |
| password | string | POST | Password | N/A | YES |
| name | string | POST | Name | N/A | YES |
| handle | string | POST | Handle | N/A | YES |
| kind | string | POST | Kind (normal, bot) | N/A | NO |
| meta | types.JSONText | POST | Meta data | N/A | NO |
| satosaID | string | POST | Satosa ID | N/A | NO |
| organisationID | uint64 | POST | Organisation ID | N/A | NO |

## Update user details

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/users/{userID}` | HTTP/S | PUT | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| userID | uint64 | PATH | User ID | N/A | YES |
| email | string | POST | Email | N/A | YES |
| username | string | POST | Username | N/A | YES |
| password | string | POST | Password | N/A | YES |
| name | string | POST | Name | N/A | YES |
| handle | string | POST | Handle | N/A | YES |
| kind | string | POST | Kind (normal, bot) | N/A | NO |
| meta | types.JSONText | POST | Meta data | N/A | NO |
| satosaID | string | POST | Satosa ID | N/A | NO |
| organisationID | uint64 | POST | Organisation ID | N/A | NO |

## Read user details and memberships

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/users/{userID}` | HTTP/S | GET | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| userID | uint64 | PATH | User ID | N/A | YES |

## Remove user

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/users/{userID}` | HTTP/S | DELETE | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| userID | uint64 | PATH | User ID | N/A | YES |

## Suspend user

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/users/{userID}/suspend` | HTTP/S | POST | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| userID | uint64 | PATH | User ID | N/A | YES |

## Unsuspend user

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/users/{userID}/unsuspend` | HTTP/S | POST | Client ID, Session ID |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| userID | uint64 | PATH | User ID | N/A | YES |