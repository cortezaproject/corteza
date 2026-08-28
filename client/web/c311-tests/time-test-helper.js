export function formatC311DateTime (input) {
  const date = new Date(input)
  if (Number.isNaN(date.getTime())) return ''
  const parts = new Intl.DateTimeFormat('en-US', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: 'numeric',
    minute: '2-digit',
    hour12: true,
    timeZone: 'America/New_York',
    timeZoneName: 'short',
  }).formatToParts(date).reduce((out, part) => ({ ...out, [part.type]: part.value }), {})
  return `${parts.month}/${parts.day}/${parts.year} ${parts.hour}:${parts.minute} ${parts.dayPeriod} ${parts.timeZoneName}`
}
