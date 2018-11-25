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
