import { currentLanguage } from './locale'

export function number (input: number, options: Intl.NumberFormatOptions = { maximumFractionDigits: 6 }): string {
  return new Intl.NumberFormat(currentLanguage(), options).format(input)
}
