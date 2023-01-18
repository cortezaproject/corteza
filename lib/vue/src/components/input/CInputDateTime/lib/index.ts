import moment from 'moment'

export function getDate (value: string|undefined): string | undefined {
  if (!value) {
    return undefined
  }

  if (value === 'Invalid date') {
    // Make sure this weird value does not cause us problems
    return undefined
  }

  if (moment(value, 'YYYY-MM-DDTHH:mm:ssZ', true).isValid()) {
    return moment.utc(value).local().format('YYYY-MM-DD')
  }

  return value
}

export function getTime (value: string|undefined): string | undefined {
  if (!value) {
    return undefined
  }

  if (moment(value, 'YYYY-MM-DDTHH:mm:ssZ', true).isValid()) {
    return moment.utc(value).local().format('HH:mm')
  }

  return value
}

export function setDate (date: string|undefined, value: string|undefined, noDate = false, noTime = false): string | undefined {
  if (noDate || !date || !date.length) {
    return undefined
  }

  if (noTime) {
    return moment(date, 'YYYY-MM-DD').format('YYYY-MM-DD')
  }

  const time = getTime(value) || '00:00'

  return moment(date + ' ' + time, 'YYYY-MM-DD HH:mm').utc().format()
}

export function setTime (time: string|undefined, value: string|undefined, noDate = false, noTime = false): string | undefined {
  if (noTime || !time || !time.length) {
    return undefined
  }

  if (noDate) {
    return moment(time, 'HH:mm').format('HH:mm')
  }

  // Default to today if date not provided
  const date = getDate(value) || moment().local().format('YYYY-MM-DD')

  return moment(date + ' ' + time, 'YYYY-MM-DD HH:mm').utc().format()
}
