// eslint-disable-next-line
import { default as component } from './FileUpload.vue'
// import { components } from '@cortezaproject/corteza-vue'
// const { checkbox, input } = components.C3.controls

// "TypeError: Cannot read property 'recordImportInitEndpoint' of undefined
const props = {
  namespace: {},
  module: {},
}

export default {
  name: 'File upload WIP',
  group: ['Public Record'],
  component,
  props,
  controls: [],
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
