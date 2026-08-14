import { describe, it, expect } from 'vitest'
import {
  boundsRing,
  COORD_EPSILON,
  parseBounds,
  parseLatLng,
  roundBounds,
  roundLatLng,
  sameBounds,
  sameLatLng,
} from './geo'

describe('parseLatLng', () => {
  it('rejects anything that is not a finite pair', () => {
    expect(parseLatLng([46.05, 14.5])).toEqual([46.05, 14.5])
    expect(parseLatLng([46.05])).toBe(null)
    expect(parseLatLng([46.05, '14.5'])).toBe(null)
    expect(parseLatLng([NaN, 14.5])).toBe(null)
    expect(parseLatLng(null)).toBe(null)
  })
})

describe('parseBounds', () => {
  it('normalises corners so the same area always compares equal', () => {
    const drawnNorthWest = parseBounds([
      [47, 16],
      [46, 14],
    ])
    const drawnSouthEast = parseBounds([
      [46, 14],
      [47, 16],
    ])
    expect(drawnNorthWest).toEqual([
      [46, 14],
      [47, 16],
    ])
    expect(drawnNorthWest).toEqual(drawnSouthEast)
  })

  it('rejects a half-formed box', () => {
    expect(parseBounds([[46, 14]])).toBe(null)
    expect(parseBounds([[46, 14], null])).toBe(null)
    expect(parseBounds(null)).toBe(null)
  })
})

describe('sameLatLng', () => {
  // This is the guard that stops a map echoing its own emitted centre back into
  // a pan. Leaflet returns a centre derived from the pixel origin, so the value
  // that comes back is a hair off the one we persisted.
  it('treats a sub-epsilon difference as the same place', () => {
    expect(sameLatLng([46.056946, 14.505751], [46.0569455, 14.5057513])).toBe(true)
  })

  it('treats a difference above epsilon as a move', () => {
    expect(sameLatLng([46.056946, 14.505751], [46.056946 + COORD_EPSILON * 10, 14.505751])).toBe(
      false,
    )
  })

  // "nothing changed" is the question being asked, so two absent positions are
  // equal — but an absent one is never the same as a real one.
  it('never mistakes an absent position for a real one', () => {
    expect(sameLatLng(null, [46, 14])).toBe(false)
    expect(sameLatLng([46, 14], undefined)).toBe(false)
    expect(sameLatLng(null, null)).toBe(true)
  })
})

describe('sameBounds', () => {
  it('compares normalised corners', () => {
    expect(
      sameBounds(
        [
          [47, 16],
          [46, 14],
        ],
        [
          [46, 14],
          [47.0000001, 16],
        ],
      ),
    ).toBe(true)
  })

  it('sees a real viewport change', () => {
    expect(
      sameBounds(
        [
          [46, 14],
          [47, 16],
        ],
        [
          [46, 14],
          [48, 16],
        ],
      ),
    ).toBe(false)
  })
})

describe('rounding', () => {
  it('keeps six decimals', () => {
    expect(roundLatLng([46.0569456789, 14.5057512345])).toEqual([46.056946, 14.505751])
    expect(
      roundBounds([
        [46.0000004, 14.0000004],
        [47.0000006, 16.0000006],
      ]),
    ).toEqual([
      [46, 14],
      [47.000001, 16.000001],
    ])
  })

  it('passes unusable input through as null', () => {
    expect(roundLatLng(['a', 'b'])).toBe(null)
    expect(roundBounds([[46, 14]])).toBe(null)
  })

  it('stays within epsilon of the value it rounded', () => {
    const raw: [number, number] = [46.0569456789, 14.5057512345]
    expect(sameLatLng(roundLatLng(raw), raw)).toBe(true)
  })
})

describe('boundsRing', () => {
  it('walks the four corners so a polygon can draw the locked area', () => {
    expect(
      boundsRing([
        [46, 14],
        [47, 16],
      ]),
    ).toEqual([
      [46, 14],
      [47, 14],
      [47, 16],
      [46, 16],
    ])
  })

  it('has nothing to draw without a box', () => {
    expect(boundsRing(null)).toBe(null)
  })
})
