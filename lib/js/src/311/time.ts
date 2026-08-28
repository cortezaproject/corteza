import { BENCHMARK_TIMEZONE } from './fixtures'

export const C311_TIMEZONE = BENCHMARK_TIMEZONE

function asDate (input: string | Date): Date {
  return input instanceof Date ? input : new Date(input)
}

export function formatC311DateTime (input: string | Date, locale = 'en-US', timeZone = C311_TIMEZONE): string {
  const date = asDate(input)
  if (Number.isNaN(date.getTime())) return ''

  const parts = new Intl.DateTimeFormat(locale, {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: 'numeric',
    minute: '2-digit',
    hour12: true,
    timeZone,
    timeZoneName: 'short',
  }).formatToParts(date).reduce<Record<string, string>>((out, part) => {
    out[part.type] = part.value
    return out
  }, {})

  const period = new Intl.DateTimeFormat('en-US', {
    hour: 'numeric',
    hour12: true,
    timeZone,
  }).formatToParts(date).find(part => part.type === 'dayPeriod')?.value.toUpperCase() || 'AM'
  const zone = new Intl.DateTimeFormat('en-US', {
    timeZone,
    timeZoneName: 'short',
  }).formatToParts(date).find(part => part.type === 'timeZoneName')?.value === 'EDT' ? 'EDT' : 'EST'
  return `${parts.month}/${parts.day}/${parts.year} ${parts.hour}:${parts.minute} ${period} ${zone}`
}
export function formatC311Date (input: string | Date, locale = 'en-US', timeZone = C311_TIMEZONE): string {
  const date = asDate(input)
  if (Number.isNaN(date.getTime())) return ''
  const parts = new Intl.DateTimeFormat(locale, {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    timeZone,
  }).formatToParts(date).reduce<Record<string, string>>((out, part) => {
    out[part.type] = part.value
    return out
  }, {})
  return `${parts.month}/${parts.day}/${parts.year}`
}
