import { describe, it, expect, vi } from 'vitest'
import { defineComponent, ref, nextTick } from 'vue'
import { mount, flushPromises } from '@vue/test-utils'
import type { Map as LeafletMap } from 'leaflet'
import CMap from './CMap.vue'
import { roundLatLng } from './geo'

/**
 * Mirrors what a configurator does: it stores the centre and zoom the map
 * reports and hands them straight back as props. Leaflet answers with a centre
 * derived from the pixel origin, never bit-identical to the one it was given,
 * so without a tolerance every echo reads as a new position to move to.
 */
const Parent = defineComponent({
  components: { CMap },
  props: {
    disablePan: { type: Boolean, default: false },
  },
  setup() {
    const center = ref<[number, number]>([30, 30])
    const zoom = ref(3)
    const writes = ref(0)

    function onCenter(next: [number, number]) {
      writes.value++
      center.value = roundLatLng(next) as [number, number]
    }

    return { center, zoom, writes, onCenter }
  },
  template: `
    <CMap
      :center="center"
      :zoom="zoom"
      :disable-pan="disablePan"
      style="width: 400px; height: 300px"
      @ready="$emit('map-ready', $event)"
      @update:center="onCenter"
      @update:zoom="zoom = $event"
    />`,
})

async function mountMap(props: Record<string, unknown> = {}) {
  const wrapper = mount(Parent, { props, attachTo: document.body })
  // `ready` hands over the leaflet object, which is how a test drives the map
  // the way a user's gestures would.
  const map = wrapper.emitted('map-ready')?.[0]?.[0] as LeafletMap
  await flushPromises()
  return { wrapper, map }
}

describe('CMap', () => {
  it('settles after a parent echoes the centre it emitted', async () => {
    const { wrapper, map } = await mountMap()
    expect(map).toBeTruthy()

    // A user move: leaflet reports the new viewport, and the parent writes the
    // rounded centre straight back into the prop.
    map!.setView([46.056946, 14.505751], 10)
    const writesAfterMove = wrapper.vm.writes
    expect(writesAfterMove).toBeGreaterThan(0)

    // Spy before the echo reaches the watcher: a map that pans on its own
    // echoed centre pans forever.
    const moveSpy = vi.spyOn(map!, 'panTo')
    await flushPromises()
    await nextTick()

    expect(moveSpy).not.toHaveBeenCalled()
    expect(wrapper.vm.writes).toBe(writesAfterMove)

    wrapper.unmount()
  })

  it('emits nothing when the view has not actually changed', async () => {
    const { wrapper, map } = await mountMap()

    map!.setView([46.056946, 14.505751], 10)
    await flushPromises()
    const writes = wrapper.vm.writes

    // Same view again: leaflet still fires moveend, but there is nothing new to
    // report, so the parent must not be touched.
    map!.setView([46.056946, 14.505751], 10)
    await flushPromises()

    expect(wrapper.vm.writes).toBe(writes)

    wrapper.unmount()
  })

  it('still follows a centre the parent genuinely changes', async () => {
    const { wrapper, map } = await mountMap()

    wrapper.vm.center = [46.056946, 14.505751]
    await nextTick()
    await flushPromises()

    const c = map!.getCenter()
    expect(c.lat).toBeCloseTo(46.056946, 3)
    expect(c.lng).toBeCloseTo(14.505751, 3)

    wrapper.unmount()
  })

  it('freezes panning but keeps zoom when the view is locked', async () => {
    const { wrapper, map } = await mountMap({ disablePan: true })

    expect(map!.dragging.enabled()).toBe(false)
    expect(map!.boxZoom.enabled()).toBe(false)
    // Zoom stays usable — anchored on the centre so it cannot shift the view.
    expect(map!.options.scrollWheelZoom).toBe('center')
    expect(map!.options.doubleClickZoom).toBe('center')

    await wrapper.setProps({ disablePan: false })
    await nextTick()

    expect(map!.dragging.enabled()).toBe(true)
    expect(map!.options.scrollWheelZoom).toBe(true)

    wrapper.unmount()
  })
})
