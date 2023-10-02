// eslint-disable-next-line
import { default as component } from './AutomationTab.vue'
import { compose } from '@cortezaproject/corteza-js'

// Error in render: 'TypeError: Cannot read property 'Find' of undefined
const props = {
  buttons: [
    {
      label: 'Dummy',
      variant: 'danger',
    },
    {
      button: {
        script: null,
      },
    },
  ],
  namespace: new compose.Namespace(),
  module: new compose.Module(),
  field: new compose.ModuleFieldBool(),
  blockIndex: -1,
  page: new compose.Page(),
  block: new compose.PageBlock(),
  // Uncaught Error: invalid module used to initialize a record
  // record: new compose.Record(),
  mode: '',
}

export default {
  name: 'Automation tab WIP',
  group: ['PageBlocks', 'Shared'],
  component,
  props,
  controls: [],
  plugins: {
    $UIHooks: {
      Find () { return [] },
    },
  },
  // controls: [
  //   checkbox('Done', 'done'),
  //   checkbox('Disabled', 'disabled'),
  //   checkbox('Small', 'small'),
  //   input('Number of step', 'stepNumber'),
  // ],

  // scenarios: [
  //   {
  //     label: 'Full form',
  //     props,
  //   },
  //   {
  //     label: 'Empty form',
  //     props: {
  //       ...props,
  //       canGrant: false,
  //       hasRecords: false,
  //     },
  //   },
  // ],
}
