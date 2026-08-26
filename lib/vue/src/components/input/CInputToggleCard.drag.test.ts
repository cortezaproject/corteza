import { describe, it, expect } from 'vitest'
import { h } from 'vue'
import { mount } from '@vue/test-utils'
import Slider from 'primevue/slider'
import CInputToggleCard from './CInputToggleCard.vue'

// A drag that starts on a control inside the card and ends anywhere else fires
// its click on the two targets' common ancestor — the card — so the card sees a
// click it was never given. Releasing the temperature slider must not switch
// temperature off.

const mountCard = () =>
  mount(CInputToggleCard, {
    props: { modelValue: true, label: 'Temperature' },
    slots: {
      description: () => [
        h(Slider, { modelValue: 0.7, min: 0, max: 1, step: 0.1, id: 'temperature' }),
      ],
    },
    attachTo: document.body,
  })

const card = (w: ReturnType<typeof mountCard>) => w.element as HTMLElement
const handle = (w: ReturnType<typeof mountCard>) =>
  w.find('[data-pc-section="handle"]').element as HTMLElement

const emitted = (w: ReturnType<typeof mountCard>) => w.emitted('update:modelValue') || []

describe('CInputToggleCard with a control in a slot', () => {
  it('does not toggle when a drag started inside the card releases over it', async () => {
    const w = mountCard()
    handle(w).dispatchEvent(new MouseEvent('pointerdown', { bubbles: true }))
    // the drag ends off the slider, so the browser fires the click on the card
    card(w).dispatchEvent(new MouseEvent('click', { bubbles: true }))
    await w.vm.$nextTick()
    expect(emitted(w).length).toBe(0)
  })

  it('still toggles on a plain click on the card', async () => {
    const w = mountCard()
    const label = w.find('span').element as HTMLElement
    label.dispatchEvent(new MouseEvent('pointerdown', { bubbles: true }))
    label.dispatchEvent(new MouseEvent('click', { bubbles: true }))
    await w.vm.$nextTick()
    expect(emitted(w).length).toBe(1)
    expect(emitted(w)[0]).toEqual([false])
  })

  it('stays inert when disabled', async () => {
    const w = mount(CInputToggleCard, {
      props: { modelValue: true, label: 'Temperature', disabled: true },
      attachTo: document.body,
    })
    const label = w.find('span').element as HTMLElement
    label.dispatchEvent(new MouseEvent('pointerdown', { bubbles: true }))
    label.dispatchEvent(new MouseEvent('click', { bubbles: true }))
    await w.vm.$nextTick()
    expect(w.emitted('update:modelValue')).toBeUndefined()
  })
})
