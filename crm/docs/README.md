# Fields

CRM input field definitions

## List available fields

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/field/` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |

## Get field details

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/field/{id}` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| id | string | PATH | Type ID | N/A | YES |




# Modules

CRM module definitions

## List modules

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/module/` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| query | string | GET | Search query | N/A | NO |

## Create module

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/module/` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| name | string | POST | Module Name | N/A | YES |
| fields | types.JSONText | POST | Fields JSON | N/A | YES |

## Read module

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/module/{id}` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| id | uint64 | PATH | Module ID | N/A | YES |

## Edit module

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/module/{id}` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| id | uint64 | PATH | Module ID | N/A | YES |
| name | string | POST | Module Name | N/A | YES |
| fields | types.JSONText | POST | Fields JSON | N/A | YES |

## Delete module

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/module/{id}` | HTTP/S | DELETE |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| id | uint64 | PATH | Module ID | N/A | YES |

## List/read contents from module section

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/module/{module}/content` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| module | uint64 | PATH | Module ID | N/A | YES |

## List/read contents from module section

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/module/{module}/content` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| module | uint64 | PATH | Module ID | N/A | YES |
| fields | types.JSONText | POST | Content JSON | N/A | YES |

## Read contents by ID from module section

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/module/{module}/content/{id}` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| module | uint64 | PATH | Module ID | N/A | YES |
| id | uint64 | PATH | Content ID | N/A | YES |

## Add/update contents in module section

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/module/{module}/content/{id}` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| module | uint64 | PATH | Module ID | N/A | YES |
| id | uint64 | PATH | Content ID | N/A | YES |
| fields | types.JSONText | POST | Content JSON | N/A | YES |

## Delete content row from module section

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/module/{module}/content/{id}` | HTTP/S | DELETE |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| module | uint64 | PATH | Module ID | N/A | YES |
| id | uint64 | PATH | Content ID | N/A | YES |