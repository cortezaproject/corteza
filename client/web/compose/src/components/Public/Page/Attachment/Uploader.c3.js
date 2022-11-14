// eslint-disable-next-line
import { default as component } from './Uploader.vue'
// import { components } from '@cortezaproject/corteza-vue'
// const { checkbox, input } = components.C3.controls

// Cannot read property 'accessToken' of undefined
const props = {
  endpoint: 'Endpoint',
  acceptedFiles: [],
  maxFilesize: 50,
  label: 'Label',
}

export default {
  name: 'Attachment uploader WIP',
  group: ['Public'],
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
