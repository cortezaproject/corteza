// eslint-disable-next-line
import { default as component } from './NamespaceSidebar.vue'

// TypeError: Cannot read property 'getters' of undefined
const props = {
  namespaces: [
    {
      namespaceID: '1111111111',
      meta: {
        subtitle: 'Subtitle',
        description: 'Lorem ipsum dolor',
      },
      name: 'CRM',
    },
  ],
}

export default {
  name: 'Sidebar',
  group: ['Namespaces'],
  component,
  props,
  controls: [],
}
