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




# JSON for Charts

When the API generates a JSON like this we can generate charts (with some JS on the front-end). The Groupby, Sum and Count fields in the form in the chart module in the CRM are there to format the values in the data array.

## Example with no GroupBy, Sum or Count:

Example chart: select Opportunity value for all opportunities between 1/1/2018 and 31/12/2018 and show them in a line chart.

```
{
  "kind": "line",
  "data": [
    [
      "Close date",
      "2018-01-01",
      "2018-01-02",
      "2018-01-03",
      "2018-01-04"
    ],
    [
      "Opportunity value",
      30,
      200,
      100,
      400
    ]
  ],
  "fields": {
    "Close date": {
      "kind": "datetime",
      "min": "1/1/2018",
      "max": "31/12/2018"
    },
    "Opportunity value": {
      "name": "(Field name)",
      "kind": "(Field kind)"
    }
  },
  "meta": {
    "name": "Line chart",
    "description": "Description for the chart",
    "kind": "line",
    "module": 62511111111111110,
    "x": 62522222222222220,
    "xmin": "1/1/2018",
    "xmax": "31/12/2018",
    "y": 62533333333333336
  }
}
```

## Example with Count:

Example pie chart: select opportunity.state count by opportunity.state.

This will need to return:

Open, 10 (where the 10 is a count)
Won, 5
Lost, 2

```
{
  "kind": "pie",
  "data": [
    [
      "Open",
      10
    ],
    [
      "Won",
      5
    ],
    [
      "Lost",
      2
    ]
  ],
  "fields": {
    "state": {
      "kind": "string"
    }
  },
  "meta": {
    "name": "Pie chart",
    "description": "",
    "kind": "pie",
    "module": 62511111111111110,
    "x": 62522222222222220,
    "y": 62533333333333336,
    "count": 62566666666666664
  }
}
```

## Example with Sum:

Example donut chart: select opportunity.type sum by opportunity.value.

This will need to return a sum of all the opportunies by type (very similar to count, but it sums instead of counts)

```
{
  "kind": "donut",
  "data": [
    [
      "Open",
      120000
    ],
    [
      "Won",
      1000000
    ],
    [
      "Lost",
      30000
    ]
  ],
  "fields": {
    "type": {
      "kind": "string"
    }
  },
  "meta": {
    "name": "Donut chart",
    "description": "Description for the chart",
    "kind": "donut",
    "module": 62511111111111110,
    "x": 62522222222222220,
    "y": 62533333333333336,
    "sum": 62555555555555550
  }
}
```

## Example with GroupBy and Count

Example bar chart with number of leads per country: select lead.country from leads Count lead.id GroupBy lead.country

```
{
  "kind": "bar",
  "data": [
    [
      "France",
      12
    ],
    [
      "Slovenia",
      10
    ],
    [
      "Spain",
      3
    ],
    [
      "United Kingdom",
      5
    ],
    [
      "Germany",
      20
    ]
  ],
  "fields": {
    "country": {
      "kind": "string"
    }
  },
  "meta": {
    "name": "Bar chart",
    "description": "",
    "kind": "bar",
    "module": 62511111111111110,
    "x": 62522222222222220,
    "y": 62533333333333336,
    "groupby": 62544444444444450,
    "count": 62566666666666664
  }
}
```

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
| fields | types.ModuleFieldSet | POST | Fields JSON | N/A | YES |

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
| metrics | string | GET | Metrics (syntax: alias:expression;...) | N/A | YES |
| dimensions | string | GET | Dimensions (syntax: alias:field|modifier|modifier2;...) | N/A | YES |
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
| sort | string | GET | Sort field (default id desc) | N/A | NO |
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
| fields | sqlxTypes.JSONText | POST | Content JSON | N/A | YES |

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
| fields | sqlxTypes.JSONText | POST | Content JSON | N/A | YES |

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