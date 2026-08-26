import { BENCHMARK_TIMEZONE } from './fixtures'

export const C311_TIMEZONE = BENCHMARK_TIMEZONE

function asDate (input: string | Date): Date {
  return input instanceof Date ? input : new Date(input)
}

export function formatC311DateTime (input: string | Date, locale = 'en-US', timeZone = C311_TIMEZONE): string {
  const date = asDate(input)
  if (Number.isNaN(date.getTime())) return ''
  return new Intl.DateTimeFormat(locale, {
    dateStyle: 'medium',
    timeStyle: 'short',
    timeZone,
  }).format(date)
}
export function formatC311Date (input: string | Date, locale = 'en-US', timeZone = C311_TIMEZONE): string {
  const date = asDate(input)
  if (Number.isNaN(date.getTime())) return ''
  return new Intl.DateTimeFormat(locale, {
    dateStyle: 'medium',
    timeZone,
  }).format(date)
}
