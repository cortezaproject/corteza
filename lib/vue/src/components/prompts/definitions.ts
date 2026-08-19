import { automation } from '@cortezaproject/corteza-js'

const variants = [
  { value: 'primary', text: 'Primary' },
  { value: 'secondary', text: 'Secondary' },
  { value: 'success', text: 'Success' },
  { value: 'warning', text: 'Warning' },
  { value: 'danger', text: 'Danger' },
  { value: 'info', text: 'Info' },
  { value: 'light', text: 'Light' },
  { value: 'dark', text: 'Dark' },
]

const openModeVariants = [
  { value: 'sameTab', text: 'Open link in the same tab' },
  { value: 'newTab', text: 'Open link in a new tab' },
]

export const prompts = Object.freeze([
  {
    ref: 'redirect',
    meta: { short: 'Redirect user to an outside URL' },
    parameters: [
      { name: 'owner', types: ['User', 'ID'], required: false, meta: { label: 'Owner' } },
      { name: 'url', types: ['String'], required: true, meta: { label: 'URL' } },
      { name: 'delay', types: ['Integer'], meta: { label: 'Delay', description: 'Redirection delay in seconds' } },
      { name: 'openMode', types: ['String'], meta: { label: 'Open mode', visual: { input: { type: 'select', properties: { options: openModeVariants }, default: 'sameTab' } } } },
    ],
  },
  {
    ref: 'reroute',
    meta: { short: 'Redirect user to an internal application route' },
    parameters: [
      { name: 'owner', types: ['User', 'ID'], required: false, meta: { label: 'Owner' } },
      { name: 'name', types: ['String'], required: true, meta: { label: 'Name' } },
      { name: 'params', types: ['KV'], meta: { label: 'Parameters' } },
      { name: 'query', types: ['KV'], meta: { label: 'Query' } },
      { name: 'delay', types: ['Integer'], meta: { label: 'Delay', description: 'Redirection delay in seconds' } },
      { name: 'openMode', types: ['String'], meta: { label: 'Open mode', visual: { input: { type: 'select', properties: { options: openModeVariants }, default: 'sameTab' } } } },
    ],
  },
  {
    ref: 'recordPage',
    meta: {
      short: 'Redirect user to the record page',
      webapps: ['compose'],
    },
    parameters: [
      { name: 'owner', types: ['User', 'ID'], required: false, meta: { label: 'Owner' } },
      { name: 'module', types: ['ID', 'Handle', 'ComposeModule'], meta: { label: 'Module' } },
      { name: 'namespace', types: ['ID', 'Handle', 'ComposeNamespace'], meta: { label: 'Namespace' } },
      { name: 'record', types: ['ID', 'ComposeRecord'], meta: { label: 'Record' } },
      { name: 'values', types: ['KV'], meta: { label: 'Values', description: 'Prefilled record values for new records' } },
      { name: 'edit', types: ['Boolean'], meta: { label: 'Edit' } },
      { name: 'delay', types: ['Integer'], meta: { label: 'Delay', description: 'Redirection delay in seconds' } },
      { name: 'openMode', types: ['String'], meta: { label: 'Open mode', visual: { input: { type: 'select', properties: { options: [...openModeVariants, { value: 'modal', text: 'Open in a modal' }] }, default: 'modal' } } } },
    ],
  },
  {
    ref: 'refetchRecords',
    meta: { short: 'Refresh all record values on the page' },
  },
  {
    ref: 'notification',
    meta: { short: 'Show non-blocking message to user' },
    parameters: [
      { name: 'owner', types: ['User', 'ID'], required: false, meta: { label: 'Owner' } },
      { name: 'title', types: ['String'], meta: { label: 'Title' } },
      { name: 'message', types: ['String'], required: true, meta: { label: 'Message' } },
      { name: 'variant', types: ['String'], meta: { label: 'Variant', visual: { input: { type: 'select', properties: { options: variants }, default: 'primary' } } } },
      { name: 'timeout', types: ['Integer'], meta: { label: 'Timeout', description: 'How long do we show the notification in seconds' } },
    ],
  },
  {
    ref: 'alert',
    meta: { short: 'Prompt user with an alert' },
    parameters: [
      { name: 'owner', types: ['User', 'ID'], required: false, meta: { label: 'Owner' } },
      { name: 'title', types: ['String'], meta: { label: 'Title' } },
      { name: 'message', types: ['String'], required: true, meta: { label: 'Message' } },
      { name: 'buttonLabel', types: ['String'], meta: { label: 'Button label' } },
      { name: 'buttonVariant', types: ['String'], meta: { label: 'Button variant', visual: { input: { type: 'select', properties: { options: variants }, default: 'primary' } } } },
    ],
  },
  {
    ref: 'choice',
    meta: { short: 'Prompt user with choice' },
    parameters: [
      { name: 'owner', types: ['User', 'ID'], required: false, meta: { label: 'Owner' } },
      { name: 'title', types: ['String'], meta: { label: 'Title' } },
      { name: 'message', types: ['String'], required: true, meta: { label: 'Message' } },
      { name: 'confirmButtonLabel', types: ['String'], meta: { label: 'Confirm button label' } },
      { name: 'confirmButtonVariant', types: ['String'], meta: { label: 'Confirm button variant', visual: { input: { type: 'select', properties: { options: variants }, default: 'primary' } } } },
      { name: 'confirmButtonValue', types: ['Any'], meta: { label: 'Confirm button value' } },
      { name: 'rejectButtonLabel', types: ['String'], meta: { label: 'Reject button label' } },
      { name: 'rejectButtonVariant', types: ['String'], meta: { label: 'Reject button variant', visual: { input: { type: 'select', properties: { options: variants }, default: 'danger' } } } },
      { name: 'rejectButtonValue', types: ['Any'], meta: { label: 'Reject button value' } },
    ],
    results: [
      { name: 'value', types: ['Any'] },
    ],
  },
  {
    ref: 'composeRecordPicker',
    meta: { short: 'Prompt user to select a Compose Record' },
    parameters: [
      { name: 'owner', types: ['User', 'ID'], required: false, meta: { label: 'Owner' } },
      { name: 'title', types: ['String'], meta: { label: 'Title' } },
      { name: 'message', types: ['String'], required: true, meta: { label: 'Message' } },
      { name: 'namespace', types: ['ID', 'Handle', 'ComposeNamespace'], required: true, meta: { label: 'Namespace' } },
      { name: 'module', types: ['ID', 'Handle', 'ComposeModule'], required: true, meta: { label: 'Module' } },
      { name: 'labelField', types: ['Handle'], required: true, meta: { label: 'Label field' } },
      { name: 'queryFields', types: ['Array'], meta: { label: 'Query fields' } },
      { name: 'prefilter', types: ['String'], meta: { label: 'Prefilter' } },
      { name: 'label', types: ['String'], meta: { label: 'Label' } },
      { name: 'placeholder', types: ['String'], meta: { label: 'Placeholder' } },
      { name: 'buttonLabel', types: ['String'], meta: { label: 'Button label' } },
    ],
    results: [
      { name: 'value', types: ['ComposeRecord'] },
    ],
  },
  {
    ref: 'input',
    meta: { short: 'Prompt user with a single input' },
    parameters: [
      { name: 'owner', types: ['User', 'ID'], required: false, meta: { label: 'Owner' } },
      { name: 'title', types: ['String'], meta: { label: 'Title' } },
      { name: 'variant', types: ['String'], meta: { label: 'Variant', visual: { input: { type: 'select', properties: { options: variants }, default: 'primary' } } } },
      { name: 'message', types: ['String'], required: true, meta: { label: 'Message' } },
      { name: 'label', types: ['String'], meta: { label: 'Label' } },
      {
        name: 'type',
        types: ['String'],
        meta: { label: 'Type',
          visual: {
            input: {
              type: 'select',
              properties: {
                options: [
                  { value: 'text', text: 'Text' },
                  { value: 'number', text: 'Number' },
                  { value: 'email', text: 'Email' },
                  { value: 'password', text: 'Password' },
                  { value: 'search', text: 'Search' },
                  { value: 'datetime', text: 'Date & Time' },
                  { value: 'date', text: 'Date' },
                  { value: 'time', text: 'Time' },
                ],
              },
              default: 'text',
            },
          },
        },
      },
      { name: 'inputValue', types: ['String'], meta: { label: 'Input value' } },
      { name: 'buttonLabel', types: ['String'], meta: { label: 'Button label' } },
    ],
    results: [
      { name: 'value', types: ['Any'] },
    ],
  },
  {
    ref: 'options',
    meta: { short: 'Prompt user with options' },
    parameters: [
      { name: 'owner', types: ['User', 'ID'], required: false, meta: { label: 'Owner' } },
      { name: 'title', types: ['String'], meta: { label: 'Title' } },
      { name: 'variant', types: ['String'], meta: { label: 'Variant', visual: { input: { type: 'select', properties: { options: variants }, default: 'primary' } } } },
      { name: 'message', types: ['String'], required: true, meta: { label: 'Message' } },
      { name: 'label', types: ['String'], meta: { label: 'Label' } },
      {
        name: 'type',
        types: ['String'],
        meta: { label: 'Type',
          visual: {
            input: {
              type: 'select',
              properties: {
                options: [
                  { value: 'select', text: 'Select' },
                  { value: 'radio', text: 'Radio' },
                ],
              },
              default: 'select',
            },
          },
        },
      },
      { name: 'value', types: ['String', 'Array'], meta: { label: 'Value' } },
      { name: 'placeholder', types: ['String'], meta: { label: 'Placeholder' } },
      { name: 'options', types: ['KV'], meta: { label: 'Options' } },
      { name: 'multiselect', types: ['Boolean'], meta: { label: 'Multiselect' } },
      { name: 'buttonLabel', types: ['String'], meta: { label: 'Button label' } },
    ],
    results: [
      { name: 'value', types: ['Any'] },
    ],
  },
].map(f => new automation.Function({ ...f, kind: 'prompt' })))
