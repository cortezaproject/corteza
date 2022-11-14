import { fmt } from '@cortezaproject/corteza-js'

export function locFullDateTime (d: unknown): string {
  // Thursday, September 4, 1986 8:30 PM
  // return moment(d).format('LLLL')
  return fmt.fullDateTime(d)
}

export function locDate (d: unknown): string {
  // 09/04/1986
  return fmt.date(d)
}
