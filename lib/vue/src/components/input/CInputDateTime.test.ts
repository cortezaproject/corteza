import { describe, it, expect, beforeAll } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import CInputDateTime from './CInputDateTime.vue'

beforeAll(() => {
  // PrimeVue's DatePicker binds a matchMedia listener on mount; jsdom has none.
  // @ts-expect-error assigning a stub over a property jsdom does not define
  window.matchMedia = (q: string) => ({
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

const shown = (props: Record<string, unknown>) => {
  const w = mount(CInputDateTime, { props })
  return (w.find('input').element as HTMLInputElement).value
}

// Every date/time input in the app renders through this component, so the two
// things worth pinning are what it puts on screen for each accepted input shape
// and what it hands back to the caller.
describe('CInputDateTime rendering', () => {
  it('renders a string and the equivalent Date identically', () => {
    const iso = '2026-08-20T14:30:00Z'
    expect(shown({ modelValue: iso })).toBe(shown({ modelValue: new Date(iso) }))
  })

  it('shows the time for a datetime value', () => {
    // Handed the raw string, PrimeVue would format it through a date-only
    // branch (and throw once show-time is on) — the conversion is what keeps
    // the field from reformatting itself as soon as it is edited.
    expect(shown({ modelValue: '2026-08-20T14:30:00Z' })).toMatch(/\d{1,2}:\d{2}/)
  })

  it('drops the time for a date-only field', () => {
    expect(shown({ modelValue: '2026-08-20', onlyDate: true })).not.toMatch(/\d{1,2}:\d{2}/)
  })

  it('renders an empty value as an empty input', () => {
    expect(shown({ modelValue: '' })).toBe('')
    expect(shown({ modelValue: null })).toBe('')
  })
})

describe('CInputDateTime value contract', () => {
  const emitFrom = async (props: Record<string, unknown>, value: Date | null) => {
    const w = mount(CInputDateTime, { props })
    await flushPromises()
    w.findComponent({ name: 'DatePicker' }).vm.$emit('update:modelValue', value)
    await flushPromises()
    return w.emitted('update:modelValue')?.at(-1)?.[0]
  }

  it('emits an ISO string for a datetime field', async () => {
    const picked = new Date('2026-08-20T14:30:00Z')
    expect(await emitFrom({ modelValue: '' }, picked)).toBe(picked.toISOString())
  })

  it('emits YYYY-MM-DD for a date-only field', async () => {
    const picked = new Date(2026, 7, 20, 14, 30)
    expect(await emitFrom({ modelValue: '', onlyDate: true }, picked)).toBe('2026-08-20')
  })

  it('emits HH:mm:ss for a time-only field', async () => {
    const picked = new Date(2026, 7, 20, 14, 30, 5)
    expect(await emitFrom({ modelValue: '', timeOnly: true }, picked)).toBe('14:30:05')
  })

  it('emits the Date itself when the caller asked for one', async () => {
    const picked = new Date('2026-08-20T14:30:00Z')
    expect(await emitFrom({ modelValue: null, valueType: 'date' }, picked)).toBe(picked)
  })

  it('clears to the shape the caller uses', async () => {
    expect(await emitFrom({ modelValue: '2026-08-20T14:30:00Z' }, null)).toBe('')
    expect(await emitFrom({ modelValue: new Date(), valueType: 'date' }, null)).toBeNull()
  })
})
