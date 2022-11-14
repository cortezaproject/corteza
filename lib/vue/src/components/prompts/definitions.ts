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

export const prompts = Object.freeze([
  {
    ref: 'redirect',
    meta: { short: 'Redirect user to an outside URL' },
    parameters: [
      { name: 'url', types: ['String'], required: true },
      { name: 'delay', types: ['Integer'], meta: { description: 'Redirection delay in seconds' } },
    ],
  },
  {
    ref: 'reroute',
    meta: { short: 'Redirect user to an internal application route' },
    parameters: [
      { name: 'name', types: ['String'], required: true },
      { name: 'params', types: ['KV'] },
      { name: 'query', types: ['KV'] },
      { name: 'delay', types: ['Integer'], meta: { description: 'Redirection delay in seconds' } },
    ],
  },
  {
    ref: 'recordPage',
    meta: {
      short: 'Redirect user to the record page',
      webapps: ['compose'],
    },
    parameters: [
      { name: 'module', types: ['ID', 'Handle', 'ComposeModule'] },
      { name: 'namespace', types: ['ID', 'Handle', 'ComposeNamespace'] },
      { name: 'record', types: ['ID', 'ComposeRecord'] },
      { name: 'edit', types: ['Boolean'] },
      { name: 'delay', types: ['Integer'], meta: { description: 'Redirection delay in seconds' } },
    ],
  },
  {
    ref: 'notification',
    meta: { short: 'Show non-blocking message to user' },
    parameters: [
      { name: 'title', types: ['String'] },
      { name: 'message', types: ['String'], required: true },
      { name: 'variant', types: ['String'], meta: { visual: { input: { type: 'select', properties: { options: variants } } } } },
      { name: 'timeout', types: ['Integer'], meta: { description: 'How long do we show the notification in seconds' } },
    ],
  },
  {
    ref: 'alert',
    meta: { short: 'Prompt user with an alert' },
    parameters: [
      { name: 'title', types: ['String'] },
      { name: 'message', types: ['String'], required: true },
      { name: 'buttonLabel', types: ['String'] },
      { name: 'buttonVariant', types: ['String'], meta: { visual: { input: { type: 'select', properties: { options: variants } } } } },
      { name: 'buttonValue', types: ['Any'] },
    ],
  },
  {
    ref: 'choice',
    meta: { short: 'Prompt user with choice' },
    parameters: [
      { name: 'title', types: ['String'] },
      { name: 'message', types: ['String'], required: true },
      { name: 'confirmButtonLabel', types: ['String'] },
      { name: 'confirmButtonVariant', types: ['String'], meta: { visual: { input: { type: 'select', properties: { options: variants } } } } },
      { name: 'confirmButtonValue', types: ['Any'] },
      { name: 'rejectButtonLabel', types: ['String'] },
      { name: 'rejectButtonVariant', types: ['String'], meta: { visual: { input: { type: 'select', properties: { options: variants } } } } },
      { name: 'rejectButtonValue', types: ['Any'] },
    ],
    results: [
      { name: 'value', types: ['Any'] },
    ],
  },
  {
    ref: 'composeRecordPicker',
    meta: { short: 'Prompt user to select a Compose Record' },
    parameters: [
      { name: 'title', types: ['String'] },
      { name: 'message', types: ['String'], required: true },
      { name: 'namespace', types: ['ID', 'Handle', 'ComposeNamespace'], required: true },
      { name: 'module', types: ['ID', 'Handle', 'ComposeModule'], required: true },
      { name: 'labelField', types: ['Handle'], required: true },
      { name: 'queryFields', types: ['Array'] },
      { name: 'prefilter', types: ['String'] },
    ],
    results: [
      { name: 'value', types: ['ComposeRecord'] },
    ],
  },
  {
    ref: 'input',
    meta: { short: 'Prompt user with a single input' },
    parameters: [
      { name: 'title', types: ['String'] },
      { name: 'variant', types: ['String'], meta: { visual: { input: { type: 'select', properties: { options: variants } } } } },
      { name: 'message', types: ['String'], required: true },
      { name: 'label', types: ['String'] },
      {
        name: 'type',
        types: ['String'],
        meta: {
          visual: {
            options: [
              { value: 'text', text: 'Text' },
              { value: 'number', text: 'Number' },
              { value: 'email', text: 'Email' },
              { value: 'password', text: 'Password' },
              { value: 'search', text: 'Search' },
              { value: 'date', text: 'Date' },
              { value: 'time', text: 'Time' },
            ],
          },
        },
      },
      { name: 'inputValue', types: ['String'] },
    ],
    results: [
      { name: 'value', types: ['Any'] },
    ],
  },
  {
    ref: 'options',
    meta: { short: 'Prompt user with options' },
    parameters: [
      { name: 'title', types: ['String'] },
      { name: 'variant', types: ['String'], meta: { visual: { input: { type: 'select', properties: { options: variants } } } } },
      { name: 'message', types: ['String'], required: true },
      { name: 'label', types: ['String'] },
      {
        name: 'type',
        types: ['String'],
        meta: {
          visual: {
            options: [
              { value: 'select', text: 'Select' },
              { value: 'radio', text: 'Radio' },
            ],
          },
        },
      },
      { name: 'value', types: ['String', 'Array'] },
      { name: 'options', types: ['KV'] },
      { name: 'multiselect', types: ['Boolean'] },
    ],
    results: [
      { name: 'value', types: ['Any'] },
    ],
  },
].map(f => new automation.Function({ ...f, kind: 'prompt' })))
