---
name: chart_building
description: Making a Compose chart and putting it on a page — what a report needs before it draws anything, and the two omissions that render an empty frame.
triggers:
  - compose_chart_create
  - compose_chart_update
  - compose_chart_lookup
  - compose_chart_delete
  - compose_chart_undelete
---

# Charts

A chart is its own resource, not a page block. Two steps, always in this order:

1. `compose_chart_create` in the namespace, with the module the data comes from.
2. `compose_page_create` / `compose_page_update` with a Chart block naming it.

Charts are fully creatable through these tools. Never tell a user that charts
have to be made by hand.

## The shape of a report

A chart holds `reports`; each report is one module, one filter, its dimensions
and its metrics.

```
{"colorScheme":"tableau.Tableau10","reports":[{"moduleID":"<numeric ID>","filter":"","dimensions":[{"field":"stage","modifier":"(no grouping / buckets)","conditions":{}}],"metrics":[{"field":"count","type":"doughnut"}]}]}
```

- **`moduleID` must be the numeric ID, as a string.** A handle is refused, and
  so is an ID no module in the namespace has — the chart would store cleanly and
  draw nothing. Get the ID from `compose_module_lookup`.
- **A dimension is what the data is grouped by.** A Select field groups well.
  `modifier` is `(no grouping / buckets)` for the field's own values, or DATE,
  WEEK, MONTH, QUARTER or YEAR to bucket a date field.
- **A metric is what is measured.** `{"field":"count"}` counts records and
  takes no aggregate. Any other field is a numeric field being reduced and
  **must** carry `"aggregate"`: SUM, MAX, MIN or AVG. Leaving it out is the
  usual way to get "Metrics aggregate not defined" instead of a chart.
- **`type` on the metric is what shapes the chart** — pie, doughnut, bar, line,
  funnel, gauge, radar, scatter. There is no separate chart-type field.

`colorScheme` is optional and its swatch count is part of the name, so
`tableau.ClassicOrangeBlue13` exists and `…Blue7` does not. An unknown name is
refused with the near matches rather than stored, because the webapp would
resolve it to no palette and draw a legend with no visible series.

## Putting it on a page

```json
[
  {
    "kind": "Chart",
    "title": "By stage",
    "options": {
      "chartID": "by_stage"
    }
  }
]
```

`chartID` takes the chart's name, handle or ID and the server resolves it. It
is required: a Chart block without one saves and draws an empty panel, and
`compose_page_create` says so in its returned note.

## Reading the result

`compose_chart_create` echoes the stored config with a `reportID`, `metricID`
and `dimensionID` minted for you. Send those back on `compose_chart_update` —
the whole `reports` array replaces, so an update that omits them re-creates the
report under new IDs.
