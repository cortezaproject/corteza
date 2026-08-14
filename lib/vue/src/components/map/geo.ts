/**
 * Coordinate helpers shared by CMap and by anything that stores a map view
 * (the geometry field and page-block configurators).
 *
 * They live outside the SFC because the comparison maths is what keeps a
 * two-way bound map from oscillating, and that is worth testing without a
 * browser or a leaflet instance in the way.
 */

export type LatLng = [number, number]
export type Bounds = [LatLng, LatLng]

/**
 * A coordinate that has been through leaflet comes back derived from the pixel
 * origin, so an echoed value is never bit-identical to the one we sent. 1e-6
 * degrees is ~11cm: below that two coordinates are the same place, and treating
 * them as different is exactly what turns a prop echo into an endless
 * pan → moveend → emit → pan loop.
 */
export const COORD_EPSILON = 1e-6

/** Decimals kept when a coordinate is persisted — matches COORD_EPSILON. */
export const COORD_PRECISION = 6

function isFiniteNumber(v: unknown): v is number {
  return typeof v === 'number' && Number.isFinite(v)
}

export function isLatLng(v: unknown): v is LatLng {
  return Array.isArray(v) && v.length === 2 && isFiniteNumber(v[0]) && isFiniteNumber(v[1])
}

/** [lat, lng] or null — never a partially valid pair. */
export function parseLatLng(v: unknown): LatLng | null {
  return isLatLng(v) ? [v[0], v[1]] : null
}

/**
 * [[swLat, swLng], [neLat, neLng]] or null. Corners are normalised so that a
 * viewport captured while dragging north-west compares equal to the same area
 * captured while dragging south-east.
 */
export function parseBounds(v: unknown): Bounds | null {
  if (!Array.isArray(v) || v.length !== 2) return null
  const a = parseLatLng(v[0])
  const b = parseLatLng(v[1])
  if (!a || !b) return null
  return [
    [Math.min(a[0], b[0]), Math.min(a[1], b[1])],
    [Math.max(a[0], b[0]), Math.max(a[1], b[1])],
  ]
}

/**
 * Both callers use this to decide "has anything actually changed?", so two
 * absent positions count as equal and an absent one never equals a real one.
 */
export function sameLatLng(a: unknown, b: unknown, epsilon: number = COORD_EPSILON): boolean {
  const l = parseLatLng(a)
  const r = parseLatLng(b)
  if (!l || !r) return !l && !r
  return Math.abs(l[0] - r[0]) <= epsilon && Math.abs(l[1] - r[1]) <= epsilon
}

export function sameBounds(a: unknown, b: unknown, epsilon: number = COORD_EPSILON): boolean {
  const l = parseBounds(a)
  const r = parseBounds(b)
  if (!l || !r) return !l && !r
  return sameLatLng(l[0], r[0], epsilon) && sameLatLng(l[1], r[1], epsilon)
}

export function roundCoord(n: number, precision: number = COORD_PRECISION): number {
  const factor = 10 ** precision
  return Math.round(n * factor) / factor
}

/** Persisted coordinates are rounded so a stored view stays human-readable. */
export function roundLatLng(v: unknown, precision: number = COORD_PRECISION): LatLng | null {
  const pair = parseLatLng(v)
  return pair ? [roundCoord(pair[0], precision), roundCoord(pair[1], precision)] : null
}

export function roundBounds(v: unknown, precision: number = COORD_PRECISION): Bounds | null {
  const bounds = parseBounds(v)
  if (!bounds) return null
  return [roundLatLng(bounds[0], precision) as LatLng, roundLatLng(bounds[1], precision) as LatLng]
}

/**
 * The four corners of a bounds box, as a closed ring a polygon layer can draw.
 * Used to show which area a locked view actually covers.
 */
export function boundsRing(v: unknown): LatLng[] | null {
  const bounds = parseBounds(v)
  if (!bounds) return null
  const [[swLat, swLng], [neLat, neLng]] = bounds
  return [
    [swLat, swLng],
    [neLat, swLng],
    [neLat, neLng],
    [swLat, neLng],
  ]
}
