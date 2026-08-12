import { describe, it, expect, beforeAll, vi } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import PrimeVue from 'primevue/config'
import DatePicker from 'primevue/datepicker'
import InputText from 'primevue/inputtext'

vi.mock('vue-i18n', () => ({
  useI18n: () => ({ t: k => k, locale: { value: 'en' } }),
}))

import GovernanceForm from './GovernanceForm.vue'

beforeAll(() => {
  // PrimeVue's DatePicker binds a matchMedia listener on mount; jsdom has none.
  window.matchMedia = q => ({
    matches: false,
    media: q,
    onchange: null,
    addEventListener() {},
    removeEventListener() {},
    addListener() {},
    removeListener() {},
    dispatchEvent: () => false,
  })
})

// The dashboard read path (stores/events.js mapEventRow) spreads the API row
// straight through, so a saved `dateDue`/`completedDate` reaches this form as
// the plain ISO string the backend stores. Only the WRITE path (normalizeDates)
// knows about Date objects. Handed that string, a bare PrimeVue DatePicker with
// show-time throws outright ('Invalid arguments' out of parseDate), and without
// show-time it renders through a branch that drops the time — so the field
// reformats itself the moment it is touched. CInputDateTime converts first.
const isoString = '2026-08-20T14:30:00Z'
const equivalentDate = new Date(isoString)

const schema = [
  {
    fields: [
      { key: 'when', labelKey: 'when', type: 'datetime' },
      { key: 'day', labelKey: 'day', type: 'date' },
    ],
  },
]

const mountForm = modelValue =>
  mount(GovernanceForm, {
    props: { schema, modelValue },
    global: {
      plugins: [[PrimeVue, { unstyled: true }]],
      components: { DatePicker, InputText },
      mocks: { $t: k => k },
      stubs: {
        EventBadge: true,
        RiskPips: true,
        ValidationMessage: true,
        CFormGroup: { template: '<div><slot /></div>' },
        Textarea: true,
        InputNumber: true,
        Select: true,
        MultiSelect: true,
      },
    },
  })

describe('GovernanceForm date fields', () => {
  it('renders a stored ISO string exactly as it renders the equivalent Date', async () => {
    const fromString = mountForm({ when: isoString, day: isoString })
    const fromDate = mountForm({ when: equivalentDate, day: equivalentDate })
    await flushPromises()

    const values = w => w.findAll('input').map(i => i.element.value)

    expect(values(fromString)).toEqual(values(fromDate))
  })

  it('keeps the datetime field showing its time', async () => {
    const wrapper = mountForm({ when: isoString })
    await flushPromises()

    // 16:30 local for a +02:00 zone, but the assertion only cares that a time
    // is present at all: the string branch of PrimeVue's formatValue drops it.
    expect(wrapper.find('input').element.value).toMatch(/\d{1,2}:\d{2}/)
  })
})
