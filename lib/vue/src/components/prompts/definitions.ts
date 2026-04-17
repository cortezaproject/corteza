import { automation } from '@planetcrust/human-js'

const variants = [
  { value: 'primary', text: 'Primary' },
  { value: 'secondary', text: 'Secondary' },
  { value: 'success', text: 'Success' },
  { value: 'warning', text: 'Warning' },
  { value: 'danger', text: 'Danger' },
  { value: 'info', text: 'Info' },
]

const openModeVariants = [
  { value: 'sameTab', text: 'Open link in the same tab' },
  { value: 'newTab', text: 'Open link in a new tab' },
  { value: 'modal', text: 'Open in a modal' },
]

export const prompts = Object.freeze([
  {
    ref: 'redirect',
    meta: { short: 'Redirect user to an outside URL' },
  },
  {
    ref: 'reroute',
    meta: { short: 'Redirect user to an internal application route' },
  },
  {
    ref: 'recordPage',
    meta: { short: 'Redirect user to the record page', webapps: ['compose'] },
  },
  {
    ref: 'refetchRecords',
    meta: { short: 'Refresh all record values on the page', webapps: ['compose'] },
  },
  {
    ref: 'notification',
    meta: { short: 'Show non-blocking message to user' },
    parameters: [
      { name: 'title', types: ['String'] },
      { name: 'message', types: ['String'], required: true },
      { name: 'variant', types: ['String'], meta: { visual: { input: { type: 'select', properties: { options: variants }, default: 'primary' } } } },
      { name: 'timeout', types: ['Integer'] },
    ],
  },
  {
    ref: 'alert',
    meta: { short: 'Prompt user with an alert' },
  },
  {
    ref: 'choice',
    meta: { short: 'Prompt user with choice' },
  },
  {
    ref: 'composeRecordPicker',
    meta: { short: 'Prompt user to select a Compose Record', webapps: ['compose'] },
  },
  {
    ref: 'input',
    meta: { short: 'Prompt user with a single input' },
  },
  {
    ref: 'options',
    meta: { short: 'Prompt user with options' },
    parameters: [
      { name: 'type', types: ['String'], meta: { visual: { input: { type: 'select', properties: { options: [{ value: 'select', text: 'Select' }, { value: 'radio', text: 'Radio' }] } } } } },
      { name: 'openMode', types: ['String'], meta: { visual: { input: { type: 'select', properties: { options: openModeVariants }, default: 'sameTab' } } } },
    ],
  },
].map(definition => new automation.Function({ ...definition, kind: 'prompt' })))
