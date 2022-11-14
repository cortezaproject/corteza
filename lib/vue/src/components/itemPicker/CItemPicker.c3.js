// eslint-disable-next-line
import { default as component } from './CItemPicker.vue'
import { components } from '@cortezaproject/corteza-vue'
const { checkbox } = components.C3.controls

const props = {
  options: [
    { value: 'pl', text: 'pl' },
    { value: 'pl2', text: 'pl' },
    { value: 'pl3', text: 'pl' },
    { value: 'oo', text: 'oo' },
  ],
  value: [
    'pl',
    'pl2',
    'pl3',
  ],
  labels: {
    searchPlaceholder: 'Type here to search among module fields',
    availableItems: 'Available fields',
    selectAllItems: 'Select all',
    selectedItems: 'Selected fields',
    systemItem: '(system field)',
    unselectAllItems: 'Unselect all',
    noItemsFound: 'No Items Found',
  },

  hideFilter: false,
  hideIcons: false,
  disabled: false,
  disabledFilter: false,
  disabledSorting: false,
  disabledDragging: false,
}

export default {
  name: 'CItemPicker',
  group: ['Picker'],
  component,
  props,
  controls: [
    checkbox('Hide filter', 'hideFilter'),
    checkbox('Hide icons', 'hideIcons'),
    checkbox('Disabled', 'disabled'),
    checkbox('Disabled filter', 'disabledFilter'),
    checkbox('Disabled sorting', 'disabledSorting'),
    checkbox('Disabled dragging', 'disabledDragging'),
  ],
}
