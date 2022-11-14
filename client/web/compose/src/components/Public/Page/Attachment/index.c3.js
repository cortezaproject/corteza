// eslint-disable-next-line
import { default as component } from './index.vue'
// import { components } from '@cortezaproject/corteza-vue'
// const { checkbox, input } = components.C3.controls

// Property or method "isImage" is not defined on the instance but referenced during render
const props = {
  kind: '',
  mode: 'list',
  value: {
    download: 'dfs',
    name: 'ss',
    size: 32,
    changedAt: 'dfc',
    previewUrl: 'fsdfs',
  },
}

export default {
  name: 'Index WIP',
  group: ['Public', 'Attachment'],
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
