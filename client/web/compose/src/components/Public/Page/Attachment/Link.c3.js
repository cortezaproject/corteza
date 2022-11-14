// eslint-disable-next-line
import { default as component } from './Link.vue'

const props = {
  attachment: {
    name: 'File',
    download: 'download link',
    meta: {
      original: '',
    },
  },
}

export default {
  name: 'Attachment link',
  group: ['Public'],
  component,
  props,
  controls: [],
}
