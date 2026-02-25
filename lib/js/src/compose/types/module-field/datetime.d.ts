import moment, { Moment } from 'moment'
import { ModuleField, Options } from './base'
interface DateTimeOptions extends Options {
  format: string
  onlyDate: boolean
  onlyTime: boolean
  onlyPastValues: boolean
  onlyFutureValues: boolean
  outputRelative: boolean
  multiDelimiter: string
}
export declare class ModuleFieldDateTime extends ModuleField {
  readonly kind = 'DateTime'
  options: DateTimeOptions
  constructor(i?: Partial<ModuleFieldDateTime>)
  applyOptions(o?: Partial<DateTimeOptions>): void
  formatValue(value: string | Moment | Date): string | null
  /**
   * Checks if given value is in the future
   * @param {String|Array<String>} v Value (in DateTime) to check
   * @param {Moment} now Time reference
   * @returns {undefined|String} undefined if valid, Error string if invalid
   */
  checkFuture(v: string | string[], now?: moment.Moment): undefined | string
  /**
   * Checks if given value is in the past
   * @param {String|Array<String>} v Value (in DateTime) to check
   * @param {Moment} now Time reference
   * @returns {undefined|String} undefined if valid, Error string if invalid
   */
  checkPast(v: string | string[], now?: moment.Moment): undefined | string
  /**
   * Checks if given value is valid for this field
   * @param {String} v Value (in DateTime) to check
   * @param {Moment} now Reference time used to compare
   * @returns {Array<>} Array of issues; empty if none
   */
  validate(v: string | string[], now?: moment.Moment): string[]
}
export {}
