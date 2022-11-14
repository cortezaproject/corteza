// eslint-disable-next-line
import { default as component } from './ListLoader.vue'
import { compose } from '@cortezaproject/corteza-js'

// TypeError: Cannot read property 'attachmentRead' of undefined
const props = {
  enableDelete: true,
  enableOrder: true,
  namespace: new compose.Namespace({
    namespaceID: '',
    attachmentID: '',
  }),
  kind: 'page',
  mode: 'list',
  set: ['244165968555671554'],
  hideFileName: false,
}

export default {
  name: 'List loader WIP',
  group: ['Public', 'Attachment'],
  component,
  props,
  controls: [],
}
