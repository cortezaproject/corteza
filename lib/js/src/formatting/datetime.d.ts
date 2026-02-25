import { MomentInput } from 'moment'
declare type DateTimeInput = unknown | MomentInput
declare type DateTimeFormatOptions = Intl.DateTimeFormatOptions & {
  dateStyle?: 'full' | 'long' | 'medium' | 'short'
  timeStyle?: 'full' | 'long' | 'medium' | 'short'
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
export declare function fullDateTime(input: DateTimeInput, options?: DateTimeFormatOptions): string
/**
 * Outputs locally formatted date without time
 *
 * Example:
 * 09/04/1986
 *
 * @param input
 * @param options
 */
export declare function date(input: DateTimeInput, options?: DateTimeFormatOptions): string
/**
 * Outputs locally formatted time
 *
 * Example:
 * 8:30 PM
 *
 * @param input
 * @param options
 */
export declare function time(input: DateTimeInput, options?: DateTimeFormatOptions): string
export {}
