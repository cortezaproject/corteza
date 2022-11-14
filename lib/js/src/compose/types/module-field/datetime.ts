// @todo option to allow only time entry
// @todo option to allow multiple entries
// @todo option to allow duplicates
import moment, { Moment } from 'moment'
import { ModuleField, Registry, Options, defaultOptions } from './base'
import { Apply } from '../../../cast'
import * as fmt from '../../../formatting'

const kind = 'DateTime'

interface DateTimeOptions extends Options {
  format: string;
  onlyDate: boolean;
  onlyTime: boolean;
  onlyPastValues: boolean;
  onlyFutureValues: boolean;
  outputRelative: boolean;
  multiDelimiter: string;
}

const defaults = (): Readonly<DateTimeOptions> => Object.freeze({
  ...defaultOptions(),
  format: '',
  multiDelimiter: '\n',
  onlyDate: false,
  onlyFutureValues: false,
  onlyPastValues: false,
  onlyTime: false,
  outputRelative: false,
})

export class ModuleFieldDateTime extends ModuleField {
  readonly kind = kind

  options: DateTimeOptions = { ...defaults() }

  constructor (i?: Partial<ModuleFieldDateTime>) {
    super(i)
    this.applyOptions(i?.options)
  }

  applyOptions (o?: Partial<DateTimeOptions>): void {
    if (!o) return
    super.applyOptions(o)

    Apply(this.options, o, String, 'format', 'multiDelimiter')
    Apply(this.options, o, Boolean, 'onlyDate', 'onlyTime', 'onlyPastValues', 'onlyFutureValues', 'outputRelative')
  }

  formatValue (value: string|Moment|Date): string|null {
    if (value === 'Invalid date') {
      return null
    }

    const o = this.options
    let m = moment(value)

    if (o.onlyTime) {
      m = moment(value, ['YYYY-MM-DD HH:mm', 'YYYY-MM-DD', 'HH:mm'])
    }

    if (o.outputRelative) {
      return m.fromNow()
    } else if (o.format.length > 0) {
      return m.format(o.format)
    } else if (o.onlyTime) {
      return fmt.time(m)
    } else if (o.onlyDate) {
      return fmt.date(m)
    } else {
      return fmt.fullDateTime(m)
    }
  }

  /**
   * Checks if given value is in the future
   * @param {String|Array<String>} v Value (in DateTime) to check
   * @param {Moment} now Time reference
   * @returns {undefined|String} undefined if valid, Error string if invalid
   */
  checkFuture (v: string|string[], now = moment()): undefined|string {
    if (!this.options.onlyFutureValues) {
      return undefined
    }

    if (!Array.isArray(v)) {
      v = [v]
    }
    if (v.find(v => moment(v) < now)) {
      return 'notification.field-datetime.valueNotFuture'
    }
    return undefined
  }

  /**
   * Checks if given value is in the past
   * @param {String|Array<String>} v Value (in DateTime) to check
   * @param {Moment} now Time reference
   * @returns {undefined|String} undefined if valid, Error string if invalid
   */
  checkPast (v: string|string[], now = moment()): undefined|string {
    if (!this.options.onlyPastValues) {
      return undefined
    }

    if (!Array.isArray(v)) {
      v = [v]
    }
    if (v.find(v => moment(v) > now)) {
      return 'notification.field-datetime.valueNotPast'
    }
  }

  /**
   * Checks if given value is valid for this field
   * @param {String} v Value (in DateTime) to check
   * @param {Moment} now Reference time used to compare
   * @returns {Array<>} Array of issues; empty if none
   */
  validate (v: string|string[], now = moment()): string[] {
    let err = this.checkFuture(v, now)
    err = err || this.checkPast(v, now)

    if (err) {
      return [err]
    }

    return []
  }
}

Registry.set(kind, ModuleFieldDateTime)
