export function Truncate(s: string, l: number): string {
  if (!s) {
    return ''
  }
  return s.length > l ? s.substring(0, l) + '...' : s
}