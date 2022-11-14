// eslint-disable-next-line
import { default as component } from './Tree.vue'
import { compose } from '@cortezaproject/corteza-js'
// import { components } from '@cortezaproject/corteza-vue'
// const { checkbox, input } = components.C3.controls

// TypeError: Cannot read property 'getters' of undefined at VueComponent.mappedGetter
const props = {
  namespace: new compose.Namespace({
    canCreatePage: true,
    canGrant: true,
    namespaceID: '',
  }),
  value: [
    {
      pageID: '44444444444',
      moduleID: '324234324',
      title: 'Home',
      handle: 'Home',
      visible: true,
      canGrant: true,
      canUpdatePage: true,
      blocks: [{
        kind: 'RecordList',
        title: 'My New Leads',
      }],
    },
    {
      pageID: '456454644444',
      moduleID: '004234324',
      title: 'Test',
      handle: 'Test',
      visible: true,
      canGrant: true,
      canUpdatePage: true,
      blocks: [{
        kind: 'RecordList',
        title: 'My New Leads',
      }],
    },
  ],
  parentID: '',
  level: 0,
}

export default {
  name: 'Tree WIP',
  group: ['Admin', 'Page'],
  component,
  props,
  controls: [],
  // controls: [
  //   checkbox('', ''),
  //   input(' of step', ''),
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
  //     },
  //   },
  // ],
}
