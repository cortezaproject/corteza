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
| `/field/{typeID}` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| typeID | string | PATH | Type ID | N/A | YES |




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
| `/module/{moduleID}` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| moduleID | uint64 | PATH | Module ID | N/A | YES |

## Analyze data for chart

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/module/{moduleID}/chart` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| name | string | GET | The chart name | N/A | YES |
| description | string | GET | The chart description | N/A | YES |
| xAxis | string | GET | X axis value | N/A | YES |
| xMin | string | GET | Min value | N/A | NO |
| xMax | string | GET | Max value | N/A | NO |
| yAxis | string | GET | Y axis value | N/A | YES |
| groupBy | string | GET | Group by field | N/A | YES |
| sum | string | GET | Sum values field | N/A | YES |
| count | string | GET | Count values field | N/A | YES |
| kind | string | GET | Chart kind (line, spline, step, area, area-spline, area-step, bar, scatter, pie, donut, gauge) | N/A | YES |
| moduleID | uint64 | PATH | Module ID | N/A | YES |

## Edit module

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/module/{moduleID}` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| moduleID | uint64 | PATH | Module ID | N/A | YES |
| name | string | POST | Module Name | N/A | YES |
| fields | types.JSONText | POST | Fields JSON | N/A | YES |

## Delete module

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/module/{moduleID}` | HTTP/S | DELETE |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| moduleID | uint64 | PATH | Module ID | N/A | YES |

## List/read contents from module section

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/module/{moduleID}/content` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| query | string | GET | Search query | N/A | NO |
| page | int | GET | Page number (0 based) | N/A | NO |
| perPage | int | GET | Returned items per page (default 50) | N/A | NO |
| moduleID | uint64 | PATH | Module ID | N/A | YES |

## List/read contents from module section

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/module/{moduleID}/content` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| moduleID | uint64 | PATH | Module ID | N/A | YES |
| fields | types.JSONText | POST | Content JSON | N/A | YES |

## Read contents by ID from module section

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/module/{moduleID}/content/{contentID}` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| moduleID | uint64 | PATH | Module ID | N/A | YES |
| contentID | uint64 | PATH | Content ID | N/A | YES |

## Add/update contents in module section

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/module/{moduleID}/content/{contentID}` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| moduleID | uint64 | PATH | Module ID | N/A | YES |
| contentID | uint64 | PATH | Content ID | N/A | YES |
| fields | types.JSONText | POST | Content JSON | N/A | YES |

## Delete content row from module section

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/module/{moduleID}/content/{contentID}` | HTTP/S | DELETE |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| moduleID | uint64 | PATH | Module ID | N/A | YES |
| contentID | uint64 | PATH | Content ID | N/A | YES |




# Pages

CRM module pages

## List available pages

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/page/` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| selfID | uint64 | GET | Parent page ID | N/A | NO |

## Create page

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/page/` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| selfID | uint64 | POST | Parent Page ID | N/A | NO |
| moduleID | uint64 | POST | Module ID (optional) | N/A | NO |
| title | string | POST | Title | N/A | YES |
| description | string | POST | Description | N/A | NO |
| visible | bool | POST | Visible in navigation | N/A | NO |
| blocks | types.JSONText | POST | Blocks JSON | N/A | YES |

## Get page details

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/page/{pageID}` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| pageID | uint64 | PATH | Page ID | N/A | YES |

## Edit page

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/page/{pageID}` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| pageID | uint64 | PATH | Page ID | N/A | YES |
| selfID | uint64 | POST | Parent Page ID | N/A | NO |
| moduleID | uint64 | POST | Module ID (optional) | N/A | NO |
| title | string | POST | Title | N/A | YES |
| description | string | POST | Description | N/A | NO |
| visible | bool | POST | Visible in navigation | N/A | NO |
| blocks | types.JSONText | POST | Blocks JSON | N/A | YES |

## Reorder pages

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/page/{selfID}/reorder` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| selfID | uint64 | PATH | Parent page ID | N/A | YES |
| pageIDs | []uint64 | POST | Page ID order | N/A | YES |

## Delete page

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/page/{pageID}` | HTTP/S | Delete |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| pageID | uint64 | PATH | Page ID | N/A | YES |