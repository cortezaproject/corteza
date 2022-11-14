import hr from 'hex-rgb'

const isLightThreshold = 100
const bgAlpha = 0.7

export const rgbaRegex = /^rgba\((\d+),.*?(\d+),.*?(\d+),.*?(\d*\.?\d*)\)$/

const ln = (n: number): number => Math.round(n < 0 ? 255 + n : (n > 255) ? n - 255 : n)
export const toRGBA = ([r, g, b, a]: number[]): string =>
  `rgba(${ln(r)}, ${ln(g)}, ${ln(b)}, ${a})`

interface Colors {
  backgroundColor: string;
  borderColor: string;
  isLight: boolean;
}

/**
 * Helper to determine event's colors
 * @param {String} hex Base color in HEX format
 * @returns {Object} { backgroundColor: String, borderColor: String, isLight: Boolean }
 */
export function makeColors (hex: string): Colors {
  const bg = hr(hex, { format: 'array' })
  const br = [...bg]
  bg[3] = bgAlpha

  return {
    backgroundColor: `rgba(${bg.join(',')})`,
    borderColor: `rgba(${br.join(',')})`,
    isLight: (bg.slice(0, 3).reduce((acc, cur) => acc + cur, 0) / (bg.length - 1)) > isLightThreshold,
  }
}

export interface Event {
  groupId?: string;
  id: string;
  title: string;
  start?: string;
  end?: string;
  allDay: boolean;
  backgroundColor: string;
  borderColor: string;
  classNames: string[];
  extendedProps: object;
}
