import { currentLanguage } from './locale'

export function number (input: number, options: Intl.NumberFormatOptions = { maximumFractionDigits: 6 }): string {
  return new Intl.NumberFormat(currentLanguage(), options).format(input)
}

export function accountingNumber (value: number): string {
  if (value === 0) {
    return '-'
  }

  if (value < 0) {
    return `(${number(Math.abs(value))})`
  }

  return number(value)
}
