const stdEventTypes = [
  'onManual',
]

const recordPageEventTypes = [
  'beforeFormSubmit',
  'onFormSubmitError',
  'afterFormSubmit',
  'beforeDelete',
  'afterDelete',
]

export default {
  compose: [
    ...stdEventTypes,
  ],
  'compose:namespace': [
    ...stdEventTypes,
  ],
  'compose:module': [
    ...stdEventTypes,
  ],
  'compose:record': [
    ...stdEventTypes,
  ],
  'ui:compose:record-page': [
    ...stdEventTypes,
    ...recordPageEventTypes,
  ],
  'ui:compose:admin-record-page': [
    ...stdEventTypes,
    ...recordPageEventTypes,
  ],
  'ui:compose': [
    ...stdEventTypes,
  ],
}
