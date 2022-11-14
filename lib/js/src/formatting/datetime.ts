import moment, { MomentInput } from 'moment'
import { currentLanguage } from './locale'

declare type DateTimeInput = unknown | MomentInput

declare type DateTimeFormatOptions = Intl.DateTimeFormatOptions & {
  dateStyle?: 'full' | 'long' | 'medium' | 'short';
  timeStyle?: 'full' | 'long' | 'medium' | 'short';
}

/**
 * Parses input into Date using Moment library
 *
 * @param input
 */
function parse (input: DateTimeInput): Date {
  return moment(input as MomentInput).toDate()
}

function format (input: DateTimeInput, options: DateTimeFormatOptions): string {
  return (new Intl.DateTimeFormat(currentLanguage(), options)).format(parse(input))
}

/**
 * Outputs locally formatted date and time, no seconds
 *
 * Examples:
 * "Wednesday, September 8, 2021 at 9:41 AM"
 * "sreda, 08. september 2021 09:41"
 * "srijeda, 8. rujna 2021. u 09:42"
 *
 * @param input
 * @param options
 */
export function fullDateTime (input: DateTimeInput, options: DateTimeFormatOptions = { dateStyle: 'full', timeStyle: 'short' }): string {
  return format(input, options)
}

/**
 * Outputs locally formatted date without time
 *
 * Example:
 * 09/04/1986
 *
 * @param input
 * @param options
 */
export function date (input: DateTimeInput, options: DateTimeFormatOptions = { dateStyle: 'short' }): string {
  return format(input, options)
}

/**
 * Outputs locally formatted time
 *
 * Example:
 * 8:30 PM
 *
 * @param input
 * @param options
 */
export function time (input: DateTimeInput, options: DateTimeFormatOptions = { timeStyle: 'short' }): string {
  return format(input, options)
}
