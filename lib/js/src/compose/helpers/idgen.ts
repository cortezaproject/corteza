export function generateUID (): string {
  const uid = Math.random().toString(36).substring(2) + (new Date()).getTime().toString(36)
  return `tempID-${uid}`
}
