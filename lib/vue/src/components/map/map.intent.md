---
kind: folder
covers: recursive
backfilled: true
owner: fe
depends-on: []
touched-by: []
tests:
  - lib/vue/src/components/map/CMap.test.ts
  - lib/vue/src/components/map/geo.test.ts
---

# map/

## Intention

The one shared Leaflet map (OSM tiles) so geometry fields, location inputs,
and map blocks all get identical geosearch, themed controls, and marker
rendering. Exported from the components barrel like every other shared
component (historically deep-imported; fixed as an oversight).

## Map

- `CMap.vue` — owns a leaflet map directly (`L.map()`, no wrapper library):
  OSM tile layer, geosearch box (leaflet-geosearch), PrimeVue zoom/locate
  controls, imperative marker/polygon rendering into a layerGroup.
- `geo.ts` — pure coordinate helpers: parse, normalise, compare, round a
  lat-lng or a bounds box, and walk a box's corners for drawing. Re-exported
  from the package entry as `mapGeo`, so whatever stores a map view stores it
  in the shape CMap reads back.

## Contracts consumers rely on

- Declarative in, evented out. The `markers` and `polygons` props are
  re-rendered on change and invalid entries are silently dropped; a marker is
  `{value, title?, color?}`, a polygon
  `{latLngs, color?, weight?, dashArray?, fillOpacity?, interactive?}`. Emits
  `ready` (the leaflet object), `map-click`, `marker-click`, `location-found`,
  `update:center/zoom/bounds`.
- `update:bounds` fires once as soon as the map exists, before any gesture.
  The viewport is a function of container size, so only the map knows it, and
  a consumer that captures an area needs it without having to move the map
  first.
- `maxBounds` goes straight to leaflet, which means a bounded view rather than
  a frozen one: dragging past the box snaps back, and zooming in leaves room to
  move around inside it. There is no separate "don't move" mode — a caller that
  wants the view held still hands over the box it wants held.
- A box also sets the zoom floor (`getBoundsZoom`, recomputed whenever the box,
  the authored limits or the container size change), since below it the map
  would show ground the box exists to keep off screen. An authored `minZoom`
  can raise that floor, never lower it.
- Geosearch picks emit a synthetic `map-click` — pick-a-location flows get it
  for free; value commit stays the caller's job.
- Exposes `invalidateSize()` / `fitBounds()`; size is auto-invalidated via
  ResizeObserver.

## When changing this

- Leaflet owns the viewport and props only nudge it. Two rules keep that
  one-way: a value is emitted only when it differs from the last one the
  parent was given, and a prop that matches where the map already is never
  moves it (`geo.ts`, epsilon 1e-6 ≈ 11cm). Both are load-bearing. A centre
  that round-trips through leaflet comes back derived from the pixel origin
  and never bit-identical, so a parent storing what it was told hands back
  what reads as a new position — the map pans, emits, is told again, forever
  ("Maximum recursive updates exceeded"). `maxBounds` makes it permanent by
  re-panning on every `moveend`. Any new two-way prop needs the same pair.
- The leaflet object lives in a `shallowRef`: a deep reactive proxy breaks
  leaflet's identity-sensitive internals.
- Marker clicks stop propagation so they don't double-fire `map-click`.
