export function formatValueAsAccountingNumber (value: number): string {
  let result = ''

  if (value < 0) {
    result = `(${Math.abs(value)})`
  } else if (value === 0) {
    result = '-'
  }

  return result
}
