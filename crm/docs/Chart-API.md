# JSON for Charts

When the API generates a JSON like this we can generate charts (with some JS on the front-end). The Groupby, Sum and Count fields in the form in the chart module in the CRM are there to format the values in the data array.

## Example with no GroupBy, Sum or Count:

Example chart: select Opportunity value for all opportunities between 1/1/2018 and 31/12/2018 and show them in a line chart. 

{
  "kind": "line",
  "data": [
      ["Close date", "2018-01-01", "2018-01-02", "2018-01-03", "2018-01-04"],
      ["Opportunity value", 30, 200, 100, 400]
    ]
  ,
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
  }
}

## Example with Count:

Example pie chart: select opportunity.state count by opportunity.state. 

This will need to return:

Open, 10 (where the 10 is a count)
Won, 5
Lost, 2

{
  "kind": "pie",
  "data": [
      ["Open", 10],
      ["Won", 5],
      ["Lost", 2]
    ]
  ,
  "fields": {
    "opportunity.state": {
      "kind": "datetime",
    }
  }
}

## Example with Sum:

Example donut chart: select opportunity.type sum by opportunity.value. 

This will need to return a sum of all the opportunies by type (very similar to count, but it sums instead of counts)

{
  "kind": "donut",
  "data": [
      ["Open", 120000],
      ["Won", 1000000],
      ["Lost", 30000]
    ]
  ,
  "fields": {
    "value": {
      "kind": "currency",
    }
  }
}

## Example with GroupBy and Count

Example bar chart with number of leads per country: select lead.country from leads Count lead.id GroupBy lead.country

{
  "kind": "bar",
  "data": [
      ["France", 12],
      ["Slovenia", 10],
      ["Spain", 3],
      ["United Kingdom", 5],
      ["Germany", 20']
    ]
  ,
  "fields": {
    "country": {
      "kind": "string",
    }
  }
}
