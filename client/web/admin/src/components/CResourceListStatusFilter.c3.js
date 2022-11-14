import { default as component } from './CResourceListStatusFilter'
import { components } from '@cortezaproject/corteza-vue'
const { input, select } = components.C3.controls

const options = [
  { value: '0', text: 'Excluded' },
  { value: '1', text: 'Inclusive' },
  { value: '2', text: 'Exclusive' },
]

const controls = [
  [
    input('Label', 'label'),
    select('Filter', 'value', options),
  ],
  [
    input('Excluded label', 'excluded-label'),
    input('Inclusive label', 'inclusive-label'),
    input('Exclusive label', 'exclusive-label'),
  ],
]

export default {
  name: 'Resource list status filter',
  group: ['Root components'],
  component,
  controls,
}
