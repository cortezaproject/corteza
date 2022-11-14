// eslint-disable-next-line
import { default as component } from './NamespaceItem.vue'
import { components } from '@cortezaproject/corteza-vue'
const { checkbox, input } = components.C3.controls

const props = {
  namespace: {
    meta: {
      subtitle: 'Subtitle',
      description: 'Lorem ipsum dolor',
    },
    name: 'CRM',
    enabled: true,
    canUpdateNamespace: true,
  },
}

export default {
  name: 'Item',
  group: ['Namespaces'],
  component,
  props,
  controls: [
    input('Subtitle', 'namespace.meta.subtitle'),
    input('Description', 'namespace.meta.description'),
    input('Name', 'namespace.name'),
    checkbox('Enable visit namespace button', 'namespace.enabled'),
    checkbox('Update namespace button', 'namespace.canUpdateNamespace'),
  ],

  scenarios: [
    {
      label: 'Full form',
      props,
    },
    {
      label: 'Empty form',
      props: {
        ...props,
        meta: {
          subtitle: '',
          description: '',
        },
        name: 'CRM',
        enabled: false,
        canUpdateNamespace: false,
      },
    },
  ],
}
