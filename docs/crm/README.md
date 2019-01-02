# Charts

## List/read charts from module section

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/chart/` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |

## List/read charts from module section

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/chart/` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| config | sqlxTypes.JSONText | POST | Chart JSON | N/A | YES |
| name | string | POST | Chart name | N/A | YES |

## Read charts by ID from module section

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/chart/{chartID}` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| chartID | uint64 | PATH | Chart ID | N/A | YES |

## Add/update charts in module section

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/chart/{chartID}` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| chartID | uint64 | PATH | Chart ID | N/A | YES |
| config | sqlxTypes.JSONText | POST | Chart JSON | N/A | YES |
| name | string | POST | Chart name | N/A | YES |

## Delete chart

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/chart/{chartID}` | HTTP/S | DELETE |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| chartID | uint64 | PATH | Chart ID | N/A | YES |




# Jobs

Workflow Jobs

## List jobs

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/job/` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| status | string | GET | Job status (`ok`, `error`, `running`, `cancelled` or `queued`) | N/A | NO |
| page | int | GET | Page number (0 based) | N/A | NO |
| perPage | int | GET | Returned items per page (default 50) | N/A | NO |

## Create a new job

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/job/` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| workflowID | string | POST | Workflow ID | N/A | YES |
| startAt | string | POST | Start datetime for a delayed job | N/A | NO |
| parameters | types.JobParameterSet | POST | Extra job parameters (map[string]string) | N/A | NO |

## Get job details

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/job/{jobID}` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| jobID | string | PATH | Job ID | N/A | YES |

## Get job logs

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/job/{jobID}/logs` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| jobID | string | PATH | Job ID | N/A | YES |
| page | int | PATH | Page number (0 based) | N/A | NO |
| perPage | int | PATH | Returned items per page (default 50) | N/A | NO |

## Update job details

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/job/{jobID}` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| jobID | string | PATH | Job ID | N/A | YES |
| status | string | POST | Job status (`ok`, `error`, `running`, `cancelled` or `queued`) | N/A | NO |
| log | sqlxTypes.JSONText | POST | Job log item (append-only) | N/A | NO |
| workflowID | string | POST | Workflow ID | N/A | NO |
| startAt | string | POST | Start datetime for a delayed job | N/A | NO |
| parameters | types.JobParameterSet | POST | Extra job parameters (map[string]string) | N/A | NO |

## Cancel job

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/job/{jobID}` | HTTP/S | DELETE |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| jobID | string | PATH | Job ID | N/A | YES |




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
| fields | types.ModuleFieldSet | POST | Fields JSON | N/A | YES |
| meta | sqlxTypes.JSONText | POST | Module meta data | N/A | YES |

## Read module

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/module/{moduleID}` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
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
| fields | types.ModuleFieldSet | POST | Fields JSON | N/A | YES |
| meta | sqlxTypes.JSONText | POST | Module meta data | N/A | YES |

## Delete module

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/module/{moduleID}` | HTTP/S | DELETE |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| moduleID | uint64 | PATH | Module ID | N/A | YES |

## Generates report from module records

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/module/{moduleID}/report` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| metrics | string | GET | Metrics (eg: 'COUNT(*) AS count, SUM(money)') | N/A | YES |
| dimensions | string | GET | Dimensions (eg: 'DATE(foo), status') | N/A | YES |
| filter | string | GET | Filter (eg: 'DATE(foo) > 2010') | N/A | YES |
| moduleID | uint64 | PATH | Module ID | N/A | YES |

## List/read records from module section

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/module/{moduleID}/record` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| filter | string | GET | Filtering condition | N/A | NO |
| page | int | GET | Page number (0 based) | N/A | NO |
| perPage | int | GET | Returned items per page (default 50) | N/A | NO |
| sort | string | GET | Sort field (default id desc) | N/A | NO |
| moduleID | uint64 | PATH | Module ID | N/A | YES |

## Create record in module section

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/module/{moduleID}/record` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| moduleID | uint64 | PATH | Module ID | N/A | YES |
| fields | sqlxTypes.JSONText | POST | Record JSON | N/A | YES |

## Read records by ID from module section

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/module/{moduleID}/record/{recordID}` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| moduleID | uint64 | PATH | Module ID | N/A | YES |
| recordID | uint64 | PATH | Record ID | N/A | YES |

## Add/update records in module section

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/module/{moduleID}/record/{recordID}` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| moduleID | uint64 | PATH | Module ID | N/A | YES |
| recordID | uint64 | PATH | Record ID | N/A | YES |
| fields | sqlxTypes.JSONText | POST | Record JSON | N/A | YES |

## Delete record row from module section

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/module/{moduleID}/record/{recordID}` | HTTP/S | DELETE |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| moduleID | uint64 | PATH | Module ID | N/A | YES |
| recordID | uint64 | PATH | Record ID | N/A | YES |




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
| moduleID | uint64 | POST | Module ID | N/A | NO |
| title | string | POST | Title | N/A | YES |
| description | string | POST | Description | N/A | NO |
| visible | bool | POST | Visible in navigation | N/A | NO |
| blocks | sqlxTypes.JSONText | POST | Blocks JSON | N/A | YES |

## Get page details

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/page/{pageID}` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| pageID | uint64 | PATH | Page ID | N/A | YES |

## Get page all (non-record) pages, hierarchically

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/page/tree` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |

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
| blocks | sqlxTypes.JSONText | POST | Blocks JSON | N/A | YES |

## Reorder pages

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/page/{selfID}/reorder` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| selfID | uint64 | PATH | Parent page ID | N/A | YES |
| pageIDs | []string | POST | Page ID order | N/A | YES |

## Delete page

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/page/{pageID}` | HTTP/S | Delete |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| pageID | uint64 | PATH | Page ID | N/A | YES |




# Workflows

CRM workflow definitions

## List available workflows

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/workflow/` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |

## Create new workflow

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/workflow/` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| name | string | POST | Workflow name | N/A | YES |
| tasks | types.WorkflowTaskSet | POST | Type ID | N/A | NO |
| onError | types.WorkflowTaskSet | POST | Type ID | N/A | NO |
| timeout | int | POST | Timeout in seconds | N/A | NO |

## Get workflow details

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/workflow/{workflowID}` | HTTP/S | GET |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| workflowID | string | PATH | Workflow ID | N/A | YES |

## Update workflow details

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/workflow/{workflowID}` | HTTP/S | POST |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| workflowID | string | PATH | Workflow ID | N/A | YES |
| name | string | POST | Workflow name | N/A | YES |
| tasks | types.WorkflowTaskSet | POST | Type ID | N/A | NO |
| onError | types.WorkflowTaskSet | POST | Type ID | N/A | NO |
| timeout | int | POST | Timeout in seconds | N/A | NO |

## Delete workflow

#### Method

| URI | Protocol | Method | Authentication |
| --- | -------- | ------ | -------------- |
| `/workflow/{workflowID}` | HTTP/S | DELETE |  |

#### Request parameters

| Parameter | Type | Method | Description | Default | Required? |
| --------- | ---- | ------ | ----------- | ------- | --------- |
| workflowID | string | PATH | Workflow ID | N/A | YES |