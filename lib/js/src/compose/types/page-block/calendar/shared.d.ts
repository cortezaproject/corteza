export declare const rgbaRegex: RegExp
export declare const toRGBA: ([r, g, b, a]: number[]) => string
interface Colors {
  backgroundColor: string
  borderColor: string
  textColor: string
}
/**
 * Helper to determine event's colors
 * @param {String} hex Base color in HEX format
 * @returns {Object} { backgroundColor: String, borderColor: String, isLight: Boolean }
 */
export declare function makeColors(hex: string): Colors
export interface Event {
  groupId?: string
  id: string
  title: string
  start?: string
  end?: string
  allDay: boolean
  backgroundColor: string
  borderColor: string
  textColor: string
  classNames: string[]
  extendedProps: object
}
export {}
