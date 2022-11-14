// eslint-disable-next-line
import { default as component } from './FieldRowEdit.vue'
import { compose } from '@cortezaproject/corteza-js'
import { components } from '@cortezaproject/corteza-vue'
const { checkbox } = components.C3.controls

const props = {
  value: {
    name: 'Name',
    label: 'Label',
    kind: 'String',
    cap: {
      configurable: true,
      multi: true,
      required: true,
      private: false,
    },
    isMulti: true,
    isRequired: true,
    isPrivate: false,
    fieldId: '453534534534343',
    isValid: true,
  },
  module: new compose.Module({
    namespaceID: '3242343234',
    moduleID: '999999442',
  }),
  isDuplicate: false,
  canGrant: true,
  hasRecords: false,
}

export default {
  name: 'Field row edit',
  group: ['Admin', 'Module'],
  component,
  props,
  controls: [
    checkbox('configurable', 'cap.configurable'),
    checkbox('multi', 'value.cap.multi'),
    checkbox('required', 'value.cap.required'),
    checkbox('private', 'value.cap.private'),
    checkbox('isMulti', 'value.isMulti'),
    checkbox('isRequired', 'value.isRequired'),
    checkbox('isPrivate', 'value.isPrivate'),
    checkbox('isValid', 'value.isValid'),
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
        canGrant: false,
        hasRecords: true,
        value: {
          cap: {
            configurable: false,
            multi: false,
            required: false,
            private: false,
          },
          isMulti: false,
          isRequired: false,
          isPrivate: false,
          isValid: false,
        },
      },
    },
  ],
}
