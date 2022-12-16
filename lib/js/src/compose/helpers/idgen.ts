export function generateUID (): string {
  let uid = Math.random().toString(36).substring(2) + (new Date()).getTime().toString(36)
  return `tempID-${uid}`
}