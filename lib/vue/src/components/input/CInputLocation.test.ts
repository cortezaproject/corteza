import { describe, it, expect, vi } from 'vitest'
import { defineComponent, ref } from 'vue'
import { mount, flushPromises } from '@vue/test-utils'
import CInputLocation from './CInputLocation.vue'

// The map dialog pulls in leaflet, which jsdom cannot render.
vi.mock('../map/CMap.vue', () => ({
  default: defineComponent({ name: 'CMap', template: '<div />' }),
}))

// Mirror the real caller: the point lives outside the component and is handed
// back down on every emit.
const Parent = defineComponent({
  components: { CInputLocation },
  setup() {
    const point = ref<object | null>(null)
    return { point }
  },
  template: `<CInputLocation :model-value="point" @update:model-value="point = $event" />`,
})

const coords = (wrapper: ReturnType<typeof mount>) =>
  wrapper.findAllComponents({ name: 'InputNumber' }).map(c => c.props('modelValue'))

async function setCoord(wrapper: ReturnType<typeof mount>, index: number, value: number | null) {
  wrapper.findAllComponents({ name: 'InputNumber' })[index].vm.$emit('update:modelValue', value)
  await flushPromises()
}

describe('CInputLocation', () => {
  it('keeps the first coordinate while the second is still missing', async () => {
    const wrapper = mount(Parent)
    await flushPromises()

    // Latitude alone cannot form a point, so the model stays empty — but the
    // entered value must survive, or the longitude has nothing to pair with.
    await setCoord(wrapper, 0, 46.0569)
    expect(wrapper.vm.point).toBe(null)
    expect(coords(wrapper)).toEqual([46.0569, null])

    await setCoord(wrapper, 1, 14.5058)
    expect(wrapper.vm.point).toEqual({ type: 'Point', coordinates: [14.5058, 46.0569] })
  })

  it('builds the point when longitude is entered first', async () => {
    const wrapper = mount(Parent)
    await flushPromises()

    await setCoord(wrapper, 1, 14.5058)
    expect(coords(wrapper)).toEqual([null, 14.5058])

    await setCoord(wrapper, 0, 46.0569)
    expect(wrapper.vm.point).toEqual({ type: 'Point', coordinates: [14.5058, 46.0569] })
  })

  it('clearing one coordinate empties the point but keeps the other', async () => {
    const wrapper = mount(Parent)
    await flushPromises()

    await setCoord(wrapper, 0, 46.0569)
    await setCoord(wrapper, 1, 14.5058)
    await setCoord(wrapper, 0, null)

    expect(wrapper.vm.point).toBe(null)
    expect(coords(wrapper)).toEqual([null, 14.5058])
  })

  it('the clear button empties both coordinates', async () => {
    const wrapper = mount(Parent)
    await flushPromises()

    await setCoord(wrapper, 0, 46.0569)
    await setCoord(wrapper, 1, 14.5058)

    // Clear sits in the InputGroupAddon, ahead of the map/location buttons.
    await wrapper.findAll('button')[0].trigger('click')
    await flushPromises()

    expect(wrapper.vm.point).toBe(null)
    expect(coords(wrapper)).toEqual([null, null])
  })
})
